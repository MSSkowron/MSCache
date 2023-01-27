package server

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/MSSkowron/mscache/cache"
	"github.com/MSSkowron/mscache/client"
	"github.com/MSSkowron/mscache/logger"
	"github.com/MSSkowron/mscache/protocol"
)

type Server struct {
	listenAddr string
	leaderAddr string
	isLeader   bool
	members    map[*client.Client]struct{}
	cache      cache.Cacher
}

func New(listenAddr, leaderAddr string, isLeader bool, c cache.Cacher) *Server {
	return &Server{
		listenAddr: listenAddr,
		leaderAddr: leaderAddr,
		isLeader:   isLeader,
		cache:      c,
		members:    make(map[*client.Client]struct{}),
	}
}

func (s *Server) Run() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err.Error())
	}

	if !s.isLeader && len(s.leaderAddr) != 0 {
		go func() {
			if err := s.dialLeader(); err != nil {
				log.Println(err)
			}
		}()
	}

	logger.InfoLogger.Printf("server is running on port [%s] is leader [%t]", s.listenAddr, s.isLeader)

	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.ErrorLogger.Printf("accept error: %s", err.Error())
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) dialLeader() error {
	conn, err := net.Dial("tcp", s.leaderAddr)
	if err != nil {
		return fmt.Errorf("failed to dial leader [%s]", s.leaderAddr)
	}

	logger.InfoLogger.Printf("connected to leader [%s]", s.leaderAddr)

	if err := binary.Write(conn, binary.LittleEndian, protocol.CmdJoin); err != nil {
		return err
	}

	s.handleConn(conn)

	return nil
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	logger.InfoLogger.Printf("new connection made: %s", conn.RemoteAddr())

	for {
		cmd, err := protocol.ParseCommand(conn)
		if err != nil {
			if err != io.EOF {
				logger.ErrorLogger.Printf("parse command error: %s", err.Error())
			}

			break
		}

		go s.handleCommand(conn, cmd)
	}

	logger.InfoLogger.Printf("connection closed: %s", conn.RemoteAddr())
}

func (s *Server) handleCommand(conn net.Conn, cmd any) {
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

func (s *Server) handleSetCommand(conn net.Conn, cmd *protocol.CommandSet) {
	msg := fmt.Sprintf("SET %s to %s", cmd.Key, cmd.Value)

	logger.InfoLogger.Println(msg)

	go func() {
		for member := range s.members {
			if err := member.Set(context.Background(), cmd.Key, cmd.Value, cmd.TTL); err != nil {
				logger.ErrorLogger.Printf("forward to member [%s] error [%s]", member, err.Error())
				continue
			}

			logger.InfoLogger.Printf("forwarded message [%s] to member [%s]", msg, member)
		}
	}()

	resp := protocol.ResponseSet{}

	defer func() {
		b, err := resp.Bytes()
		if err != nil {
			logger.ErrorLogger.Printf("error sending response to %s while handling SET command error: %s", conn.RemoteAddr(), err.Error())
			return
		}

		if err := s.respond(conn, b); err != nil {
			logger.ErrorLogger.Printf("esending response to %s while handling SET command error: %s", conn.RemoteAddr(), err.Error())
			return
		}
	}()

	if err := s.cache.Set(cmd.Key, cmd.Value, time.Duration(cmd.TTL)); err != nil {
		resp.Status = protocol.StatusError
		logger.InfoLogger.Printf("handling SET command error: %s", err.Error())
		return
	}

	resp.Status = protocol.StatusOK
}

func (s *Server) handleGetCommand(conn net.Conn, cmd *protocol.CommandGet) {
	logger.InfoLogger.Printf("GET %s", cmd.Key)

	resp := protocol.ResponseGet{}

	defer func() {
		b, err := resp.Bytes()
		if err != nil {
			logger.ErrorLogger.Printf("error sending response to %s while handling GET command error: %s", conn.RemoteAddr(), err.Error())
			return
		}

		if err := s.respond(conn, b); err != nil {
			logger.ErrorLogger.Printf("error sending response to %s while handling GET command error: %s", conn.RemoteAddr(), err.Error())
			return
		}
	}()

	val, err := s.cache.Get(cmd.Key)
	if err != nil {
		resp.Status = protocol.StatusKeyNotFound
		logger.ErrorLogger.Printf("handling GET command error: %s", err.Error())
		return
	}

	resp.Status = protocol.StatusOK
	resp.Value = val
}

func (s *Server) handleDeleteCommand(conn net.Conn, cmd *protocol.CommandDelete) {
	msg := fmt.Sprintf("DELETE %s", cmd.Key)

	logger.InfoLogger.Println(msg)

	resp := protocol.ResponseDelete{}

	go func() {
		for member := range s.members {
			if err := member.Delete(context.Background(), cmd.Key); err != nil {
				logger.ErrorLogger.Printf("forward to member [%s] error [%s]", member, err.Error())
				continue
			}

			logger.InfoLogger.Printf("forwarded message [%s] to member [%s]", msg, member)
		}
	}()

	defer func() {
		b, err := resp.Bytes()
		if err != nil {
			logger.ErrorLogger.Printf("error sending response to %s while handling DELETE command error: %s", conn.RemoteAddr(), err.Error())
			return
		}

		if err := s.respond(conn, b); err != nil {
			logger.ErrorLogger.Printf("error sending response to %s while handling DELETE command error: %s", conn.RemoteAddr(), err.Error())
			return
		}
	}()

	if err := s.cache.Delete(cmd.Key); err != nil {
		resp.Status = protocol.StatusKeyNotFound
		logger.ErrorLogger.Printf("handling DELETE command error: %s", err.Error())
		return
	}

	resp.Status = protocol.StatusOK
}

func (s *Server) handleJoinCommand(conn net.Conn, cmd *protocol.CommandJoin) {
	logger.InfoLogger.Printf("member just joined the cluster [%s]", conn.RemoteAddr())

	s.members[client.NewFromConn(conn)] = struct{}{}
}

func (s *Server) respond(conn net.Conn, msg []byte) error {
	_, err := conn.Write(msg)
	return err
}
