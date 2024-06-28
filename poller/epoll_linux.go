package poller

import (
	"golang.org/x/sys/unix"
	"gpoll/event"
	"syscall"
)

type Epoll struct {
	epFd     int
	listenFd int
	events   []syscall.EpollEvent
}

func (e *Epoll) Add(fd int) error {

	return syscall.EpollCtl(e.epFd, syscall.EPOLL_CTL_ADD, fd, &syscall.EpollEvent{
		Events: event.AddListenEvent,
		Fd:     int32(fd),
	})

}
func (e *Epoll) Remove(fd int) error {
	// 移除文件描述符的监听
	if err := syscall.EpollCtl(e.epFd, syscall.EPOLL_CTL_DEL, fd, nil); err != nil {
		return err
	}

	return nil
}

func (e *Epoll) Close() error {
	return unix.Close(e.epFd)
}

func (e *Epoll) Wait() ([]syscall.EpollEvent, error) {
	n, err := syscall.EpollWait(e.epFd, e.events, 3*1000)
	if n <= 0 && err != nil {
		return nil, err
	}
	return e.events[:n], nil

}

func NewPoll() (*Epoll, error) {

	epFd, err := syscall.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &Epoll{
		epFd:   epFd,
		events: make([]syscall.EpollEvent, 100, 100),
	}, nil
}
