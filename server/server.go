package server

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/MSSkowron/MSCache/server/cache"
	"github.com/MSSkowron/MSCache/server/logger"
	"github.com/MSSkowron/MSCache/server/protocol"
)

var (
	ErrEmptyLeaderAddress = errors.New("leader address is empty")
)

type ServerNode struct {
	listener      net.Listener
	listenAddress string
	leaderAddress string
	isLeader      bool
	followers     map[net.Conn]struct{}
	leader        net.Conn
	cache         cache.Cache
}

func New(listenAddress, leaderAddress string, isLeader bool, c cache.Cache) *ServerNode {
	return &ServerNode{
		listenAddress: listenAddress,
		leaderAddress: leaderAddress,
		isLeader:      isLeader,
		cache:         c,
	}
}

func (s *ServerNode) Run() error {
	ln, err := net.Listen("tcp", s.listenAddress)
	if err != nil {
		return fmt.Errorf("running tcp listener error: %s", err.Error())
	}
	defer func() {
		_ = ln.Close()
	}()

	s.listener = ln

	if s.isLeader {
		s.followers = make(map[net.Conn]struct{})
	} else {
		if len(s.leaderAddress) == 0 {
			return ErrEmptyLeaderAddress
		}

		if err := s.dialLeader(); err != nil {
			return fmt.Errorf("connecting to leader %s error: %s", s.leaderAddress, err.Error())
		}

		logger.CustomLogger.Info.Printf("Connected to leader [%s]", s.leaderAddress)
	}

	logger.CustomLogger.Info.Printf("Server is running on [%s], is leader [%t]", s.listenAddress, s.isLeader)

	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.CustomLogger.Error.Printf("accept a new connection error: %s", err.Error())
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *ServerNode) Close() error {
	return s.listener.Close()
}

func (s *ServerNode) dialLeader() error {
	conn, err := net.Dial("tcp", s.leaderAddress)
	if err != nil {
		return err
	}

	if err := binary.Write(conn, binary.LittleEndian, protocol.CmdJoin); err != nil {
		return err
	}

	s.leader = conn

	go s.handleConnection(conn)

	return nil
}

func (s *ServerNode) handleConnection(conn net.Conn) {
	logger.CustomLogger.Info.Printf("Opened connection with [%s]", conn.RemoteAddr())

	defer func() {
		_ = conn.Close()

		if s.isLeader {
			delete(s.followers, conn)
		}
	}()

	for {
		cmd, err := protocol.ParseCommand(conn)
		if err != nil {
			break
		}

		go s.handleCommand(conn, cmd)
	}

	logger.CustomLogger.Info.Printf("Closed connection with [%s]", conn.RemoteAddr())

	if !s.isLeader && conn.RemoteAddr() == s.leader.RemoteAddr() {
		logger.CustomLogger.Error.Printf("Lost connection with leader [%s]", s.leader.RemoteAddr())

		_ = s.Close()
		os.Exit(0)
	}
}

func (s *ServerNode) handleCommand(conn net.Conn, cmd any) {
	switch v := cmd.(type) {
	case *protocol.CommandGet:
		s.handleGetCommand(conn, v)
	case *protocol.CommandSet:
		if s.isLeader || (!s.isLeader && conn.RemoteAddr() == s.leader.RemoteAddr()) {
			s.handleSetCommand(conn, v)
			return
		}

		response := protocol.ResponseSet{
			Status: protocol.StatusNotLeader,
		}

		b, err := response.Bytes()
		if err != nil {
			logger.CustomLogger.Error.Printf("sending response to %s while handling GET command error: %s", conn.RemoteAddr(), err.Error())
			return
		}

		if err := s.respond(conn, b); err != nil {
			logger.CustomLogger.Error.Printf("sending response to %s while handling GET command error: %s", conn.RemoteAddr(), err.Error())
			return
		}

		return
	case *protocol.CommandDelete:
		if s.isLeader || (!s.isLeader && conn.RemoteAddr() == s.leader.RemoteAddr()) {
			s.handleDeleteCommand(conn, v)
			return
		}

		response := protocol.ResponseDelete{
			Status: protocol.StatusNotLeader,
		}

		b, err := response.Bytes()
		if err != nil {
			logger.CustomLogger.Error.Printf("sending response to %s while handling GET command error: %s", conn.RemoteAddr(), err.Error())
			return
		}

		if err := s.respond(conn, b); err != nil {
			logger.CustomLogger.Error.Printf("sending response to %s while handling GET command error: %s", conn.RemoteAddr(), err.Error())
			return
		}

		return
	case *protocol.CommandJoin:
		s.handleJoinCommand(conn, v)
	}
}

func (s *ServerNode) handleGetCommand(conn net.Conn, cmd *protocol.CommandGet) {
	var (
		key      = cache.Key(cmd.Key)
		response protocol.ResponseGet
	)

	logger.CustomLogger.Info.Printf("Received GET key=%s from [%s]", key, conn.RemoteAddr())

	defer func() {
		b, err := response.Bytes()
		if err != nil {
			logger.CustomLogger.Error.Printf("sending response to %s while handling GET command error: %s", conn.RemoteAddr(), err.Error())
			return
		}

		if err := s.respond(conn, b); err != nil {
			logger.CustomLogger.Error.Printf("sending response to %s while handling GET command error: %s", conn.RemoteAddr(), err.Error())
			return
		}
	}()

	val, err := s.cache.Get(key)
	if err != nil {
		logger.CustomLogger.Error.Printf("getting key %s from cache error: %s", key, err.Error())
		response.Status = protocol.StatusKeyNotFound
		return
	}

	response.Status = protocol.StatusOK
	response.Value = val.Value
}

func (s *ServerNode) handleSetCommand(conn net.Conn, cmd *protocol.CommandSet) {
	var (
		key   = cache.Key(cmd.Key)
		value = string(cache.Value{
			Value: cmd.Value,
		}.Value)
		response protocol.ResponseSet
	)

	logger.CustomLogger.Info.Printf("Received SET key=%s value=%s from [%s]", key, value, conn.RemoteAddr())

	defer func() {
		b, err := response.Bytes()
		if err != nil {
			logger.CustomLogger.Error.Printf("sending response to %s while handling SET command error: %s", conn.RemoteAddr(), err.Error())
			return
		}

		if err := s.respond(conn, b); err != nil {
			logger.CustomLogger.Error.Printf("sending response to %s while handling SET command error: %s", conn.RemoteAddr(), err.Error())
			return
		}
	}()

	if err := s.cache.Set(key, cache.Value{
		Value: cmd.Value,
		TTL:   time.Second * time.Duration(cmd.TTL),
	}); err != nil {
		logger.CustomLogger.Error.Printf("setting key %s to value %s in cache error: %s", key, value, err.Error())
		response.Status = protocol.StatusError
		return
	}

	response.Status = protocol.StatusOK

	if s.isLeader {
		go func() {
			propagateSetCmd := &protocol.CommandSet{
				Key:   cmd.Key,
				Value: cmd.Value,
				TTL:   cmd.TTL,
			}

			b, err := propagateSetCmd.Bytes()
			if err != nil {
				logger.CustomLogger.Error.Printf("propagating SET command error: %s", err.Error())
				return
			}

			for follower := range s.followers {
				_, err = follower.Write(b)
				if err != nil {
					logger.CustomLogger.Error.Printf("propagating SET command to member %s error: %s", follower.RemoteAddr(), err.Error())
				}
			}
		}()
	}
}

func (s *ServerNode) handleDeleteCommand(conn net.Conn, cmd *protocol.CommandDelete) {
	var (
		key      = cache.Key(cmd.Key)
		response protocol.ResponseDelete
	)

	logger.CustomLogger.Info.Printf("Received DELETE key=%s from [%s]", key, conn.RemoteAddr())

	defer func() {
		b, err := response.Bytes()
		if err != nil {
			logger.CustomLogger.Error.Printf("sending response to %s while handling DELETE command error: %s", conn.RemoteAddr(), err.Error())
			return
		}

		if err := s.respond(conn, b); err != nil {
			logger.CustomLogger.Error.Printf("sending response to %s while handling DELETE command error: %s", conn.RemoteAddr(), err.Error())
			return
		}
	}()

	if err := s.cache.Delete(cache.Key(cmd.Key)); err != nil {
		logger.CustomLogger.Error.Printf("deleting key %s from cache error: %s", key, err.Error())
		response.Status = protocol.StatusKeyNotFound
		return
	}

	response.Status = protocol.StatusOK

	if s.isLeader {
		go func() {
			propagateDelCmd := &protocol.CommandDelete{
				Key: cmd.Key,
			}

			b, err := propagateDelCmd.Bytes()
			if err != nil {
				logger.CustomLogger.Error.Printf("propagating DELETE command error: %s", err.Error())
				return
			}

			for follower := range s.followers {
				_, err = follower.Write(b)
				if err != nil {
					logger.CustomLogger.Error.Printf("propagating DELETE command to member %s error: %s", follower, err.Error())
				}
			}
		}()
	}
}

func (s *ServerNode) handleJoinCommand(conn net.Conn, cmd *protocol.CommandJoin) {
	logger.CustomLogger.Info.Printf("New member [%s] joined the cluster", conn.RemoteAddr())

	s.followers[conn] = struct{}{}
}

func (s *ServerNode) respond(conn net.Conn, msg []byte) error {
	_, err := conn.Write(msg)
	return err
}
