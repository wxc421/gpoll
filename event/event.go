package event

import (
	"golang.org/x/sys/unix"
	"log"
)

const (
	AddListenEvent = unix.EPOLLIN | unix.EPOLLOUT | unix.EPOLLET | unix.EPOLLRDHUP
)

func IsReadableEvent(event uint32) bool {
	if event&unix.EPOLLIN != 0 {
		return true
	}
	return false
}

func IsClosedEvent(event uint32) bool {
	if event&unix.EPOLLHUP != 0 {
		return true
	}
	if event&unix.EPOLLRDHUP != 0 {
		return true
	}
	return false
}

func debugEvent(event uint32) {
	if event&unix.EPOLLHUP != 0 {
		log.Println("receive EPOLLHUP")
	}
	if event&unix.EPOLLERR != 0 {
		log.Println("receive EPOLLERR")
	}
	if event&unix.EPOLLIN != 0 {
		log.Println("receive EPOLLIN")
	}
	if event&unix.EPOLLOUT != 0 {
		log.Println("receive EPOLLOUT")
	}
	if event&unix.EPOLLRDHUP != 0 {
		log.Println("receive EPOLLRDHUP")
	}
}
