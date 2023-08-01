package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// Command represents the different types of commands.
type Command byte

const (
	// CmdNone represents an empty command.
	CmdNone Command = iota
	// CmdSet represents the Set command.
	CmdSet
	// CmdGet represents the Get command.
	CmdGet
	// CmdDel represents the Delete command.
	CmdDel
	// CmdJoin represents the Join command.
	CmdJoin
)

// Status represents the different status types for responses.
type Status byte

const (
	// StatusNone represents an empty status.
	StatusNone Status = iota
	// StatusOK represents a successful status.
	StatusOK
	// StatusError represents an error status.
	StatusError
	// StatusKeyNotFound represents a key not found status.
	StatusKeyNotFound
	// StatusNotLeader represents a not leader status.
	StatusNotLeader
)

// ResponseSet represents response for Set command.
type ResponseSet struct {
	Status Status
}

// ResponseGet represents response for Get command.
type ResponseGet struct {
	Status Status
	Value  []byte
}

// ResponseDelete represents response for Delete command.
type ResponseDelete struct {
	Status Status
}

// CommandSet represents Set command.
type CommandSet struct {
	Key   []byte
	Value []byte
	TTL   int
}

// CommandGet represents Get command.
type CommandGet struct {
	Key []byte
}

// CommandDelete represents Delete command.
type CommandDelete struct {
	Key []byte
}

// CommandJoin represents Join command.
type CommandJoin struct{}

func (s Status) String() string {
	switch s {
	case StatusOK:
		return "OK"
	case StatusError:
		return "ERROR"
	case StatusKeyNotFound:
		return "NOT FOUND"
	case StatusNotLeader:
		return "NOT LEADER"
	default:
		return "NONE"
	}
}

// Bytes returns byte representation of response to set command.
func (r *ResponseSet) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, r.Status); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Bytes returns byte representation of response to get command.
func (r *ResponseGet) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, r.Status); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.LittleEndian, int32(len(r.Value))); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.LittleEndian, r.Value); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Bytes returns byte representation of response to delete command.
func (r *ResponseDelete) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, r.Status); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Bytes returns byte representation of join command.
func (c *CommandSet) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, CmdSet); err != nil {
		return nil, err
	}

	keyLen := int32(len(c.Key))
	if err := binary.Write(buf, binary.LittleEndian, keyLen); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, c.Key); err != nil {
		return nil, err
	}

	valueLen := int32(len(c.Value))
	if err := binary.Write(buf, binary.LittleEndian, valueLen); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, c.Value); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.LittleEndian, int32(c.TTL)); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Bytes returns byte representation of get command.
func (c *CommandGet) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, CmdGet); err != nil {
		return nil, err
	}

	keyLen := int32(len(c.Key))
	if err := binary.Write(buf, binary.LittleEndian, keyLen); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, c.Key); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Bytes returns byte representation of delete command.
func (c *CommandDelete) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, CmdDel); err != nil {
		return nil, err
	}

	keyLen := int32(len(c.Key))
	if err := binary.Write(buf, binary.LittleEndian, keyLen); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, c.Key); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ParseSetResponse parses response to set command.
func ParseSetResponse(r io.Reader) (*ResponseSet, error) {
	resp := &ResponseSet{}

	if err := binary.Read(r, binary.LittleEndian, &resp.Status); err != nil {
		return nil, err
	}

	return resp, nil
}

// ParseGetResponse parses response to get command.
func ParseGetResponse(r io.Reader) (*ResponseGet, error) {
	resp := &ResponseGet{}

	if err := binary.Read(r, binary.LittleEndian, &resp.Status); err != nil {
		return nil, err
	}

	var valueLen int32
	if err := binary.Read(r, binary.LittleEndian, &valueLen); err != nil {
		return nil, err
	}
	resp.Value = make([]byte, valueLen)
	if err := binary.Read(r, binary.LittleEndian, &resp.Value); err != nil {
		return nil, err
	}

	return resp, nil
}

// ParseDeleteResponse parses response to delete command.
func ParseDeleteResponse(r io.Reader) (*ResponseDelete, error) {
	resp := &ResponseDelete{}

	if err := binary.Read(r, binary.LittleEndian, &resp.Status); err != nil {
		return nil, err
	}

	return resp, nil
}

// ParseCommand parses command.
func ParseCommand(r io.Reader) (any, error) {
	var cmd Command
	if err := binary.Read(r, binary.LittleEndian, &cmd); err != nil {
		return nil, err
	}

	switch cmd {
	case CmdSet:
		return parseSetCommand(r)
	case CmdGet:
		return parseGetCommand(r)
	case CmdJoin:
		return &CommandJoin{}, nil
	case CmdDel:
		return parseDelCommand(r)
	default:
		return nil, errors.New("invalid command type")
	}
}

func parseSetCommand(r io.Reader) (*CommandSet, error) {
	cmd := &CommandSet{}

	var keyLen int32
	if err := binary.Read(r, binary.LittleEndian, &keyLen); err != nil {
		return nil, err
	}
	cmd.Key = make([]byte, keyLen)
	if err := binary.Read(r, binary.LittleEndian, &cmd.Key); err != nil {
		return nil, err
	}

	var valueLen int32
	if err := binary.Read(r, binary.LittleEndian, &valueLen); err != nil {
		return nil, err
	}
	cmd.Value = make([]byte, valueLen)
	if err := binary.Read(r, binary.LittleEndian, &cmd.Value); err != nil {
		return nil, err
	}

	var ttl int32
	if err := binary.Read(r, binary.LittleEndian, &ttl); err != nil {
		return nil, err
	}
	cmd.TTL = int(ttl)

	return cmd, nil
}

func parseGetCommand(r io.Reader) (*CommandGet, error) {
	cmd := &CommandGet{}

	var keyLen int32
	if err := binary.Read(r, binary.LittleEndian, &keyLen); err != nil {
		return nil, err
	}
	cmd.Key = make([]byte, keyLen)
	if err := binary.Read(r, binary.LittleEndian, &cmd.Key); err != nil {
		return nil, err
	}

	return cmd, nil
}

func parseDelCommand(r io.Reader) (*CommandDelete, error) {
	cmd := &CommandDelete{}

	var keyLen int32
	if err := binary.Read(r, binary.LittleEndian, &keyLen); err != nil {
		return nil, err
	}
	cmd.Key = make([]byte, keyLen)
	if err := binary.Read(r, binary.LittleEndian, &cmd.Key); err != nil {
		return nil, err
	}

	return cmd, nil
}
