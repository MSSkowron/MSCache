package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/MSSkowron/mscache/cache"
	"github.com/MSSkowron/mscache/protocol"
)

const bufferSize = 2048

type Server struct {
	listenAddr string
	leaderAddr string
	isLeader   bool
	cache      cache.Cacher
}

func New(listenAddr, leaderAddr string, isLeader bool, c cache.Cacher) *Server {
	server := &Server{
		listenAddr: listenAddr,
		leaderAddr: leaderAddr,
		isLeader:   isLeader,
		cache:      c,
	}

	return server
}

func (s *Server) Run() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return fmt.Errorf("[Server] Listen error: %s\n", err.Error())
	}

	log.Printf("[Server] Server is running on port [%s]\n", s.listenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("[Server] Accept error: %s\n", err.Error())
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	for {
		cmd, err := protocol.ParseCommand(conn)
		if err != nil {
			if err != io.EOF {
				log.Printf("[Server] Parse command error: %s\n", err.Error())
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
	default:
		log.Println("[Server] Invalid command type")
	}
}

func (s *Server) handleSetCommand(conn net.Conn, cmd *protocol.CommandSet) error {
	if err := s.cache.Set(cmd.Key, cmd.Value, time.Duration(cmd.TTL)); err != nil {
		return err
	}

	return nil
}

func (s *Server) handleGetCommand(conn net.Conn, cmd *protocol.CommandGet) error {
	val, err := s.cache.Get(cmd.Key)
	if err != nil {
		return err
	}

	_, err = conn.Write(val)

	return err
}
