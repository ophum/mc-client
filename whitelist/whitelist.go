package whitelist

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/ophum/mc-client/rcon"
)

type Interface interface {
	List(ctx context.Context) ([]string, error)
	Add(ctx context.Context, name string) error
	Remove(ctx context.Context, name string) error
}

type Client struct {
	conn       *rcon.RetryableRcon
	serverType rcon.ServerType
	command    string
}

func New(conn *rcon.RetryableRcon, serverType rcon.ServerType) *Client {
	command := ""
	switch serverType {
	case rcon.ServerTypeVanilla:
		command = "/whitelist"
	case rcon.ServerTypeSpigot:
		command = "whitelist"
	}
	return &Client{conn, serverType, command}
}

func (c *Client) commandWithArgs(args string) string {
	return c.command + " " + args
}
func (c *Client) List(ctx context.Context) ([]string, error) {

	res, err := c.conn.Execute(c.commandWithArgs("list"))
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

func (c *Client) Add(ctx context.Context, name string) error {
	res, err := c.conn.Execute(c.commandWithArgs("add " + name))
	if err != nil {
		return err
	}

	if res == "That player does not exist" {
		return errors.New("invalid player name, not found")
	}

	if res == "Player is already whitelisted" {
		return nil
	}

	if res == fmt.Sprintf("Added %s to the whitelist", name) {
		return nil
	}

	return errors.New("Unknown error")
}

func (c *Client) Remove(ctx context.Context, name string) error {
	res, err := c.conn.Execute(c.commandWithArgs("remove " + name))
	if err != nil {
		return err
	}

	if res == "That player does not exist" {
		return errors.New("invalid player name, not found")
	}

	if res == fmt.Sprintf("Removed %s from the whitelist", name) {
		return nil
	}

	log.Println(res)
	return errors.New("Unknown error")
}
