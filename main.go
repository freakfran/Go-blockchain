package main

import (
	"MyChain/network"
	"time"
)

func main() {
	trLocal := network.NewLocalTransport("LOCAL")
	trRemote := network.NewLocalTransport("REMOTE")

	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	go func() {
		for {
			trRemote.SendMessage(trLocal.Addr(), []byte("hello"))
			time.Sleep(3 * time.Second)
		}

	}()
	opts := network.ServerOpts{
		Transports: []network.Transport{trLocal},
	}
	server := network.NewServer(opts)

	server.Start()
}
