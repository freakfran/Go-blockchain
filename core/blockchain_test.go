package core

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newBlockChainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockChain(randomBlock(0))
	assert.Nil(t, err)
	return bc
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
		block := randomBlockWithSignature(t, uint32(i+1))
		assert.Nil(t, bc.AddBlock(block))
	}
	assert.Equal(t, uint32(lenBlock), bc.Height())
	assert.Equal(t, lenBlock+1, len(bc.headers))
	assert.NotNil(t, bc.AddBlock(randomBlock(88)))
}
