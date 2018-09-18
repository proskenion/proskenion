package test_utils

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"math/rand"
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
