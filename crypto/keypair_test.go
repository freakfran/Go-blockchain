package crypto

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeypair_Sign_Verify_Success(t *testing.T) {
	privateKey := GeneratePrivateKey()
	publicKey := privateKey.PublicKey()
	address := publicKey.Address()
	fmt.Println("address:", address)
	msg := []byte("hello")
	sign, err := privateKey.Sign(msg)
	assert.Nil(t, err)
	assert.True(t, sign.Verify(publicKey, msg))
}

func TestKeypair_Sign_Verify_Fail(t *testing.T) {
	privateKey := GeneratePrivateKey()
	publicKey := privateKey.PublicKey()

	msg := []byte("hello")
	sign, err := privateKey.Sign(msg)
	assert.Nil(t, err)

	otherPrivateKey := GeneratePrivateKey()
	otherPubKey := otherPrivateKey.PublicKey()
	assert.False(t, sign.Verify(otherPubKey, msg))
	assert.False(t, sign.Verify(publicKey, []byte("no")))
}
