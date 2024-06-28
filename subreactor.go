package gpoll

import (
	"fmt"
	"gpoll/conn"
	"gpoll/logx"
	"gpoll/utils/structure"
	"syscall"
)

var (
	_ SubReactor = (*ShardSubReactor)(nil)
	_ SubReactor = (*SingleSubReactor)(nil)
)

type SubReactor interface {
	RegisterConnection(c *conn.Conn)
	GetConnection(fd int) *conn.Conn
	Offer(events ...syscall.EpollEvent)
	Polling(stopCh <-chan struct{}, callback func(int))
}

type ShardSubReactor struct {
	container structure.Sharding[*SingleSubReactor]
}

func (shardSubReactor *ShardSubReactor) GetConnection(fd int) *conn.Conn {
	return shardSubReactor.container.GetShard(fd).GetConnection(fd)
}

func NewShardSubReactor(shardingSize int, bufferSize int) *ShardSubReactor {
	return &ShardSubReactor{
		structure.NewSharding[*SingleSubReactor](shardingSize, func() *SingleSubReactor {
			return NewSingleSubReactor(bufferSize)
		}),
	}
}

func (shardSubReactor *ShardSubReactor) Offer(events ...syscall.EpollEvent) {
	for _, event := range events {
		subReactor := shardSubReactor.container.GetShard(int(event.Fd))
		subReactor.Offer(event)
	}
}
func (shardSubReactor *ShardSubReactor) Polling(stopCh <-chan struct{}, callback func(int)) {
	shardSubReactor.container.Iterator(func(subReactor *SingleSubReactor) {
		go func() {
			subReactor.Polling(stopCh, callback)
		}()
	})
}

func (shardSubReactor *ShardSubReactor) RegisterConnection(c *conn.Conn) {
	shardSubReactor.container.GetShard(c.NFd).RegisterConnection(c)
}

// SingleSubReactor represents sub reactor
type SingleSubReactor struct {
	// buffer manage active file descriptors
	buffer *structure.Queue[int]
	// container Save conn ConcurrentMap
	container *structure.ConcurrentMap[int, *conn.Conn]
}

func (singleSubReactor *SingleSubReactor) GetConnection(fd int) *conn.Conn {
	c, _ := singleSubReactor.container.Get(fd)
	return c
}

func NewSingleSubReactor(bufferSize int) *SingleSubReactor {
	return &SingleSubReactor{
		buffer:    structure.NewQueue[int](bufferSize),
		container: structure.NewConcurrentMap[int, *conn.Conn](),
	}
}

func (singleSubReactor *SingleSubReactor) RegisterConnection(c *conn.Conn) {
	logx.Log.Info(fmt.Sprintf("singleSubReactor RegisterConnection:%v", c.NFd))
	singleSubReactor.container.Set(c.NFd, c)
}

func (singleSubReactor *SingleSubReactor) Offer(events ...syscall.EpollEvent) {
	for _, event := range events {
		logx.Log.Info(fmt.Sprintf("singleSubReactor Offer:%+v", event))
		singleSubReactor.buffer.Offer(int(event.Fd))
	}
}

func (singleSubReactor *SingleSubReactor) Polling(stopCh <-chan struct{}, callback func(int)) {
	singleSubReactor.buffer.Polling(stopCh, func(fd int) {
		logx.Log.Info(fmt.Sprintf("singleSubReactor Polling:%v", fd))
		callback(fd)
	})
}
