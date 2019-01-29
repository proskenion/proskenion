package repository_test

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/repository"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func test_TxHistory_Upserts(t *testing.T, txHistory core.TxHistory, txList core.TxList) {
	hash := txList.Hash()
	_, err := txHistory.GetTxList(hash)
	require.EqualError(t, errors.Cause(err), core.ErrTxHistoryNotFound.Error())
	err = txHistory.Append(txList)
	require.NoError(t, err)

	retTxList, err := txHistory.GetTxList(hash)
	require.NoError(t, err)
	assert.Equal(t, txList.Hash(), retTxList.Hash())
}

func test_TxHistory(t *testing.T, dba core.DBA, TxHistory core.TxHistory) {
	txs := []core.TxList{
		RandomTxList(),
		RandomTxList(),
		RandomTxList(),
		RandomTxList(),
		RandomTxList(),
	}
	txs2 := []core.TxList{
		RandomTxList(),
		RandomTxList(),
		RandomTxList(),
		RandomTxList(),
		RandomTxList(),
	}
	txs3 := []core.TxList{
		RandomTxList(),
		RandomTxList(),
		RandomTxList(),
		RandomTxList(),
		RandomTxList(),
	}

	for _, tx := range txs {
		test_TxHistory_Upserts(t, TxHistory, tx)
	}
	firstHash := TxHistory.Hash()
	require.NoError(t, TxHistory.Commit())

	for _, tx := range txs2 {
		test_TxHistory_Upserts(t, TxHistory, tx)
	}
	secondHash := TxHistory.Hash()
	require.NoError(t, TxHistory.Commit())

	tx, err := dba.Begin()
	require.NoError(t, err)
	txFirstHistory, err := repository.NewTxHistory(tx, RandomFactory(), RandomCryptor(), firstHash)

	for _, tx := range txs2 {
		test_TxHistory_Upserts(t, txFirstHistory, tx)
	}
	require.NoError(t, err)
	require.NoError(t, txFirstHistory.Commit())
	secondHash2 := txFirstHistory.Hash()
	assert.Equal(t, secondHash, secondHash2)

	txFirstHistory2, err := repository.NewTxHistory(tx, RandomFactory(), RandomCryptor(), firstHash)
	for _, tx := range txs3 {
		test_TxHistory_Upserts(t, txFirstHistory2, tx)
	}
	require.NoError(t, err)
	require.NoError(t, txFirstHistory2.Commit())
	secondHash3 := txFirstHistory2.Hash()
	assert.NotEqual(t, secondHash, secondHash3)

	for _, tx := range txs2 {
		// History2 では tx2 が保存されていないことになっている
		_, err := txFirstHistory2.GetTxList(tx.Hash())
		assert.EqualError(t, errors.Cause(err), core.ErrTxHistoryNotFound.Error())

		retTx, err := txFirstHistory.GetTxList(tx.Hash())
		assert.Equal(t, tx.Hash(), retTx.Hash())
	}
}

func TestTxHistory(t *testing.T) {
	dba := RandomDBA()
	tx, err := dba.Begin()
	require.NoError(t, err)
	TxHistory, err := repository.NewTxHistory(tx, RandomFactory(), RandomCryptor(), model.Hash(nil))
	require.NoError(t, err)
	test_TxHistory(t, dba, TxHistory)
}
