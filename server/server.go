package server

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/MSSkowron/MSCache/client"
	"github.com/MSSkowron/MSCache/server/cache"
	"github.com/MSSkowron/MSCache/server/logger"
	"github.com/MSSkowron/MSCache/server/protocol"
)

type ServerNode struct {
	listenAddress string
	leaderAddress string
	isLeader      bool
	followers     map[*client.Client]struct{}
	leader        *client.Client
	cache         cache.Cache
}

func New(listenAddress, leaderAddress string, isLeader bool, c cache.Cache) *ServerNode {
	return &ServerNode{
		listenAddress: listenAddress,
		leaderAddress: leaderAddress,
		isLeader:      isLeader,
		followers:     nil,
		leader:        nil,
		cache:         c,
	}
}

func (s *ServerNode) Run() error {
	ln, err := net.Listen("tcp", s.listenAddress)
	if err != nil {
		return fmt.Errorf("listen error: %s", err.Error())
	}

	if s.isLeader {
		s.followers = make(map[*client.Client]struct{})
	} else {
		if len(s.leaderAddress) == 0 {
			return fmt.Errorf("leader address is empty")
		}

		if err := s.dialLeader(); err != nil {
			return err
		}
	}

	logger.CustomLogger.Info.Printf("server is running on port [%s] is leader [%t]", s.listenAddress, s.isLeader)

	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.CustomLogger.Error.Printf("accept error: %s", err.Error())
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *ServerNode) dialLeader() error {
	conn, err := net.Dial("tcp", s.leaderAddress)
	if err != nil {
		return fmt.Errorf("failed to dial leader [%s]", s.leaderAddress)
	}

	logger.CustomLogger.Info.Printf("connected to leader [%s]", s.leaderAddress)

	if err := binary.Write(conn, binary.LittleEndian, protocol.CmdJoin); err != nil {
		return err
	}

	s.leader = client.NewFromConn(conn)

	go s.handleConnection(conn)

	return nil
}

func (s *ServerNode) handleConnection(conn net.Conn) {
	defer conn.Close()

	logger.CustomLogger.Info.Printf("new connection made: %s", conn.RemoteAddr())

	for {
		cmd, err := protocol.ParseCommand(conn)
		if err != nil {
			if err != io.EOF {
				logger.CustomLogger.Error.Printf("parse command error: %s", err.Error())
			}

			break
		}

		go s.handleCommand(conn, cmd)
	}

	logger.CustomLogger.Info.Printf("connection closed: %s", conn.RemoteAddr())
}

func (s *ServerNode) handleCommand(conn net.Conn, cmd any) {
	switch v := cmd.(type) {
	case *protocol.CommandSet:
		s.handleSetCommand(conn, v)
	case *protocol.CommandGet:
		s.handleGetCommand(conn, v)
	case *protocol.CommandDelete:
		s.handleDeleteCommand(conn, v)
	case *protocol.CommandJoin:
		s.handleJoinCommand(conn, v)
	}
}

func (s *ServerNode) handleSetCommand(conn net.Conn, cmd *protocol.CommandSet) {
	msg := fmt.Sprintf("SET %s to %s", cmd.Key, cmd.Value)

	logger.CustomLogger.Info.Println(msg)

	var response protocol.ResponseSet

	defer func() {
		b, err := response.Bytes()
		if err != nil {
			logger.CustomLogger.Error.Printf("error sending response to %s while handling SET command error: %s", conn.RemoteAddr(), err.Error())
			return
		}

		if err := s.respond(conn, b); err != nil {
			logger.CustomLogger.Error.Printf("esending response to %s while handling SET command error: %s", conn.RemoteAddr(), err.Error())
			return
		}
	}()

	if err := s.cache.Set(cache.Key(cmd.Key), cache.Value{
		Value: cmd.Value,
		TTL:   time.Second * time.Duration(cmd.TTL),
	}); err != nil {
		response.Status = protocol.StatusError
		logger.CustomLogger.Info.Printf("handling SET command error: %s", err.Error())
		return
	}

	response.Status = protocol.StatusOK

	go func() {
		for follower := range s.followers {
			if err := follower.Set(context.Background(), cmd.Key, cmd.Value, cmd.TTL); err != nil {
				logger.CustomLogger.Error.Printf("forward to member [%s] error [%s]", follower, err.Error())
				continue
			}

			logger.CustomLogger.Info.Printf("forwarded message [%s] to member [%s]", msg, follower)
		}
	}()
}

func (s *ServerNode) handleGetCommand(conn net.Conn, cmd *protocol.CommandGet) {
	logger.CustomLogger.Info.Printf("GET %s", cmd.Key)

	var response protocol.ResponseGet

	defer func() {
		b, err := response.Bytes()
		if err != nil {
			logger.CustomLogger.Error.Printf("error sending response to %s while handling GET command error: %s", conn.RemoteAddr(), err.Error())
			return
		}

		if err := s.respond(conn, b); err != nil {
			logger.CustomLogger.Error.Printf("error sending response to %s while handling GET command error: %s", conn.RemoteAddr(), err.Error())
			return
		}
	}()

	val, err := s.cache.Get(cache.Key(cmd.Key))
	if err != nil {
		response.Status = protocol.StatusKeyNotFound
		logger.CustomLogger.Error.Printf("handling GET command error: %s", err.Error())
		return
	}

	response.Status = protocol.StatusOK
	response.Value = val.Value
}

func (s *ServerNode) handleDeleteCommand(conn net.Conn, cmd *protocol.CommandDelete) {
	msg := fmt.Sprintf("DELETE %s", cmd.Key)

	logger.CustomLogger.Info.Println(msg)

	resp := protocol.ResponseDelete{}

	defer func() {
		b, err := resp.Bytes()
		if err != nil {
			logger.CustomLogger.Error.Printf("error sending response to %s while handling DELETE command error: %s", conn.RemoteAddr(), err.Error())
			return
		}

		if err := s.respond(conn, b); err != nil {
			logger.CustomLogger.Error.Printf("error sending response to %s while handling DELETE command error: %s", conn.RemoteAddr(), err.Error())
			return
		}
	}()

	if err := s.cache.Delete(cache.Key(cmd.Key)); err != nil {
		resp.Status = protocol.StatusKeyNotFound
		logger.CustomLogger.Error.Printf("handling DELETE command error: %s", err.Error())
		return
	}

	resp.Status = protocol.StatusOK

	go func() {
		for follower := range s.followers {
			if err := follower.Delete(context.Background(), cmd.Key); err != nil {
				logger.CustomLogger.Error.Printf("forward to member [%s] error [%s]", follower, err.Error())
				continue
			}

			logger.CustomLogger.Info.Printf("forwarded message [%s] to member [%s]", msg, follower)
		}
	}()
}

func (s *ServerNode) handleJoinCommand(conn net.Conn, cmd *protocol.CommandJoin) {
	logger.CustomLogger.Info.Printf("member just joined the cluster [%s]", conn.RemoteAddr())

	s.followers[client.NewFromConn(conn)] = struct{}{}
}

func (s *ServerNode) respond(conn net.Conn, msg []byte) error {
	_, err := conn.Write(msg)
	return err
}
