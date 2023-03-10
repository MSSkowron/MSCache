package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type Command byte

const (
	CmdNone Command = iota
	CmdSet
	CmdGet
	CmdDel
	CmdJoin
)

type Status byte

const (
	StatusNone Status = iota
	StatusOK
	StatusError
	StatusKeyNotFound
)

type ResponseSet struct {
	Status Status
}

type ResponseGet struct {
	Status Status
	Value  []byte
}

type ResponseDelete struct {
	Status Status
}

type CommandSet struct {
	Key   []byte
	Value []byte
	TTL   int
}

type CommandGet struct {
	Key []byte
}

type CommandDelete struct {
	Key []byte
}

type CommandJoin struct{}

func (s Status) String() string {
	switch s {
	case StatusOK:
		return "OK"
	case StatusError:
		return "ERROR"
	case StatusKeyNotFound:
		return "NOT FOUND"
	default:
		return "NONE"
	}
}

func (r *ResponseSet) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, r.Status); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *ResponseGet) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, r.Status); err != nil {
		return nil, err
	}

	valueLen := int32(len(r.Value))
	if err := binary.Write(buf, binary.LittleEndian, valueLen); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, r.Value); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *ResponseDelete) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, r.Status); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

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

func ParseSetResponse(r io.Reader) (*ResponseSet, error) {
	resp := &ResponseSet{}

	if err := binary.Read(r, binary.LittleEndian, &resp.Status); err != nil {
		return nil, err
	}

	return resp, nil
}

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

func ParseDeleteResponse(r io.Reader) (*ResponseDelete, error) {
	resp := &ResponseDelete{}

	if err := binary.Read(r, binary.LittleEndian, &resp.Status); err != nil {
		return nil, err
	}

	return resp, nil
}

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
