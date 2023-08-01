package client

import (
	"context"
	"fmt"
	"net"

	"github.com/MSSkowron/mscache/protocol"
)

// Client is a client for the cache server.
type Client struct {
	conn net.Conn
}

// New creates a new client.
func New(endpoint string) (*Client, error) {
	conn, err := net.Dial("tcp", endpoint)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

// Close closes the connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

// Get sends a get command to the server.
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

// Set sends a set command to the server.
// ttl is in seconds.
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

// Delete sends a delete command to the server.
func (c *Client) Delete(ctx context.Context, key []byte) error {
	cmd := &protocol.CommandDelete{
		Key: key,
	}

	b, err := cmd.Bytes()
	if err != nil {
		return err
	}

	_, err = c.conn.Write(b)
	if err != nil {
		return err
	}

	resp, err := protocol.ParseDeleteResponse(c.conn)
	if err != nil {
		return err
	}

	if resp.Status != protocol.StatusOK {
		return fmt.Errorf("server responded with non OK status [%s]", resp.Status)
	}

	return nil
}

// String returns the string representation of the client which is the client's address.
func (c Client) String() string {
	return c.conn.RemoteAddr().String()
}
