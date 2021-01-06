# pool

Pool is a thread safe connection pool for net.Conn interface. It can be used to manage and reuse connections.

## ## Base Usage:
```go
    // create a factory() to be used with channel based pool
    factory := func() (net.Conn, error) { return net.Dial("tcp", "127.0.0.1:8080") }

    // create a new channel based pool with an initial capacity of 5 and maximum
    // capacity of 30. The factory will create 5 initial connections and put it
    // into the pool. The connection idle time is 1 hour
    p, err := pool.New(&pool.Config{InitialCap: 5, MaxCap: 30, Factory: factory, IdleTimeout: 1 * time.Hour})

    // now you can get a connection from the pool, if there is no connection
    // available it will create a new one via the factory function.
    conn, err := p.Get()

    // do something with conn and put it back to the pool by closing the connection
    // (this doesn't close the underlying connection instead it's putting it back
    // to the pool).
    conn.Close()

    // close the underlying connection instead of returning it to pool
    // it is useful when acceptor has already closed connection and conn.Write() returns error
    if pc, ok := conn.(*pool.PoolConn); ok {
        pc.MarkUnusable()
        pc.Close()
    }

    // close pool any time you want, this closes all the connections inside a pool
    p.Close()
    
    // currently available connections in the pool
    current := p.Len()
```
