package network

import (
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
	"time"
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

	go func() {
		for rpc := range trb.Consume() {
			buf := make([]byte, len(msg))
			n, err := rpc.Payload.Read(buf)
			assert.Nil(t, err)
			assert.Equal(t, n, len(msg))
			assert.Equal(t, buf, msg)
			assert.Equal(t, rpc.From, tra.Addr())
		}
	}()

	assert.Nil(t, tra.SendMessage(trb.Addr(), msg))
	time.Sleep(5 * time.Second)

}
func TestLocalTransport_Broadcast(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")
	trc := NewLocalTransport("C")
	tra.Connect(trb)
	tra.Connect(trc)

	msg := []byte("hello")

	go func() {
		for rpcb := range trb.Consume() {
			b, err := io.ReadAll(rpcb.Payload)
			assert.Nil(t, err)
			assert.Equal(t, b, msg)
		}
	}()

	go func() {
		for rpcc := range trc.Consume() {
			c, err := io.ReadAll(rpcc.Payload)
			assert.Nil(t, err)
			assert.Equal(t, c, msg)
		}
	}()
	assert.Nil(t, tra.Broadcast(msg))

	time.Sleep(5 * time.Second)
}
