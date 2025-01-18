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
	Raw(command string) (string, error)

	Whitelist() whitelist.Interface

	list.Interface

	Close() error
}

type Client struct {
	conn       *rcon.RetryableRcon
	serverType rcon.ServerType
}

var _ Interface = (*Client)(nil)

type Options struct {
	serverType rcon.ServerType
}

type Option func(o *Options)

func WithServerType(serverType rcon.ServerType) Option {
	return func(o *Options) {
		o.serverType = serverType
	}
}

func New(host string, port int, password string, options ...Option) (*Client, error) {
	opts := &Options{
		serverType: rcon.ServerTypeVanilla,
	}
	for _, opt := range options {
		opt(opts)
	}

	c, err := rcon.New(net.JoinHostPort(host, strconv.Itoa(port)), password)
	if err != nil {
		return nil, err
	}
	return &Client{c, opts.serverType}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Raw(command string) (string, error) {
	return c.conn.Execute(command)
}

func (c *Client) Whitelist() whitelist.Interface {
	return whitelist.New(c.conn, c.serverType)
}

func (c *Client) List(ctx context.Context) ([]string, error) {
	return list.New(c.conn, c.serverType).List(ctx)
}
