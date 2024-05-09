package network

import (
	"MyChain/core"
	"MyChain/types"
)

type TxPool struct {
	transactions map[types.Hash]*core.Transaction
}

func NewTxPool() *TxPool {
	return &TxPool{
		transactions: make(map[types.Hash]*core.Transaction),
	}
}

func (p *TxPool) Len() int {
	return len(p.transactions)
}

func (p *TxPool) Flush() {
	p.transactions = make(map[types.Hash]*core.Transaction)
}

// Has 检查给定哈希的交易是否存在于交易池中。
// 如果找到该交易，则返回true，否则返回false。
//
// 参数：
// - hash：要检查的交易的唯一标识符。
//
// 返回：
// - 一个布尔值，表示交易是否在池中存在。
func (p *TxPool) Has(hash types.Hash) bool {
	// 通过查找其哈希来检查交易是否存在于池中。
	_, ok := p.transactions[hash]
	return ok
}

// Add 将一个交易添加到交易池中。
// 如果该交易已经存在于交易池中，则不进行任何操作。
// 参数：
//
//	tx *core.Transaction: 需要添加的交易。
//
// 返回值：
//
//	error: 如果添加过程中遇到错误，则返回错误信息；否则返回nil。
func (p *TxPool) Add(tx *core.Transaction) error {
	// 生成交易的哈希值
	hash := tx.Hash(core.TxHasher{})
	// 检查交易是否已经存在于交易池中
	if p.Has(hash) {
		return nil
	}
	// 将新交易添加到交易池
	p.transactions[hash] = tx
	return nil
}
