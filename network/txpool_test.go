package network

import (
	"MyChain/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTxPool(t *testing.T) {
	p := NewTxPool()
	assert.Equal(t, 0, p.Len())
}

func TestTxPool_Add(t *testing.T) {
	p := NewTxPool()
	tx := core.NewTransaction([]byte("foo"))
	assert.Nil(t, p.Add(tx))
	assert.Equal(t, 1, p.Len())

	p.Flush()
	assert.Equal(t, 0, p.Len())
}
