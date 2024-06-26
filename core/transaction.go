package core

import (
	"MyChain/crypto"
	"MyChain/types"
	"fmt"
)

type Transaction struct {
	Data      []byte
	From      crypto.PublicKey
	Signature *crypto.Signature

	//cached
	hash types.Hash
	//firstSeen is the timestamp when the transaction is first seen locally
	firstSeen int64
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data: data,
	}
}

func (tx *Transaction) Hash(hasher Hasher[*Transaction]) types.Hash {
	if tx.hash.IsZero() {
		tx.hash = hasher.Hash(tx)
	}
	return tx.hash
}

// Sign 对交易数据进行签名。
//
// 参数:
// - privateKey: 执行签名的私钥。
//
// 返回值:
// - error: 执行过程中遇到的错误，如果签名成功则为nil。
func (tx *Transaction) Sign(privateKey crypto.PrivateKey) error {
	// 使用私钥对交易数据进行签名
	sign, err := privateKey.Sign(tx.Data)
	if err != nil {
		return err // 返回签名过程中遇到的任何错误
	}

	// 设置公钥和签名值
	tx.From = privateKey.PublicKey()
	tx.Signature = sign

	return nil // 成功完成签名过程，返回nil
}

// Verify 验证交易的签名有效性。
// 如果交易签名为空，返回一个错误。
// 如果签名无法通过公钥和交易数据验证，返回一个错误。
// 若验证成功，返回 nil。
func (tx *Transaction) Verify() error {
	// 检查交易签名是否为空
	if tx.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}

	// 验证签名，如果无效则返回错误
	if !tx.Signature.Verify(tx.From, tx.Data) {
		return fmt.Errorf("invalid signature")
	}

	// 验证成功，返回 nil
	return nil
}

func (tx *Transaction) SetFirstSeen(firstSeen int64) {
	tx.firstSeen = firstSeen
}

func (tx *Transaction) FirstSeen() int64 {
	return tx.firstSeen
}

func (tx *Transaction) Decode(dec Decoder[*Transaction]) error {
	return dec.Decode(tx)
}

func (tx *Transaction) Encode(enc Encoder[*Transaction]) error {
	return enc.Encode(tx)
}
