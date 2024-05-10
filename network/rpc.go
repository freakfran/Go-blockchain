package network

import (
	"MyChain/core"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
)

type MessageType byte

const (
	MessageTypeTx MessageType = iota + 1
	MessageTypeBlock
)

type RPC struct {
	From    NetAddr
	Payload io.Reader
}

type Message struct {
	Header MessageType
	Data   []byte
}

func NewMessage(t MessageType, data []byte) *Message {
	return &Message{Header: t, Data: data}
}

func (m *Message) Bytes() []byte {
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(m)
	return buf.Bytes()
}

type RPCHandler interface {
	HandleRPC(rpc RPC) error
}
type DefaultRPCHandler struct {
	p RPCProcessor
}

func NewDefaultRPCHandler(p RPCProcessor) *DefaultRPCHandler {
	return &DefaultRPCHandler{p: p}
}

func (h *DefaultRPCHandler) HandleRPC(rpc RPC) error {
	msg := Message{}
	if err := gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil {
		return fmt.Errorf("failed to decode message from %s : %w", rpc.From, err)
	}

	switch msg.Header {
	case MessageTypeTx:
		tx := new(core.Transaction)
		if err := tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data))); err != nil {
			return err
		}
		return h.p.ProcessTransaction(rpc.From, tx)
	default:
		return fmt.Errorf("unknown message type: %v", msg.Header)
	}
}

type RPCProcessor interface {
	ProcessTransaction(NetAddr, *core.Transaction) error
}
