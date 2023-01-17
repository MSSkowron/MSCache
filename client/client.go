package client

import (
	"context"
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

func (c *Client) Set(ctx context.Context, key, value []byte, ttl int) ([]byte, error) {
	cmd := &protocol.CommandSet{
		Key:   key,
		Value: value,
		TTL:   ttl,
	}

	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return nil, err
	}

	return nil, err
}
