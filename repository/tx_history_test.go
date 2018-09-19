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

func test_TxHistory_Upserts(t *testing.T, TxHistory core.TxHistory, tx model.Transaction) {
	hash := MustHash(tx)
	_, err := TxHistory.Query(hash)
	require.EqualError(t, errors.Cause(err), core.ErrTxHistoryNotFound.Error())
	err = TxHistory.Append(tx)
	require.NoError(t, err)

	retTx, err := TxHistory.Query(hash)
	require.NoError(t, err)
	assert.Equal(t, MustHash(tx), MustHash(retTx))
}

func test_TxHistory(t *testing.T, TxHistory core.TxHistory) {
	txs := []model.Transaction{
		RandomTx(),
		RandomTx(),
		RandomTx(),
		RandomTx(),
		RandomTx(),
	}

	for _, tx := range txs {
		test_TxHistory_Upserts(t, TxHistory, tx)
	}
	require.NoError(t, TxHistory.Commit())
}

func TestTxHistory(t *testing.T) {
	TxHistory, err := repository.NewTxHistory(RandomDBATx(), NewTestFactory(), RandomCryptor(), model.Hash(nil))
	require.NoError(t, err)
	test_TxHistory(t, TxHistory)
}
