package server

import "time"

type Command string

const (
	CMDSet Command = "SET"
)

type Message struct {
	Cmd   Command
	Key   []byte
	Value []byte
	TTL   time.Duration
}
