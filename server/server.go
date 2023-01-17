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

	log.Printf("[Server] New connection made: %s", conn.RemoteAddr())

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

func (s *Server) handleSetCommand(conn net.Conn, cmd *protocol.CommandSet) {
	var status protocol.Status

	defer func() {
		response := protocol.ResponseSet{
			Status: status,
		}

		b, err := response.Bytes()
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
		status = protocol.StatusError
		log.Printf("[Server] Handling SET command error: %s\n", err.Error())
		return
	}

	status = protocol.StatusOK
}

func (s *Server) handleGetCommand(conn net.Conn, cmd *protocol.CommandGet) {
	var (
		status protocol.Status
		value  []byte
	)

	defer func() {
		response := protocol.ResponseGet{
			Status: status,
			Value:  value,
		}

		b, err := response.Bytes()
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
		status = protocol.StatusError
		log.Printf("[Server] Handling GET command error: %s\n", err.Error())
		return
	}

	status = protocol.StatusOK
	value = val
}

func (s *Server) respond(conn net.Conn, msg []byte) error {
	_, err := conn.Write(msg)
	return err
}
