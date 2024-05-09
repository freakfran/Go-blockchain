package core

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
)

type Blockchain struct {
	store     Storage
	lock      sync.RWMutex
	headers   []*Header
	validator Validator
}

// NewBlockChain 创建一个新的区块链实例。
//
// 参数:
//
//	genesis *Block: 用于初始化区块链的创世区块。
//
// 返回值:
//
//	*Blockchain: 初始化后的区块链实例。
//	error: 如果在初始化过程中遇到错误，则返回错误信息；否则返回nil。
func NewBlockChain(genesis *Block) (*Blockchain, error) {
	// 初始化Blockchain结构体，包括空的区块头切片和一个新的内存存储实例
	bc := &Blockchain{
		headers: []*Header{},
		store:   NewMemoryStorage(),
	}
	// 尝试添加创世区块，不进行验证
	err := bc.addBlockWithoutValidation(genesis)
	if err != nil {
		return nil, err // 如果添加创世区块失败，则返回错误
	}
	// 为区块链实例设置区块验证器
	bc.validator = NewBlockValidator(bc)
	return bc, nil
}

// addBlockWithoutValidation 方法用于将一个区块添加到区块链中，但不进行验证。
// 此方法直接将区块头添加到区块链的头部列表，并通过存储接口将区块存储起来。
// 参数:
//
//	b *Block - 需要被添加到区块链的区块。
//
// 返回值:
//
//	error - 添加过程中遇到的错误，如果没有错误则为 nil。
func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {
	bc.lock.Lock()
	// 将新区块的头添加到区块链的头列表中
	bc.headers = append(bc.headers, b.Header)
	bc.lock.Unlock()

	logrus.WithFields(logrus.Fields{
		"height": b.Height,
		"hash":   b.Hash(BlockHasher{}),
	}).Infoln("add a new block")
	// 将区块存储起来，返回可能发生的错误
	return bc.store.Put(b)
}

func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

// HasBlock 检查区块链中是否存在指定高度的区块
// 参数：
//
//	height - 指定的区块高度
//
// 返回值：
//
//	bool - 如果区块链中存在该高度的区块则返回true，否则返回false
func (bc *Blockchain) HasBlock(height uint32) bool {
	// 判断指定的区块高度是否小于当前区块链的高度
	return height <= bc.Height()
}

// AddBlock 将一个新的区块添加到区块链中。
// 此函数首先会验证区块的有效性，如果验证失败，则返回相应的错误。
// 如果验证成功，则会调用内部函数 addBlockWithoutValidation 来实际添加区块。
//
// 参数:
//
//	b *Block - 需要添加到区块链的区块。
//
// 返回值:
//
//	error - 如果验证失败或其他原因导致添加失败，返回错误信息；否则返回 nil。
func (bc *Blockchain) AddBlock(b *Block) error {
	// 验证区块的有效性
	err := bc.validator.ValidateBlock(b)
	if err != nil {
		return err // 验证失败，返回错误
	}
	return bc.addBlockWithoutValidation(b) // 验证成功，添加区块
}

// Height 返回Blockchain当前的高度，即区块头的数量减一。
// 该函数不接受参数。
// 返回值：
// uint32：Blockchain的高度。
func (bc *Blockchain) Height() uint32 {
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	// 计算并返回Blockchain的高度
	return uint32(len(bc.headers) - 1)
}

func (bc *Blockchain) GetHeader(height uint32) (*Header, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("blockchain height is %d, but get %d", bc.Height(), height)
	}
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	return bc.headers[height], nil
}
