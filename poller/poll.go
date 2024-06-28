package poller

import "syscall"

type Poller interface {
	// Wait linux: EpollWait
	Wait() ([]syscall.EpollEvent, error)
	// Add linux: EpollCtl
	Add(fd int) error
	// Close linux: EpollCtl
	Close() error
	Remove(fd int) error
}
