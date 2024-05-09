package core

import "fmt"

type Validator interface {
	ValidateBlock(block *Block) error
}

type BlockValidator struct {
	bc *Blockchain
}

func NewBlockValidator(bc *Blockchain) *BlockValidator {
	return &BlockValidator{bc: bc}
}

// ValidateBlock 验证给定的Block是否有效。如果Block已经在区块链中存在，或者Block验证失败，将返回错误。
// 参数:
//
//	block: 需要验证的Block对象。
//
// 返回:
//
//	error: 如果Block已存在于区块链（根据Height判断），或者Block验证过程出错，返回一个错误；否则返回nil。
func (v *BlockValidator) ValidateBlock(block *Block) error {
	// 检查区块链中是否已经存在该Block的高度
	if v.bc.HasBlock(block.Height) {
		return fmt.Errorf("block already exists block with height %d", block.Height)
	}
	// 验证Block本身的有效性
	if err := block.Verify(); err != nil {
		return err
	}
	// 如果一切正常，返回nil
	return nil
}
