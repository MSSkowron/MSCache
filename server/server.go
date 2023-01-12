package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

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
		log.Println("Got message:", string(msg))

		go s.handleCommand(conn, buf[:n])
	}
}

func (s *Server) handleCommand(conn net.Conn, rawCMD []byte) {
	msg, err := parseCommand(rawCMD)
	if err != nil {
		log.Printf("Failed to parse command: %s\n", err.Error())
		//respond
		return
	}

	switch msg.Cmd {
	case CMDSet:
		if err := s.handleSetCommand(conn, msg); err != nil {
			log.Printf("Failed to handle set command: %s\n", err.Error())
			//respond
			return
		}
	}
}

func (s *Server) handleSetCommand(conn net.Conn, msg *Message) error {
	if err := s.cache.Set(msg.Key, msg.Value, msg.TTL); err != nil {
		return err
	}

	go s.sendToFollowers(context.TODO(), msg)

	return nil
}

func (s *Server) sendToFollowers(ctx context.Context, msg *Message) error {

	return nil
}

func parseCommand(raw []byte) (*Message, error) {
	rawStr := string(raw)

	parts := strings.Split(rawStr, " ")
	if len(parts) < 2 {
		return nil, errors.New("invalid protocol format")
	}

	msg := &Message{
		Cmd: Command(parts[0]),
		Key: []byte(parts[1]),
	}

	if msg.Cmd == CMDSet {
		if len(parts) < 4 {
			return nil, errors.New("invalid SET command")
		}

		msg.Value = []byte(parts[2])

		ttl, err := strconv.Atoi(parts[3])
		if err != nil {
			return nil, errors.New("invalid TTL in SET command")
		}
		msg.TTL = time.Duration(ttl)
	}

	return msg, nil
}
