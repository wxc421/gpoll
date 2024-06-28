package conn

import (
	"fmt"
	"gpoll/logx"
	"gpoll/poller"
	"syscall"
)

type Conn struct {
	ListenFd   int
	NFd        int // 文件描述符
	poller     poller.Poller
	SocketAddr *syscall.SockaddrInet4
}

// ConnectionHandler represents a connection handler
type ConnectionHandler func(conn *Conn)

func (c *Conn) Close() error {
	if err := c.poller.Remove(c.NFd); err != nil {
		return err
	}
	// 关闭文件描述符
	return syscall.Close(c.NFd)
}

func (c *Conn) Read() {
	logx.Log.Info("into read")
	buffer := make([]byte, 1024)
	n, err := syscall.Read(c.NFd, buffer)
	if err != nil {
		logx.Log.Warn("read error:", err)
		return
	}
	logx.Log.Info(fmt.Sprintf("read:%v", string(buffer[:n])))
	fmt.Println("read:", string(buffer[:n]))
}
