package core

import (
	"MyChain/types"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newBlockChainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockChain(randomBlock(0, types.Hash{}))
	assert.Nil(t, err)
	return bc
}

func getPrevBlockHash(t *testing.T, height uint32, bc *Blockchain) types.Hash {
	prevHeader, err := bc.GetHeader(height - 1)
	assert.Nil(t, err)
	return BlockHasher{}.Hash(prevHeader)
}

func TestNewBlockChain(t *testing.T) {
	bc := newBlockChainWithGenesis(t)
	assert.NotNil(t, bc.validator)
	fmt.Println(bc.Height())
}
func TestBlockchain_HasBlock(t *testing.T) {
	bc := newBlockChainWithGenesis(t)
	assert.True(t, bc.HasBlock(0))
}
func TestBlockchain_AddBlock(t *testing.T) {
	bc := newBlockChainWithGenesis(t)
	lenBlock := 1000
	for i := 0; i < lenBlock; i++ {
		block := randomBlockWithSignature(t, uint32(i+1), getPrevBlockHash(t, uint32(i+1), bc))
		assert.Nil(t, bc.AddBlock(block))
	}
	assert.Equal(t, uint32(lenBlock), bc.Height())
	assert.Equal(t, lenBlock+1, len(bc.headers))
	assert.NotNil(t, bc.AddBlock(randomBlock(88, types.Hash{})))
}

func TestBlockchain_AddBlock_High(t *testing.T) {
	bc := newBlockChainWithGenesis(t)
	block := randomBlockWithSignature(t, 1000, types.Hash{})
	assert.NotNil(t, bc.AddBlock(block))
}

func TestBlockchain_GetHeader(t *testing.T) {
	bc := newBlockChainWithGenesis(t)
	lenBlock := 1000
	for i := 0; i < lenBlock; i++ {
		block := randomBlockWithSignature(t, uint32(i+1), getPrevBlockHash(t, uint32(i+1), bc))
		assert.Nil(t, bc.AddBlock(block))
		header, err := bc.GetHeader(uint32(i + 1))
		assert.Nil(t, err)
		assert.Equal(t, block.Header, header)
	}
}
