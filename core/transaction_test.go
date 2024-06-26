package core

import (
	"MyChain/crypto"
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func randomTxWithSignature(t *testing.T) *Transaction {
	privateKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: []byte("foo"),
	}
	assert.Nil(t, tx.Sign(privateKey))
	return tx
}

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
	tx.From = otherPrivKey.PublicKey()
	assert.NotNil(t, tx.Verify())
}

func TestTransaction_Encode_Decode(t *testing.T) {
	tx := randomTxWithSignature(t)
	buf := bytes.Buffer{}
	assert.Nil(t, tx.Encode(NewGobTxEncoder(&buf)))

	dec := new(Transaction)
	assert.Nil(t, dec.Decode(NewGobTxDecoder(&buf)))

	assert.Equal(t, tx, dec)
}
