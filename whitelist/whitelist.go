package whitelist

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gorcon/rcon"
)

type Interface interface {
	List(ctx context.Context) ([]string, error)
	Add(ctx context.Context, name string) error
	Remove(ctx context.Context, name string) error
}

type Client struct {
	conn *rcon.Conn
}

func New(conn *rcon.Conn) *Client {
	return &Client{conn}
}

func (c *Client) List(ctx context.Context) ([]string, error) {
	res, err := c.conn.Execute("/whitelist list")
	if err != nil {
		return nil, err
	}

	_, usersStr, found := strings.Cut(res, ": ")
	if !found {
		return nil, errors.New("invalid response")
	}

	users := strings.Split(usersStr, ", ")

	return users, nil
}

func (c *Client) Add(ctx context.Context, name string) error {
	res, err := c.conn.Execute("/whitelist add " + name)
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
	res, err := c.conn.Execute("/whitelist remove " + name)
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
