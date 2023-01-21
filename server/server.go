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
	"github.com/MSSkowron/mscache/protocol"
)

type Server struct {
	listenAddr string
	leaderAddr string
	isLeader   bool
	cache      cache.Cacher
	members    map[*client.Client]struct{}
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
		return fmt.Errorf("listen error: %s\n", err.Error())
	}

	if !s.isLeader && len(s.leaderAddr) != 0 {
		go func() {
			if err := s.dialLeader(); err != nil {
				log.Println(err)
			}
		}()
	}

	log.Printf("[Server] Server is running on port [%s] is leader [%t]\n", s.listenAddr, s.isLeader)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %s\n", err.Error())
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

	log.Printf("[Server] Connected to leader [%s]\n", s.leaderAddr)

	if err := binary.Write(conn, binary.LittleEndian, protocol.CmdJoin); err != nil {
		return err
	}

	s.handleConn(conn)

	return nil
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	log.Printf("[Server] New connection made: %s", conn.RemoteAddr())

	for {
		cmd, err := protocol.ParseCommand(conn)
		if err != nil {
			if err != io.EOF {
				log.Printf("parse command error: %s\n", err.Error())
			}

			break
		}

		go s.handleCommand(conn, cmd)
	}

	log.Printf("[Server] Connection closed: %s", conn.RemoteAddr())
}

func (s *Server) handleCommand(conn net.Conn, cmd any) {
	switch v := cmd.(type) {
	case *protocol.CommandSet:
		s.handleSetCommand(conn, v)
	case *protocol.CommandGet:
		s.handleGetCommand(conn, v)
	case *protocol.CommandJoin:
		s.handleJoinCommand(conn, v)
	}
}

func (s *Server) handleSetCommand(conn net.Conn, cmd *protocol.CommandSet) {
	log.Printf("[Server] SET %s to %s\n", cmd.Key, cmd.Value)

	go func() {
		for member := range s.members {
			if err := member.Set(context.Background(), cmd.Key, cmd.Value, cmd.TTL); err != nil {
				log.Println("[Server] Forward to member error:", err)
			}

			log.Printf("[Server] Forwarding message to members")
		}
	}()

	resp := protocol.ResponseSet{}

	defer func() {
		b, err := resp.Bytes()
		if err != nil {
			log.Printf("[Server] Error sending response to %s while handling SET command error: %s\n", conn.RemoteAddr(), err.Error())
			return
		}

		if err := s.respond(conn, b); err != nil {
			log.Printf("[Server] Error sending response to %s while handling SET command error: %s\n", conn.RemoteAddr(), err.Error())
			return
		}
	}()

	if err := s.cache.Set(cmd.Key, cmd.Value, time.Duration(cmd.TTL)); err != nil {
		resp.Status = protocol.StatusError
		log.Printf("[Server] Handling SET command error: %s\n", err.Error())
		return
	}

	resp.Status = protocol.StatusOK
}

func (s *Server) handleGetCommand(conn net.Conn, cmd *protocol.CommandGet) {
	log.Printf("[Server] GET %s\n", cmd.Key)

	resp := protocol.ResponseGet{}

	defer func() {
		b, err := resp.Bytes()
		if err != nil {
			log.Printf("[Server] Error sending response to %s while handling GET command error: %s\n", conn.RemoteAddr(), err.Error())
			return
		}

		if err := s.respond(conn, b); err != nil {
			log.Printf("[Server] Error sending response to %s while handling GET command error: %s\n", conn.RemoteAddr(), err.Error())
			return
		}
	}()

	val, err := s.cache.Get(cmd.Key)
	if err != nil {
		resp.Status = protocol.StatusKeyNotFound
		log.Printf("[Server] Handling GET command error: %s\n", err.Error())
		return
	}

	resp.Status = protocol.StatusOK
	resp.Value = val
}

func (s *Server) handleJoinCommand(conn net.Conn, cmd *protocol.CommandJoin) {
	log.Printf("[Server] Member just joined the cluster [%s]\n", conn.RemoteAddr())

	s.members[client.NewFromConn(conn)] = struct{}{}
}

func (s *Server) respond(conn net.Conn, msg []byte) error {
	_, err := conn.Write(msg)
	return err
}
