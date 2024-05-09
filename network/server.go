package network

import (
	"MyChain/core"
	"MyChain/crypto"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

type ServerOpts struct {
	Transports []Transport
	BlockTime  time.Duration
	PrivateKey *crypto.PrivateKey
}

type Server struct {
	ServerOpts
	blockTime   time.Duration
	memPool     *TxPool
	isValidator bool
	rpcChan     chan RPC
	quitChan    chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	return &Server{
		ServerOpts:  opts,
		blockTime:   opts.BlockTime,
		memPool:     NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcChan:     make(chan RPC),
		quitChan:    make(chan struct{}, 1),
	}
}

// HandleTransaction 处理一个交易，首先验证交易的有效性，然后检查交易是否已经存在于内存池中。
// 如果交易无效或已存在，则不进行处理；否则，将交易添加到内存池中。
// 参数:
//
//	tx *core.Transaction: 需要处理的交易对象。
//
// 返回值:
//
//	error: 如果处理过程中出现错误，则返回错误对象；否则返回nil。
func (s *Server) HandleTransaction(tx *core.Transaction) error {
	// 验证交易的有效性
	if err := tx.Verify(); err != nil {
		return err
	}
	// 计算交易的哈希值
	hash := tx.Hash(core.TxHasher{})
	// 检查交易是否已存在于内存池中
	if s.memPool.Has(hash) {
		logrus.WithFields(logrus.Fields{"hash": hash}).Infoln("Transaction already exists in memPool")
		return nil
	}
	// 将新交易添加到内存池
	logrus.WithFields(logrus.Fields{"hash": hash}).Infoln("Add new transaction to memPool")
	return s.memPool.Add(tx)
}

// Start 启动服务器，初始化传输层，然后进入一个循环，不断监听RPC请求、退出信号和定时器事件。
// 在接收到退出信号时，会退出循环并关闭服务器。
func (s *Server) Start() {
	// 初始化传输层
	s.initTransports()
	// 创建一个定时器，按照设定的块间隔时间触发
	ticker := time.NewTicker(s.blockTime)

free:
	for {
		select {
		case rpc := <-s.rpcChan:
			// 处理RPC请求
			fmt.Printf("%+v\n", rpc)
		case <-s.quitChan:
			// 接收到退出信号，退出循环
			break free
		case <-ticker.C:
			// 定时器触发，如果是验证者，则创建新块
			if s.isValidator {
				s.createNewBlock()
			}

		}
	}

	// 打印服务器关闭信息
	fmt.Println("Server shutdown")
}

func (s *Server) createNewBlock() error {
	fmt.Printf("creating a new block")
	return nil
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
