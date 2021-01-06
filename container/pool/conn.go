package pool

import (
	"net"
	"sync"
	"time"
)

type Conn struct {
	net.Conn
	mu       sync.RWMutex
	c        *channelPool
	unusable bool
	t 		time.Time
}

// Close() puts the given connects back to the pool instead of closing it.
func (c *Conn) Close() error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.unusable {
		if c.Conn != nil {
			return c.Conn.Close()
		}
		return nil
	}
	return c.c.put(c)
}

// MarkUnusable() marks the connection not usable any more, to let the pool close it instead of returning it to pool.
func (c *Conn) MarkUnusable() {
	c.mu.RLock()
	c.unusable = true
	c.mu.RUnlock()
}

// newConn wraps a standard net.Conn to a poolConn net.Conn.
func (p *channelPool) wrapConn(conn net.Conn) *Conn {
	c := &Conn{c: p, Conn: conn, t: time.Now()}
	return c
}