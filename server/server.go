package server

import (
	"fmt"
	"log"
	"net"

	"github.com/MSSkowron/mscache/cache"
)

const bufferSize = 2048

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
}

type Server struct {
	ServerOpts

	cache cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      c,
	}
}

func (s *Server) Run() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("Listen error: %s\n", err.Error())
	}

	log.Printf("Server is running on port [%s]\n", s.ListenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Accept error: %s\n", err.Error())
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, bufferSize)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("Conn read error: %s\n", err.Error())
			break
		}

		msg := buf[:n]
		fmt.Println(string(msg))
	}
}
