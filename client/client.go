package client

import (
	"context"
	"fmt"
	"net"

	"github.com/MSSkowron/mscache/protocol"
)

type Client struct {
	conn net.Conn
}

func New(endpoint string) (*Client, error) {
	conn, err := net.Dial("tcp", endpoint)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Set(ctx context.Context, key, value []byte, ttl int) error {
	cmd := &protocol.CommandSet{
		Key:   key,
		Value: value,
		TTL:   ttl,
	}

	b, err := cmd.Bytes()
	if err != nil {
		return err
	}

	_, err = c.conn.Write(b)
	if err != nil {
		return err
	}

	resp, err := protocol.ParseSetResponse(c.conn)
	if err != nil {
		return err
	}

	if resp.Status != protocol.StatusOK {
		return fmt.Errorf("server responded with non OK status [%s]", resp.Status)
	}

	return nil
}

func (c *Client) Get(ctx context.Context, key []byte) ([]byte, error) {
	cmd := &protocol.CommandGet{
		Key: key,
	}

	b, err := cmd.Bytes()
	if err != nil {
		return nil, err
	}

	_, err = c.conn.Write(b)
	if err != nil {
		return nil, err
	}

	resp, err := protocol.ParseGetResponse(c.conn)
	if err != nil {
		return nil, err
	}

	if resp.Status != protocol.StatusOK {
		return nil, fmt.Errorf("server responded with non OK status [%s]", resp.Status)
	}

	return resp.Value, nil
}
