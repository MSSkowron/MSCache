package server

import "time"

type Command string

const (
	CMDSet Command = "SET"
	CMDGet Command = "GET"
)

type Message struct {
	Cmd Command
	Key     []byte
	Value   []byte
	TTL     time.Duration
}
