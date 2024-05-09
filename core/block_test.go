package core

import (
	"MyChain/crypto"
	"MyChain/types"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func randomBlock(height uint32, prevBlockHash types.Hash) *Block {
	header := &Header{
		Version:       0,
		PrevBlockHash: prevBlockHash,
		Timestamp:     time.Now().UnixNano(),
		Height:        height,
	}

	return NewBlock(header, []Transaction{})
}

func randomBlockWithSignature(t *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	privateKey := crypto.GeneratePrivateKey()
	b := randomBlock(height, prevBlockHash)
	assert.Nil(t, b.Sign(privateKey))
	b.AddTransaction(randomTxWithSignature(t))
	return b
}

func TestBlock_Hash(t *testing.T) {
	block := randomBlock(0, types.Hash{})
	fmt.Println(block.Hash(BlockHasher{}))
}

func TestBlock_Sign(t *testing.T) {
	privateKey := crypto.GeneratePrivateKey()
	block := randomBlock(0, types.Hash{})

	assert.Nil(t, block.Sign(privateKey))
	assert.Nil(t, block.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	block.Validator = otherPrivKey.PublicKey()
	assert.NotNil(t, block.Verify())
}
