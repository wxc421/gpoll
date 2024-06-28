package gpoll

import (
	"errors"
	"gpoll/acceptor"
	"gpoll/conn"
	"gpoll/logx"
	"strconv"
	"strings"
)

type GPoll struct {
	reactor  *Reactor
	ip       [4]byte
	port     int
	acceptor acceptor.Acceptor
	options  *Options
	stopChan chan struct{}
}

func (gpoll *GPoll) Close() {
	close(gpoll.stopChan)
}

func (gpoll *GPoll) Start() error {
	err := gpoll.acceptor.Listen()
	if err != nil {
		return err
	}
	gpoll.reactor = NewReactor(gpoll.acceptor.Fd(), gpoll.stopChan)
	if err := gpoll.reactor.Init(); err != nil {
		return err
	}
	onAccept := func(c *conn.Conn) {
		logx.Log.Info("before reactor onAccept")
		gpoll.reactor.onAccept(c)
		logx.Log.Info("after reactor onAccept")
	}
	go func() {
		gpoll.acceptor.StartAccept(onAccept)
	}()

	handler := func(c *conn.Conn) {
		// handler conn
		c.Read()
	}

	go func() {
		// reactor start
		gpoll.reactor.Run(handler)
	}()

	return nil
}

func New(address string, opts ...Option) *GPoll {
	options := completeOptions(opts...)
	ip, port, err := getIPPort(address)
	if err != nil {
		return nil
	}
	return &GPoll{
		ip:       ip,
		port:     port,
		options:  options,
		stopChan: make(chan struct{}),
		acceptor: &acceptor.TcpAcceptor{},
	}
}

func getIPPort(addr string) (ip [4]byte, port int, err error) {
	strs := strings.Split(addr, ":")
	if len(strs) != 2 {
		err = errors.New("addr error")
		return
	}

	if len(strs[0]) != 0 {
		ips := strings.Split(strs[0], ".")
		if len(ips) != 4 {
			err = errors.New("addr error")
			return
		}
		for i := range ips {
			data, err := strconv.Atoi(ips[i])
			if err != nil {
				return ip, 0, err
			}
			ip[i] = byte(data)
		}
	}

	port, err = strconv.Atoi(strs[1])
	return
}
