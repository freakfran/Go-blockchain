package crypto

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGeneratePrivateKey(t *testing.T) {
	privateKey := GeneratePrivateKey()
	publicKey := privateKey.PublicKey()
	address := publicKey.Address()
	fmt.Println("address:", address)
	msg := []byte("hello")
	sign, err := privateKey.Sign(msg)
	assert.Nil(t, err)
	assert.True(t, sign.Verify(publicKey, msg))

}