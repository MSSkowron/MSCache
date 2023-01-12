package server

import (
	"fmt"
	"time"
)

type Command string

const (
	CMDSet Command = "SET"
	CMDGet Command = "GET"
)

type Message struct {
	Cmd   Command
	Key   []byte
	Value []byte
	TTL   time.Duration
}

func (m *Message) ToBytes() []byte {
	switch m.Cmd {
	case CMDSet:
		return []byte(fmt.Sprintf("%s %s %s %d", m.Cmd, m.Key, m.Value, m.TTL))
	case CMDGet:
		return []byte(fmt.Sprintf("%s %s", m.Cmd, m.Key))
	default:
		panic("unknow command")
	}
}
