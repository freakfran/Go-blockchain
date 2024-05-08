package network

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestLocalTransport_Connect
//
//	@Description: 测试LocalTransport的Connect方法
//	@param t
func TestLocalTransport_Connect(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")

	tra.Connect(trb)
	trb.Connect(tra)

	assert.Equal(t, tra.peers[trb.addr], trb)
	assert.Equal(t, trb.peers[tra.addr], tra)
}

// TestLocalTransport_SendMessage
//
//	@Description: 测试LocalTransport的SendMessage方法
//	@param t
func TestLocalTransport_SendMessage(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")
	tra.Connect(trb)
	trb.Connect(tra)
	msg := []byte("hello")
	// 显式启动一个goroutine来消费trb的消息
	go func() {
		for {
			rpc := <-trb.Consume()
			// 这里可以添加处理消息的逻辑，或者简单记录消息以确保通道不被阻塞
			assert.Equal(t, rpc.Payload, msg)
			assert.Equal(t, rpc.From, tra.Addr())
		}
	}()

	assert.Nil(t, tra.SendMessage(trb.Addr(), msg))
}
