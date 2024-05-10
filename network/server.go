package network

import (
	"MyChain/core"
	"MyChain/crypto"
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

var defaultBlockTime = 5 * time.Second

// ServerOpts represents the options used to create a new Server.
type ServerOpts struct {
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	Transports    []Transport
	BlockTime     time.Duration
	PrivateKey    *crypto.PrivateKey
}

// Server represents a server that listens for incoming connections and handles them.
type Server struct {
	ServerOpts
	memPool     *TxPool
	isValidator bool
	rpcChan     chan RPC
	quitChan    chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}
	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}
	s := &Server{
		ServerOpts:  opts,
		memPool:     NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcChan:     make(chan RPC),
		quitChan:    make(chan struct{}, 1),
	}
	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}
	if s.isValidator {
		s.validateLoop()
	}
	return s
}

// ProcessMessage 处理接收到的消息。
//
// 参数:
// message: 包含待处理消息的 DecodeMessage 指针。
//
// 返回值:
// 返回处理过程中可能出现的错误，如果处理成功则返回 nil。
func (s *Server) ProcessMessage(message *DecodeMessage) error {
	// 根据消息类型使用多态方式调用相应的处理函数
	switch t := message.Data.(type) {
	case *core.Transaction:
		// 处理交易消息
		return s.processTransaction(t)
	}

	// 如果消息类型不匹配任何已知类型，则不处理，返回 nil
	return nil
}

// broadcast将给定的payload广播到所有传输层。
// 参数:
// - payload: 需要广播的数据负载。
// 返回值:
// - error: 广播过程中遇到的任何错误。
func (s *Server) broadcast(payload []byte) error {
	for _, tr := range s.Transports {
		// 尝试通过每个传输层进行广播，如果有任何一个失败，则返回错误。
		if err := tr.Broadcast(payload); err != nil {
			return err
		}
	}

	return nil
}

// broadcastTx将一个交易广播到所有连接的客户端。
// 参数:
// - tx: 需要广播的交易实例。
// 返回值:
// - error: 在编码交易或广播过程中遇到的任何错误。
func (s *Server) broadcastTx(tx *core.Transaction) error {
	buf := &bytes.Buffer{}
	// 使用Gob编码器将交易编码为字节流。
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	// 创建一个包含交易消息的消息对象。
	msg := NewMessage(MessageTypeTx, buf.Bytes())
	// 调用broadcast函数将交易消息广播出去。
	return s.broadcast(msg.Bytes())
}

// processTransaction 处理一个交易，首先验证交易的有效性，然后检查交易是否已经存在于内存池中。
// 如果交易无效或已存在，则不进行处理；否则，将交易添加到内存池中。
// 参数:
//
//	tx *core.Transaction: 需要处理的交易对象。
//
// 返回值:
//
//	error: 如果处理过程中出现错误，则返回错误对象；否则返回nil。
func (s *Server) processTransaction(tx *core.Transaction) error {
	// 计算交易的哈希值
	hash := tx.Hash(core.TxHasher{})
	// 检查交易是否已存在于内存池中
	if s.memPool.Has(hash) {
		logrus.WithFields(logrus.Fields{"hash": hash, "memPool length": s.memPool.Len()}).Infoln("Transaction already exists in memPool")
		return nil
	}
	// 验证交易的有效性
	if err := tx.Verify(); err != nil {
		return err
	}

	tx.SetFirstSeen(time.Now().UnixNano())
	// 将新交易添加到内存池
	logrus.WithFields(logrus.Fields{"hash": hash}).Infoln("Add new transaction to memPool")
	// 异步广播该交易
	go func() {
		err := s.broadcastTx(tx)
		if err != nil {
			logrus.WithFields(logrus.Fields{"hash": hash}).Errorf("Broadcast tx error:%v", err)
		}
	}()
	return s.memPool.Add(tx)
}

// Start 启动服务器，初始化传输层，然后进入一个循环，不断监听RPC请求、退出信号和定时器事件。
// 在接收到退出信号时，会退出循环并关闭服务器。
func (s *Server) Start() {
	// 初始化传输层
	s.initTransports()

free:
	for {
		select {
		case rpc := <-s.rpcChan:
			// 处理RPC请求
			msg, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				logrus.Error(err)
			}
			if err := s.ProcessMessage(msg); err != nil {
				logrus.Error(err)
			}
		case <-s.quitChan:
			// 接收到退出信号，退出循环
			break free
		}
	}

	// 打印服务器关闭信息
	fmt.Println("Server shutdown")
}

func (s *Server) validateLoop() {
	ticker := time.NewTicker(s.BlockTime)
	for {
		<-ticker.C
		s.createNewBlock()
	}
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
