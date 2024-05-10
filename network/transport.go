package network

type NetAddr string

// Transport
//
//	@Description: 传输
type Transport interface {
	// Consume
	//  Consume(RPC)
	//  @Description:从RPC（Remote Procedure Call，远程过程调用）到RPC的通道
	Consume() <-chan RPC

	// Connect
	//  Connect(Transport)
	//  @Description: 连接到另一个Transport实例
	Connect(Transport) error

	// SendMessage
	//  SendMessage(NetAddr, []byte)
	//  @Description: 发送一个字节切片（通常代表消息或数据）到指定的NetAddr（网络地址）
	SendMessage(NetAddr, []byte) error

	// Broadcast
	//  Broadcast([]byte)
	//  @Description:广播一个字节切片（通常代表消息或数据）到所有连接的Transport实例
	Broadcast([]byte) error

	// Addr
	//  Addr()
	//  @Description:当前Transport的NetAddr
	Addr() NetAddr
}
