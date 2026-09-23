package httpapi

import (
	"net"
	"sync"
)

// LimitedListener closes overflow connections immediately; no permit wait queue.
type LimitedListener struct {
	net.Listener
	permits chan struct{}
}

func LimitConnections(l net.Listener) *LimitedListener {
	return &LimitedListener{l, make(chan struct{}, 64)}
}
func (l *LimitedListener) Accept() (net.Conn, error) {
	for {
		c, err := l.Listener.Accept()
		if err != nil {
			return nil, err
		}
		select {
		case l.permits <- struct{}{}:
			return &countedConn{Conn: c, release: func() { <-l.permits }}, nil
		default:
			_ = c.Close()
		}
	}
}

type countedConn struct {
	net.Conn
	once    sync.Once
	release func()
}

func (c *countedConn) Close() error { err := c.Conn.Close(); c.once.Do(c.release); return err }
