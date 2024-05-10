package network

import (
	"MyChain/core"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
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

	txx := core.NewTransaction([]byte("foo"))
	assert.Nil(t, p.Add(txx))
	assert.Equal(t, 1, p.Len())

	p.Flush()
	assert.Equal(t, 0, p.Len())
}

func TestNewTxMapSorter(t *testing.T) {
	p := NewTxPool()
	txLen := 1000
	for i := 0; i < txLen; i++ {
		tx := core.NewTransaction([]byte(strconv.FormatInt(int64(i), 10)))
		tx.SetFirstSeen(rand.Int63())
		assert.Nil(t, p.Add(tx))
	}
	assert.Equal(t, txLen, p.Len())
	transactions := p.Transactions()
	for i := 0; i < len(transactions)-1; i++ {
		assert.True(t, transactions[i].FirstSeen() < transactions[i+1].FirstSeen())
	}
}
