package cache

import (
	"github.com/gomodule/redigo/redis"
)

type conn struct {
	redis.Conn
}

// Returns a new redis connection.
func NewConn() (conn, error) {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		c.Close()
		return conn{}, err
	}
	return conn{c}, nil
}

// SET command.
func (c *conn) Set(args ...interface{}) (interface{}, error) {
	cmd := "SET"
	reply, err := c.Do(cmd, args...)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// GET command. We stringify the response.
func (c *conn) Get(args ...interface{}) (string, error) {
	cmd := "GET"
	reply, err := redis.String(c.Do(cmd, args...))
	if err != nil {
		return "", err
	}
	return reply, nil
}
