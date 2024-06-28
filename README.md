# gpoll
## 简述
gpoll是一个基于linux下epoll的网络框架，目前只能运行在Linux环境下，gpoll可以配置处理网络事件的goroutine数量，相比golang原生库，在海量链接下，可以减少goroutine的开销，从而减少系统资源占用。

还未完成