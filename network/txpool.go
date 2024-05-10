package network

import (
	"MyChain/core"
	"MyChain/types"
	"sort"
)

type TxMapSorter struct {
	transactions []*core.Transaction
}

func (s *TxMapSorter) Len() int {
	return len(s.transactions)
}

func (s *TxMapSorter) Less(i, j int) bool {
	return s.transactions[i].FirstSeen() < s.transactions[j].FirstSeen()
}

func (s *TxMapSorter) Swap(i, j int) {
	s.transactions[i], s.transactions[j] = s.transactions[j], s.transactions[i]
}

func NewTxMapSorter(txMap map[types.Hash]*core.Transaction) *TxMapSorter {
	txs := make([]*core.Transaction, len(txMap))
	i := 0
	for _, tx := range txMap {
		txs[i] = tx
		i++
	}
	s := &TxMapSorter{
		transactions: txs,
	}
	sort.Sort(s)
	return s
}

type TxPool struct {
	transactions map[types.Hash]*core.Transaction
}

func NewTxPool() *TxPool {
	return &TxPool{
		transactions: make(map[types.Hash]*core.Transaction),
	}
}

func (p *TxPool) Transactions() []*core.Transaction {
	s := NewTxMapSorter(p.transactions)
	return s.transactions
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
// 如果该交易已经存在于交易池中，调用方需要确保交易没有已经存在于交易池中
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
	// 将新交易添加到交易池
	p.transactions[hash] = tx
	return nil
}
