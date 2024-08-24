package mcclient

import (
	"net"
	"strconv"

	"github.com/gorcon/rcon"
	"github.com/ophum/mc-client/whitelist"
)

type Interface interface {
	Whitelist() whitelist.Interface

	Close() error
}

type Client struct {
	conn *rcon.Conn
}

var _ Interface = (*Client)(nil)

func New(host string, port int, password string) (*Client, error) {
	c, err := rcon.Dial(net.JoinHostPort(host, strconv.Itoa(port)), password)
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
