package core

import (
	"MyChain/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransaction_Sign(t *testing.T) {
	privateKey := crypto.GeneratePrivateKey()
	msg := []byte("test")
	tx := &Transaction{
		Data: msg,
	}
	assert.Nil(t, tx.Sign(privateKey))
	assert.NotNil(t, tx.Signature)
}

func TestTransaction_Verify(t *testing.T) {
	privateKey := crypto.GeneratePrivateKey()
	msg := []byte("test")
	tx := &Transaction{
		Data: msg,
	}

	assert.Nil(t, tx.Sign(privateKey))
	assert.Nil(t, tx.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	tx.PublicKey = otherPrivKey.PublicKey()
	assert.NotNil(t, tx.Verify())
}
