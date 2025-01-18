package list

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/ophum/mc-client/rcon"
)

type Interface interface {
	List(ctx context.Context) ([]string, error)
}

type Client struct {
	conn       *rcon.RetryableRcon
	serverType rcon.ServerType
}

var _ Interface = (*Client)(nil)

func New(conn *rcon.RetryableRcon, serverType rcon.ServerType) *Client {
	return &Client{conn, serverType}
}

func (c *Client) List(ctx context.Context) ([]string, error) {
	command := map[rcon.ServerType]string{
		rcon.ServerTypeVanilla: "/list",
		rcon.ServerTypeSpigot:  "list",
	}
	res, err := c.conn.Execute(command[c.serverType])
	if err != nil {
		return nil, err
	}

	log.Println(res)
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
