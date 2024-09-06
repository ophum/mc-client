package rcon

import (
	"sync"

	"github.com/gorcon/rcon"
)

type RetryableRcon struct {
	cache    *rcon.Conn
	address  string
	password string
	options  []rcon.Option
	mu       sync.Mutex
}

func New(address, password string, options ...rcon.Option) (*RetryableRcon, error) {
	conn, err := rcon.Dial(address, password, options...)
	if err != nil {
		return nil, err
	}
	return &RetryableRcon{
		cache:    conn,
		address:  address,
		password: password,
		options:  options,
		mu:       sync.Mutex{},
	}, nil
}

func (c *RetryableRcon) Close() error {
	return c.cache.Close()
}

// 最大2回実行しエラーになったら、rcon.Dialしなおして再度実行する
// それでもエラーになったらエラーを返す
func (c *RetryableRcon) Execute(command string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var err error
	var res string
	for i := 0; i < 2; i++ {
		res, err = c.cache.Execute(command)
		if err == nil {
			break
		}
	}

	if err != nil {
		c.cache, err = rcon.Dial(c.address, c.password, c.options...)
		if err != nil {
			return "", err
		}

		res, err = c.cache.Execute(command)
		if err != nil {
			return "", err
		}
	}

	return res, nil
}
