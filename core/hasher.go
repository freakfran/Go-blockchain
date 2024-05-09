package core

import (
	"MyChain/types"
	"crypto/sha256"
)

// Hasher 是一个泛型接口，用于定义对任意类型 T 进行哈希处理的能力。
// 其中，T 代表可以被哈希处理的任意类型。
// 该接口包含一个方法 Hash，用于计算给定输入 T 的哈希值。
type Hasher[T any] interface {
	// Hash 方法接受一个类型为 T 的参数 value，对其进行哈希处理，
	// 并返回一个 types.Hash 类型的哈希值。
	Hash(T) types.Hash
}

type BlockHasher struct {
}

// Hash 计算给定Block的哈希值。
// 参数：
// b *Block - 需要计算哈希值的区块。
// 返回值：
// types.Hash - 计算得到的哈希值。
func (BlockHasher) Hash(h *Header) types.Hash {
	// 计算字节缓冲区内容的SHA256哈希值
	b := sha256.Sum256(h.Bytes())
	// 返回计算得到的哈希值
	return types.Hash(b)
}

type TxHasher struct {
}

func (TxHasher) Hash(tx *Transaction) types.Hash {
	return types.Hash(sha256.Sum256(tx.Data))
}
