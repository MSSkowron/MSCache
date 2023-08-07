package node

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/MSSkowron/MSCache/internal/cache"
	"github.com/MSSkowron/MSCache/internal/protocol"
	"github.com/MSSkowron/MSCache/pkg/logger"
)

var (
	// ErrEmptyLeaderAddress is returned when leader address is empty.
	ErrEmptyLeaderAddress = errors.New("leader address is empty")
)

// Node represents a server node.
type Node struct {
	listener      net.Listener
	listenAddress string
	leaderAddress string
	isLeader      bool
	followers     map[net.Conn]struct{}
	leader        net.Conn
	cache         cache.Cache
}

// New creates a new Node Node.
func New(listenAddress, leaderAddress string, isLeader bool, c cache.Cache) *Node {
	return &Node{
		listenAddress: listenAddress,
		leaderAddress: leaderAddress,
		isLeader:      isLeader,
		cache:         c,
	}
}

// Run runs the Node Node.
func (s *Node) Run() error {
	ln, err := net.Listen("tcp", s.listenAddress)
	if err != nil {
		return fmt.Errorf("running tcp listener: %s", err)
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
			return fmt.Errorf("connecting to leader %s: %s", s.leaderAddress, err)
		}

		logger.Infof("Connected to leader %s", s.leaderAddress)
	}

	logger.Infof("Node is running on %s, is leader: %t", s.listenAddress, s.isLeader)

	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.Errorf("accepting a new connection: %s", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

// Close closes the Node Node.
func (s *Node) Close() error {
	return s.listener.Close()
}

func (s *Node) dialLeader() error {
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

func (s *Node) handleConnection(conn net.Conn) {
	logger.Infof("Opened connection with %s", conn.RemoteAddr())

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

	logger.Infof("Closed connection with %s", conn.RemoteAddr())

	if !s.isLeader && conn.RemoteAddr() == s.leader.RemoteAddr() {
		logger.Errorf("Lost connection with leader %s", s.leader.RemoteAddr())

		_ = s.Close()
		os.Exit(0)
	}
}

func (s *Node) handleCommand(conn net.Conn, cmd any) {
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
			logger.Errorf("responding to %s while handling GET command: %s", conn.RemoteAddr(), err)
			return
		}

		if err := s.respond(conn, b); err != nil {
			logger.Errorf("responding to %s while handling GET command: %s", conn.RemoteAddr(), err)
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
			logger.Errorf("responding to %s while handling GET command: %s", conn.RemoteAddr(), err)
			return
		}

		if err := s.respond(conn, b); err != nil {
			logger.Errorf("responding to %s while handling GET command: %s", conn.RemoteAddr(), err)
			return
		}

		return
	case *protocol.CommandJoin:
		s.handleJoinCommand(conn, v)
	}
}

func (s *Node) handleGetCommand(conn net.Conn, cmd *protocol.CommandGet) {
	var (
		key      = cache.Key(cmd.Key)
		response protocol.ResponseGet
	)

	logger.Infof("Received GET key=%s from %s", key, conn.RemoteAddr())

	defer func() {
		b, err := response.Bytes()
		if err != nil {
			logger.Errorf("responding to %s while handling GET command: %s", conn.RemoteAddr(), err)
			return
		}

		if err := s.respond(conn, b); err != nil {
			logger.Errorf("responding to %s while handling GET command: %s", conn.RemoteAddr(), err)
			return
		}
	}()

	val, err := s.cache.Get(key)
	if err != nil {
		if errors.Is(err, cache.ErrKeyNotFound) {
			response.Status = protocol.StatusKeyNotFound
			return
		}

		logger.Errorf("getting key %s from cache: %s", key, err)
		response.Status = protocol.StatusError
		return
	}

	response.Status = protocol.StatusOK
	response.Value = val.Value
}

func (s *Node) handleSetCommand(conn net.Conn, cmd *protocol.CommandSet) {
	var (
		key   = cache.Key(cmd.Key)
		value = string(cache.Value{
			Value: cmd.Value,
		}.Value)
		response protocol.ResponseSet
	)

	logger.Infof("Received SET key=%s value=%s from %s", key, value, conn.RemoteAddr())

	defer func() {
		b, err := response.Bytes()
		if err != nil {
			logger.Errorf("responding to %s while handling SET command: %s", conn.RemoteAddr(), err)
			return
		}

		if err := s.respond(conn, b); err != nil {
			logger.Errorf("responding to %s while handling SET command: %s", conn.RemoteAddr(), err)
			return
		}
	}()

	if err := s.cache.Set(key, cache.Value{
		Value: cmd.Value,
		TTL:   time.Second * time.Duration(cmd.TTL),
	}); err != nil {
		logger.Errorf("setting key %s to value %s in cache: %s", key, value, err)
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
				logger.Errorf("propagating SET command: %s", err)
				return
			}

			for follower := range s.followers {
				_, err = follower.Write(b)
				if err != nil {
					logger.Errorf("propagating SET command to member %s: %s", follower.RemoteAddr(), err)
				}
			}
		}()
	}
}

func (s *Node) handleDeleteCommand(conn net.Conn, cmd *protocol.CommandDelete) {
	var (
		key      = cache.Key(cmd.Key)
		response protocol.ResponseDelete
	)

	logger.Infof("Received DELETE key=%s from %s", key, conn.RemoteAddr())

	defer func() {
		b, err := response.Bytes()
		if err != nil {
			logger.Errorf("responding to %s while handling DELETE command: %s", conn.RemoteAddr(), err)
			return
		}

		if err := s.respond(conn, b); err != nil {
			logger.Errorf("responding to %s while handling DELETE command: %s", conn.RemoteAddr(), err)
			return
		}
	}()

	if err := s.cache.Delete(cache.Key(cmd.Key)); err != nil {
		logger.Errorf("deleting key %s from cache: %s", key, err)
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
				logger.Errorf("propagating DELETE command: %s", err)
				return
			}

			for follower := range s.followers {
				_, err = follower.Write(b)
				if err != nil {
					logger.Errorf("propagating DELETE command to member %s: %s", follower, err)
				}
			}
		}()
	}
}

func (s *Node) handleJoinCommand(conn net.Conn, cmd *protocol.CommandJoin) {
	logger.Infof("New member %s joined the cluster", conn.RemoteAddr())

	s.followers[conn] = struct{}{}
}

func (s *Node) respond(conn net.Conn, msg []byte) error {
	_, err := conn.Write(msg)
	return err
}
