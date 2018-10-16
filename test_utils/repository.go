package test_utils

import (
	"github.com/proskenion/proskenion/commit"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/repository"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

type MockKVNode struct {
	key     []byte
	account model.Account
}

func RandomKVStoreFromAccount(key []byte, ac model.Account) core.KVNode {
	return &MockKVNode{key, ac}
}

func (kv *MockKVNode) Key() []byte {
	return kv.key
}

func (kv *MockKVNode) Value() model.Marshaler {
	return kv.account
}

func (kv *MockKVNode) Next(cnt int) core.KVNode {
	return &MockKVNode{
		kv.key[cnt:],
		kv.account,
	}
}

var MOCK_ROOT_KEY byte = 0

func RandomStrKey() []byte {
	ret := make([]byte, rand.Int()%10+2)
	ret[0] = MOCK_ROOT_KEY
	for i := 1; i < len(ret); i++ {
		ret[i] = byte(rand.Intn(26))
	}
	return ret
}

func RandomQueue() core.ProposalTxQueue {
	config := NewTestConfig()
	queue := repository.NewProposalTxQueueOnMemory(config)
	for i := 0; i < 100; i++ {
		tx := RandomValidTx()
		err := queue.Push(tx)
		if err != nil {
			panic(err)
		}
	}
	return queue
}

func RandomTxList() core.TxList {
	txList := repository.NewTxList(RandomCryptor())
	for _, tx := range RandomTxs() {
		txList.Push(tx)
	}
	return txList
}

func RandomCommitableBlock(t *testing.T, top model.Block, rp core.Repository) (model.Block, core.TxList) {
	wsvHash := model.Hash(nil)
	txHistoryHash := model.Hash(nil)
	topHash := model.Hash(nil)
	topHeight := int64(0)
	if top != nil {
		wsvHash = top.GetPayload().GetWSVHash()
		txHistoryHash = top.GetPayload().GetTxHistoryHash()
		topHash = MustHash(top)
		topHeight = top.GetPayload().GetHeight()
	}

	dtx, err := rp.Begin()
	require.NoError(t, err)

	// load state
	bc, err := dtx.Blockchain(topHash)
	require.NoError(t, err)

	wsv, err := dtx.WSV(wsvHash)
	require.NoError(t, err)

	txHistory, err := dtx.TxHistory(txHistoryHash)
	require.NoError(t, err)

	txList := repository.NewTxList(RandomCryptor())

	for txList.Size() < 100 {
		tx := NewTestFactory().NewTxBuilder().
			CreateAccount("authorizer@com", RandomStr()+"@com").
			CreatedTime(RandomNow()).Build()
		// tx を構築
		for _, cmd := range tx.GetPayload().GetCommands() {
			err = cmd.Validate(wsv)
			require.NoError(t, err)

			err = cmd.Execute(wsv)
			require.NoError(t, err)
		}
		err = txHistory.Append(tx)
		require.NoError(t, err)
		txList.Push(tx)
	}

	newTxHistoryHash := MustHash(txHistory)
	newWSVHash := MustHash(wsv)

	newBlock := NewTestFactory().NewBlockBuilder().
		Round(0).
		TxsHash(txList.Top()).
		TxHistoryHash(newTxHistoryHash).
		WSVHash(newWSVHash).
		CreatedTime(commit.Now()).
		Height(topHeight + 1).
		PreBlockHash(topHash).
		Build()

	pub, pri := RandomCryptor().NewKeyPairs()
	err = newBlock.Sign(pub, pri)
	require.NoError(t, err)

	err = bc.Append(newBlock)
	require.NoError(t, err)
	return newBlock, txList
}
