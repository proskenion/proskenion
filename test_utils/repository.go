package test_utils

import (
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
	config := RandomConfig()
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

func RandomCommitableBlockAndTxList(t *testing.T, rp core.Repository) (model.Block, core.TxList) {
	block, txList, err := rp.CreateBlock(RandomQueue(), 0, RandomNow())
	require.NoError(t, err)
	return block, txList
}

func EmptyTxList() core.TxList {
	return repository.NewTxList(RandomCryptor(), RandomFactory())
}

func RandomTxList() core.TxList {
	txList := repository.NewTxList(RandomCryptor(), RandomFactory())
	for _, tx := range RandomTxs() {
		txList.Push(tx)
	}
	return txList
}

func RandomRepository() core.Repository {
	_, _, _, _, rp, _, _ := NewTestFactories()
	return rp
}

func RandomGenesisTxList(t *testing.T) core.TxList {
	NewTestFactories()
	ret, err := repository.GenesisTxListFromConf(RandomCryptor(), RandomFactory(), RandomRepository(), RandomProsl(), RandomConfig())
	require.NoError(t, err)
	return ret
}

func RandomTxListCache() core.TxListCache {
	return repository.NewTxListCache(RandomConfig())
}

func MusTop(rp core.Repository) model.Block {
	top, _ := rp.Top()
	return top
}
