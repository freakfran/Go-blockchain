package main

import (
	"MyChain/core"
	"MyChain/crypto"
	"MyChain/network"
	"bytes"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	trLocal := network.NewLocalTransport("LOCAL")
	trRemote := network.NewLocalTransport("REMOTE")

	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	go func() {
		for {
			//trRemote.SendMessage(trLocal.Addr(), []byte("hello"))
			err := sendTransaction(trRemote, trLocal.Addr())
			if err != nil {
				logrus.Error(err)
			}
			time.Sleep(3 * time.Second)
		}

	}()
	opts := network.ServerOpts{
		Transports: []network.Transport{trLocal},
	}
	server := network.NewServer(opts)

	server.Start()
}

func sendTransaction(tr network.Transport, to network.NetAddr) error {
	privateKey := crypto.GeneratePrivateKey()
	data := []byte(strconv.FormatInt(rand.Int63(), 10))
	tx := core.NewTransaction(data)
	tx.Sign(privateKey)

	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())

	return tr.SendMessage(to, msg.Bytes())
}
