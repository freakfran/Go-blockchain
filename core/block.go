package core

import (
	"MyChain/crypto"
	"MyChain/types"
	"bytes"
	"encoding/gob"
	"fmt"
)

type Header struct {
	Version       uint32
	DataHash      types.Hash
	PrevBlockHash types.Hash
	Timestamp     int64
	Height        uint32
}

func (h *Header) Bytes() []byte {
	// 创建一个 bytes.Buffer 用于存储序列化后的数据
	buf := &bytes.Buffer{}
	// 使用 gob Encoder 将 Block 的头部信息编码到 buf 中
	encoder := gob.NewEncoder(buf)
	encoder.Encode(h)
	// 返回编码后的字节切片
	return buf.Bytes()
}

type Block struct {
	*Header
	Transactions []Transaction
	Validator    crypto.PublicKey
	Signature    *crypto.Signature
	//cached version of the header hash
	hash types.Hash
}

func NewBlock(h *Header, txs []Transaction) *Block {
	return &Block{
		Header:       h,
		Transactions: txs,
	}

}

func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, *tx)
}

func (b *Block) Decode(dec Decoder[*Block]) error {
	return dec.Decode(b)
}
func (b *Block) Encode(enc Encoder[*Block]) error {
	return enc.Encode(b)
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
func (b *Block) Hash(hasher Hasher[*Header]) types.Hash {
	// 检查hash是否已经计算，若未计算则使用hasher进行计算
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b.Header)
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
	sign, err := privateKey.Sign(b.Header.Bytes())
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
	if !b.Signature.Verify(b.Validator, b.Header.Bytes()) {
		return fmt.Errorf("invalid signature")
	}

	for _, tx := range b.Transactions {
		if err := tx.Verify(); err != nil {
			return err
		}
	}
	return nil
}
