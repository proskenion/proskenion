package repository_test

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/crypto"
	. "github.com/proskenion/proskenion/repository"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func testTxList_PushAndTop(t *testing.T, list core.TxList) {
	txs := []model.Transaction{
		RandomTx(),
		RandomTx(),
		RandomTx(),
		RandomTx(),
		RandomTx(),
	}

	assert.Equal(t, model.Hash(nil), list.Hash())
	for _, tx := range txs {
		err := list.Push(tx)
		require.NoError(t, err)
	}
	for i, tx := range list.List() {
		assert.Equal(t, txs[i], tx)
	}
}

func TestTxList_PushAndTop(t *testing.T) {
	cryptor := crypto.NewEd25519Sha256Cryptor()
	fc := RandomFactory()
	testTxList_PushAndTop(t, NewTxList(cryptor, fc))
}

func TestTxList_MarshalUnmarshal(t *testing.T) {
	txList := RandomTxList()
	fc := RandomFactory()
	anothertxList := NewTxList(RandomCryptor(), fc)

	b, err := txList.Marshal()
	require.NoError(t, err)
	require.NoError(t, anothertxList.Unmarshal(b))

	assert.Equal(t, txList.Hash(), anothertxList.Hash())

	assert.Equal(t, len(txList.List()), len(anothertxList.List()))
	for i, tx := range txList.List() {
		assert.Equal(t, tx.Hash(), anothertxList.List()[i].Hash())
	}
}
