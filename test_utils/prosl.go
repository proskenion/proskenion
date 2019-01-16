package test_utils

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/repository"
	"github.com/stretchr/testify/require"
	"testing"
)

func CommitTxWrapBlock(t *testing.T, rp core.Repository, fc model.ModelFactory, tx model.Transaction) {
	dtx, err := rp.Begin()
	require.NoError(t, err)

	top, ok := rp.Top()
	topHash := model.Hash(nil)
	if !ok {
		top = fc.NewBlockBuilder().
			PreBlockHash(nil).
			TxHistoryHash(nil).
			WSVHash(nil).
			Height(-1).
			Build()
	} else {
		topHash = top.Hash()
	}

	// load state
	var bc core.Blockchain
	bc, err = dtx.Blockchain(top.GetPayload().GetPreBlockHash())
	require.NoError(t, err)
	wsv, err := dtx.WSV(top.GetPayload().GetWSVHash())
	require.NoError(t, err)
	txHistory, err := dtx.TxHistory(top.GetPayload().GetTxHistoryHash())
	require.NoError(t, err)

	// transactions execute (no validate)
	txList := repository.NewTxList(RandomCryptor())
	require.NoError(t, txList.Push(tx))
	for _, cmd := range tx.GetPayload().GetCommands() {
		err := cmd.Execute(wsv)
		require.NoError(t, err)
	}
	err = txHistory.Append(tx)
	require.NoError(t, err)

	// hash check and block 生成
	wsvHash := wsv.Hash()
	require.NoError(t, err)
	txHistoryHash := txHistory.Hash()
	require.NoError(t, err)
	block := fc.NewBlockBuilder().
		CreatedTime(RandomNow()).
		TxsHash(txList.Top()).
		PreBlockHash(topHash).
		TxHistoryHash(txHistoryHash).
		WSVHash(wsvHash).
		Round(0).
		Height(top.GetPayload().GetHeight() + 1).
		Build()

	// block を追加・
	err = bc.Append(block)
	require.NoError(t, err)

	err = dtx.Commit()
	require.NoError(t, err)

	// repository の Top を Height を返す。
	rp.(*repository.Repository).TopBlock = block
	rp.(*repository.Repository).Height = block.GetPayload().GetHeight()
}
