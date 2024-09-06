package mcclient

import (
	"context"
	"net"
	"strconv"

	"github.com/ophum/mc-client/list"
	"github.com/ophum/mc-client/rcon"
	"github.com/ophum/mc-client/whitelist"
)

type Interface interface {
	Whitelist() whitelist.Interface

	list.Interface

	Close() error
}

type Client struct {
	conn *rcon.RetryableRcon
}

var _ Interface = (*Client)(nil)

func New(host string, port int, password string) (*Client, error) {
	c, err := rcon.New(net.JoinHostPort(host, strconv.Itoa(port)), password)
	if err != nil {
		return nil, err
	}
	return &Client{c}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Whitelist() whitelist.Interface {
	return whitelist.New(c.conn)
}

func (c *Client) List(ctx context.Context) ([]string, error) {
	return list.New(c.conn).List(ctx)
}
