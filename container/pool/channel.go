package pool

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

type Config struct {
	InitialCap	int
	MaxCap		int
	Factory		Factory
	IdleTimeout time.Duration
}

// channelPool implements the Pool interface based on buffered channels.
type channelPool struct {
	// storage for our net.Conn connections
	mu    sync.RWMutex
	conns chan *Conn

	// net.Conn generator
	factory Factory

	idleTimeout time.Duration
}


// Factory is a function to create new connections.
type Factory func() (net.Conn, error)

// NewChannelPool returns a new pool based on buffered channels with an initial
// capacity and maximum capacity. Factory is used when initial capacity is
// greater than zero to fill the pool. A zero initialCap doesn't fill the Pool
// until a new Get() is called. During a Get(), If there is no new connection
// available in the pool, a new connection will be created via the Factory()
// method.
func New(c *Config) (Pool, error) {
	if c.InitialCap < 0 || c.MaxCap <= 0 || c.InitialCap > c.MaxCap {
		return nil, errors.New("invalid capacity settings")
	}

	if c.Factory == nil {
		return nil, errors.New("invalid factory func settings")
	}

	pool := &channelPool{
		conns:   make(chan *Conn, c.MaxCap),
		factory: c.Factory,
		idleTimeout: c.IdleTimeout,
	}

	// create initial connections, if something goes wrong,
	// just close the pool error out.
	for i := 0; i < c.InitialCap; i++ {
		conn, err := c.Factory()
		if err != nil {
			pool.Close()
			return nil, fmt.Errorf("factory is not able to fill the pool: %s", err)
		}
		pool.conns <- &Conn{Conn: conn, c: pool, t: time.Now()}
	}

	return pool, nil
}

func (p *channelPool) getConnsAndFactory() (chan *Conn, Factory) {
	p.mu.RLock()
	conns := p.conns
	factory := p.factory
	p.mu.RUnlock()
	return conns, factory
}

// Get implements the Pool interfaces Get() method. If there is no new
// connection available in the pool, a new connection will be created via the
// Factory() method.
func (p *channelPool) Get() (net.Conn, error) {
	conns, factory := p.getConnsAndFactory()
	if conns == nil {
		return nil, ErrClosed
	}

	// wrap our connections with out custom net.Conn implementation (wrapConn
	// method) that puts the connection back to the pool if it's closed.
	for {
		select {
		case conn := <-conns:
			if conn == nil {
				return nil, ErrClosed
			}

			if timeout := p.idleTimeout; timeout > 0 {
				if conn.t.Add(timeout).Before(time.Now()) {
					conn.MarkUnusable()
					conn.Close()
					continue
				}
			}
			conn.t = time.Now()
			return conn, nil
		default:
			conn, err := factory()
			if err != nil {
				return nil, err
			}

			return p.wrapConn(conn), nil
		}
	}

}

// put puts the connection back to the pool. If the pool is full or closed,
// conn is simply closed. A nil conn will be rejected.
func (p *channelPool) put(conn *Conn) error {
	if conn.Conn == nil {
		return errors.New("connection is nil. rejecting")
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.conns == nil {
		// pool is closed, close passed connection
		return conn.Close()
	}

	// put the resource back into the pool. If the pool is full, this will
	// block and the default case will be executed.
	select {
	case p.conns <- conn:
		return nil
	default:
		// pool is full, close passed connection
		conn.MarkUnusable()
		return conn.Close()
	}
}

func (p *channelPool) Close() {
	p.mu.Lock()
	conns := p.conns
	p.conns = nil
	p.factory = nil
	p.mu.Unlock()

	if conns == nil {
		return
	}

	close(conns)
	for conn := range conns {
		conn.Close()
	}
}

func (p *channelPool) Len() int {
	conns, _ := p.getConnsAndFactory()
	return len(conns)
}
