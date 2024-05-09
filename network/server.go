package network

import (
	"fmt"
	"time"
)

type ServerOpts struct {
	Transports []Transport
}

type Server struct {
	ServerOpts
	rpcChan  chan RPC
	quitChan chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	return &Server{
		ServerOpts: opts,
		rpcChan:    make(chan RPC),
		quitChan:   make(chan struct{}, 1),
	}
}

func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(5 * time.Second)

free:
	for {
		select {
		case rpc := <-s.rpcChan:
			fmt.Printf("%+v\n", rpc)
		case <-s.quitChan:
			break free
		case <-ticker.C:
			fmt.Println("Do something every 5 seconds")
		}
	}

	fmt.Println("Server shutdown")
}

// initTransports 初始化所有的传输介质，为每个传输介质创建一个goroutine，
// goroutine 不断地从传输介质中消费RPC请求，并将这些请求发送到服务器的rpcChan通道中。
func (s *Server) initTransports() {
	// 遍历所有的传输介质
	for _, tr := range s.Transports {
		go func(tr Transport) {
			// 持续从当前传输介质中消费RPC请求
			for rpc := range tr.Consume() {
				// 将消费到的RPC请求发送到服务器的rpcChan通道
				s.rpcChan <- rpc
			}
		}(tr)
	}
}
