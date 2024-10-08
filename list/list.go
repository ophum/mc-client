package list

import (
	"context"
	"errors"
	"strings"

	"github.com/ophum/mc-client/rcon"
)

type Interface interface {
	List(ctx context.Context) ([]string, error)
}

type Client struct {
	conn *rcon.RetryableRcon
}

var _ Interface = (*Client)(nil)

func New(conn *rcon.RetryableRcon) *Client {
	return &Client{conn}
}

func (c *Client) List(ctx context.Context) ([]string, error) {
	res, err := c.conn.Execute("/list")
	if err != nil {
		return nil, err
	}

	_, usersStr, found := strings.Cut(res, ": ")
	if !found {
		return nil, errors.New("invalid response")
	}

	if usersStr == "" {
		return []string{}, nil
	}

	users := strings.Split(usersStr, ", ")

	return users, nil
}
