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

type Server struct {
	listenAddr string
	isLeader   bool
	cache      cache.Cacher
}

func NewServer(listenAddr string, isLeader bool, c cache.Cacher) *Server {
	return &Server{
		listenAddr: listenAddr,
		isLeader:   isLeader,
		cache:      c,
	}
}

func (s *Server) Run() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return fmt.Errorf("Listen error: %s\n", err.Error())
	}

	log.Printf("[Server] Server is running on port [%s]\n", s.listenAddr)

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

		go s.handleCommand(conn, buf[:n])
	}
}

func (s *Server) handleCommand(conn net.Conn, rawCMD []byte) {
	var err error

	defer func() {
		if err != nil {
			log.Printf("Failed to handle command: %s\n", err.Error())
			_, err = conn.Write([]byte(err.Error()))
			if err != nil {
				log.Printf("Failed to respond: %s\n", err.Error())
			}
		}
	}()

	msg, err := parseCommand(rawCMD)
	if err != nil {
		return
	}

	switch msg.Cmd {
	case CMDSet:
		err = s.handleSetCommand(conn, msg)
	case CMDGet:
		err = s.handleGetCommand(conn, msg)
	}
}

func (s *Server) handleSetCommand(conn net.Conn, msg *Message) error {
	if err := s.cache.Set(msg.Key, msg.Value, msg.TTL); err != nil {
		return err
	}

	go s.sendToFollowers(context.TODO(), msg)

	return nil
}

func (s *Server) handleGetCommand(conn net.Conn, msg *Message) error {
	value, err := s.cache.Get(msg.Key)
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte(value))

	go s.sendToFollowers(context.TODO(), msg)

	return err
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
