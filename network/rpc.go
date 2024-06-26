package network

import (
	"MyChain/core"
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/sirupsen/logrus"
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

type DecodeMessage struct {
	From NetAddr
	Data any
}
type RPCDecodeFunc func(rpc RPC) (message *DecodeMessage, err error)

func DefaultRPCDecodeFunc(rpc RPC) (message *DecodeMessage, err error) {
	msg := Message{}
	if err := gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil {
		return nil, fmt.Errorf("failed to decode message from %s : %w", rpc.From, err)
	}

	logrus.WithFields(logrus.Fields{
		"from": rpc.From, "type": msg.Header,
	}).Infoln("new incoming message")

	switch msg.Header {
	case MessageTypeTx:
		tx := new(core.Transaction)
		if err := tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, err
		}
		return &DecodeMessage{From: rpc.From, Data: tx}, nil
	default:
		return nil, fmt.Errorf("unknown message type: %v", msg.Header)
	}
}

type RPCProcessor interface {
	ProcessMessage(message *DecodeMessage) error
}
