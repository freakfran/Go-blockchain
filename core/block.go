package core

import (
	"MyChain/crypto"
	"MyChain/types"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
)

type Header struct {
	Version       uint32
	DataHash      types.Hash
	PrevBlockHash types.Hash
	Timestamp     int64
	Height        uint32
}

type Block struct {
	Header
	Transactions []Transaction
	Validator    crypto.PublicKey
	Signature    *crypto.Signature
	//cached version of the header hash
	hash types.Hash
}

func NewBlock(h *Header, txs []Transaction) *Block {
	return &Block{
		Header:       *h,
		Transactions: txs,
	}

}

func (b *Block) Decode(reader io.Reader, dec Decoder[*Block]) error {
	return dec.Decode(reader, b)
}
func (b *Block) Encode(writer io.Writer, enc Encoder[*Block]) error {
	return enc.Encode(writer, b)
}

// Hash 计算给定Block的哈希值。
// 如果该Block的哈希值尚未计算，将使用提供的hasher进行计算。
// 参数：
//
//	hasher - 用于计算Block哈希的Hasher接口实例。
//
// 返回值：
//
//	计算得到的Block哈希值。
func (b *Block) Hash(hasher Hasher[*Block]) types.Hash {
	// 检查hash是否已经计算，若未计算则使用hasher进行计算
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b)
	}
	return b.hash
}

// Sign 对区块进行签名。
// 参数:
// - privateKey: 执行签名的私钥。
// 返回值:
// - error: 执行过程中遇到的错误。
func (b *Block) Sign(privateKey crypto.PrivateKey) error {
	// 使用私钥对区块头数据进行签名
	sign, err := privateKey.Sign(b.HeaderData())
	if err != nil {
		return err // 如果签名过程中出现错误，则返回错误
	}

	// 设置验证者公钥和签名信息
	b.Validator = privateKey.PublicKey()
	b.Signature = sign

	return nil // 完成签名过程，无错误返回
}

func (b *Block) Verify() error {
	if b.Signature == nil {
		return fmt.Errorf("block has no signature")
	}
	if !b.Signature.Verify(b.Validator, b.HeaderData()) {
		return fmt.Errorf("invalid signature")
	}
	return nil
}

// HeaderData 将 Block 的头部信息序列化为字节切片返回。
// 该函数不接受参数，返回值为序列化后的头部信息字节切片。
func (b *Block) HeaderData() []byte {
	// 创建一个 bytes.Buffer 用于存储序列化后的数据
	buf := &bytes.Buffer{}
	// 使用 gob Encoder 将 Block 的头部信息编码到 buf 中
	encoder := gob.NewEncoder(buf)
	encoder.Encode(b.Header)
	// 返回编码后的字节切片
	return buf.Bytes()
}
