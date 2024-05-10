package network

import (
	"bytes"
	"fmt"
	"sync"
)

type LocalTransport struct {
	addr        NetAddr
	consumeChan chan RPC
	//
	//  lock
	//  @Description: 读写锁
	//
	lock  sync.RWMutex
	peers map[NetAddr]*LocalTransport
}

// NewLocalTransport
//
//	@Description: New LocalTransport
//	@param addr
//	@return *LocalTransport
func NewLocalTransport(addr NetAddr) *LocalTransport {
	return &LocalTransport{
		addr:        addr,
		consumeChan: make(chan RPC),
		peers:       make(map[NetAddr]*LocalTransport),
	}
}

// Consume
//
//	@Description:
//	@receiver LocalTransport的方法
//	@return <-chan
func (t *LocalTransport) Consume() <-chan RPC {
	return t.consumeChan
}

// Addr
//
//	@Description: 返回地址
//	@receiver t
//	@return NetAddr
func (t *LocalTransport) Addr() NetAddr {
	return t.addr
}

// Connect
//
//	@Description:
//	@receiver LocalTransport的方法
//	@return <-chan
func (t *LocalTransport) Connect(tr Transport) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.peers[tr.Addr()] = tr.(*LocalTransport)

	return nil
}
func (t *LocalTransport) SendMessage(to NetAddr, payload []byte) error {
	t.lock.RLock()
	defer t.lock.RUnlock()

	peer, ok := t.peers[to]
	if !ok {
		return fmt.Errorf("%s:could not send message to %s", t.addr, to)
	}
	peer.consumeChan <- RPC{
		From:    t.addr,
		Payload: bytes.NewReader(payload),
	}
	return nil
}

// Broadcast 向所有的peer广播给定的payload。
//
// 参数:
//
//	payload []byte - 要广播的数据。
//
// 返回值:
//
//	error - 如果向任意peer发送消息时遇到错误，则返回该错误；否则返回nil。
func (t *LocalTransport) Broadcast(payload []byte) error {
	for _, peer := range t.peers {
		// 尝试向每个peer发送payload，如果发送失败则立即返回错误。
		if err := t.SendMessage(peer.Addr(), payload); err != nil {
			return err
		}
	}

	return nil
}
