package core

import (
	"MyChain/crypto"
	"MyChain/types"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func randomBlock(height uint32) *Block {
	header := &Header{
		Version:       0,
		PrevBlockHash: types.Hash{},
		Timestamp:     time.Now().UnixNano(),
		Height:        height,
	}
	tx := Transaction{
		Data: []byte("hello"),
	}
	return NewBlock(header, []Transaction{tx})
}

func randomBlockWithSignature(t *testing.T, height uint32) *Block {
	privateKey := crypto.GeneratePrivateKey()
	b := randomBlock(height)
	assert.Nil(t, b.Sign(privateKey))
	return b
}

func TestBlock_Hash(t *testing.T) {
	block := randomBlock(0)
	fmt.Println(block.Hash(BlockHasher{}))
}

func TestBlock_Sign(t *testing.T) {
	privateKey := crypto.GeneratePrivateKey()
	block := randomBlock(0)

	assert.Nil(t, block.Sign(privateKey))
	assert.Nil(t, block.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	block.Validator = otherPrivKey.PublicKey()
	assert.NotNil(t, block.Verify())
}
