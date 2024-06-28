package acceptor

import "gpoll/conn"

// Acceptor represents a server for accepting connections
type Acceptor interface {
	// Listen runs the thread that will receive the connection
	Listen() error
	// StartAccept start accept conn
	StartAccept(onAccept func(c *conn.Conn))
	// Fd get fd
	Fd() int
}
