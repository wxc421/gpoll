package gpoll

import (
	"fmt"
	"gpoll/conn"
	"gpoll/logx"
	"gpoll/poller"
)

type Reactor struct {
	poll     poller.Poller
	fd       int        // listen fd
	sub      SubReactor // manage connections
	stopChan chan struct{}
}

func NewReactor(fd int, stopChan chan struct{}) *Reactor {
	return &Reactor{
		fd:       fd,
		sub:      NewShardSubReactor(10, 10),
		stopChan: stopChan,
	}
}

func (reactor *Reactor) Init() error {
	poll, err := poller.NewPoll()
	if err != nil {
		return err
	}
	reactor.poll = poll
	return nil
}

func (reactor *Reactor) onAccept(c *conn.Conn) {
	if err := reactor.poll.Add(c.NFd); err != nil {
		logx.Log.Warn(fmt.Sprintf("reactor.poll.Add(%v) err:%v", c.NFd, err))
	}
	reactor.sub.RegisterConnection(c)
}

func (reactor *Reactor) wrapHandler(handler conn.ConnectionHandler) func(active int) {
	return func(active int) {
		c := reactor.sub.GetConnection(active)
		if c == nil {
			return
		}
		handler(c)
	}
}

func (reactor *Reactor) Run(handler conn.ConnectionHandler) {
	go func() {
		// start sub reactor
		reactor.sub.Polling(reactor.stopChan, reactor.wrapHandler(handler))
	}()
	// start epoll wait
	go reactor.wait()
}

func (reactor *Reactor) wait() {
	logx.Log.Info("reactor:wait")
	for {
		select {
		case <-reactor.stopChan:
			return
		default:
			events, err := reactor.poll.Wait()
			logx.Log.Info(fmt.Sprintf("event:%v", len(events)))
			if err != nil {
				logx.Log.Warn(err.Error())
				continue
			}
			if len(events) <= 0 {
				continue
			}
			reactor.sub.Offer(events...)
		}

	}
}
