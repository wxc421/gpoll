package acceptor

import (
	"golang.org/x/sys/unix"
	"gpoll/conn"
	"syscall"
)

type TcpAcceptor struct {
	fd   int
	stop chan struct{}
}

func (acceptor *TcpAcceptor) close() error {
	return unix.Close(acceptor.fd)
}

func (acceptor *TcpAcceptor) Fd() int {
	return acceptor.fd
}

func (acceptor *TcpAcceptor) Listen() error {
	var err error
	// socket()
	if acceptor.fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0); err != nil {
		return err
	}
	// bind()
	err = syscall.Bind(acceptor.fd, &syscall.SockaddrInet4{
		Port: 8080,
		Addr: [4]byte{0, 0, 0, 0},
	})
	if err != nil {
		return err
	}
	// listen()
	if err = syscall.Listen(acceptor.fd, 1024); err != nil {
		return err
	}
	return nil
}
func (acceptor *TcpAcceptor) StartAccept(onAccept func(c *conn.Conn)) {
	for {
		select {
		case <-acceptor.stop:
			_ = acceptor.close()
		default:
			nfd, sa, err := syscall.Accept(acceptor.fd)
			if err != nil {
				continue
			}
			// Set to a non-blocking state
			err = syscall.SetNonblock(nfd, true)
			if err != nil {
				continue
			}
			c := &conn.Conn{
				ListenFd:   acceptor.fd,
				NFd:        nfd,
				SocketAddr: sa.(*syscall.SockaddrInet4),
			}
			onAccept(c)
		}

	}
}
