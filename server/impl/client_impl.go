package impl

import (
	"context"
	"go_grpc_demo/server"
	"net"
	"sync"
	"time"
)

type client struct {
	id         int64 // 每一个客户端有一个唯一 ID
	hub        *hub
	ctx        context.Context
	conn       net.Conn
	auth       bool
	mu         sync.Mutex
	authTicker *time.Ticker
}

func (c *client) Id() int64 {
	return c.id
}

func (c *client) Hub() server.Hub {
	return c.hub
}

func (c *client) Context() context.Context {
	return c.ctx
}

func (c *client) SetAuthed() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.authTicker != nil {
		c.authTicker.Stop()
		c.authTicker = nil
	}
	c.auth = true
}

func (c *client) Auth() bool {
	return c.auth
}

func (c *client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil
	}
}

//// private

func (c *client) waitForAuth(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.authTicker = time.NewTicker(d)
	go func() {
		select {
		case <-c.authTicker.C:
			if !c.auth {
				c.Close()
			}
			c.authTicker = nil
			return
		}
	}()
}
