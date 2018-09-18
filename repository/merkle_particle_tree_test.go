package repository_test

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	kv.key = kv.key[cnt:]
	return kv
}

func RandomStrKey() []byte {
	str := RandomStr()
	ret := make([]byte, 0)
	for _, c := range str {
		ret = append(ret, byte(c-'a'))
	}
	return ret
}

func testUpsertFirst(t *testing.T, tree core.MerklePatriciaTree, node core.KVNode, unmarshaler model.Unmarshaler) {
	_, err := tree.Find(node.Key())
	assert.EqualError(t, errors.Cause(err), core.ErrMerklePatriciaTreeNotFoundKey.Error())

	_, err = tree.Upsert(node)
	require.NoError(t, err)

	it, err := tree.Find(node.Key())
	assert.NoError(t, err)
	err = it.Data(unmarshaler)
	assert.NoError(t, err)
	// use value is model.Account
	assert.Equal(t, node.Value().(model.Account), unmarshaler.(model.Account))
}

func testUpsertSecond(t *testing.T, tree core.MerklePatriciaTree, node core.KVNode, unmarshaler model.Unmarshaler) {
	it, err := tree.Find(node.Key())
	assert.NoError(t, err)
	err = it.Data(unmarshaler)
	assert.NoError(t, err)
	// use value is model.Account
	assert.Equal(t, node.Value().(model.Account), unmarshaler.(model.Account))

	_, err = tree.Upsert(node)
	require.NoError(t, err)
	it, err = tree.Find(node.Key())
	err = it.Data(unmarshaler)
	assert.NoError(t, err)
	// use value is model.Account
	assert.Equal(t, node.Value().(model.Account), unmarshaler.(model.Account))
}

// TODO
func testMerklePatriciaTree(t *testing.T, tree1 core.MerklePatriciaTree, tree2 core.MerklePatriciaTree) {
	acs := []model.Account{
		RandomAccount(),
		RandomAccount(),
		RandomAccount(),
		RandomAccount(),
		RandomAccount(),
	}
	nodes := []core.KVNode{
		RandomKVStoreFromAccount(RandomStrKey(), acs[0]),
		RandomKVStoreFromAccount(RandomStrKey(), acs[1]),
		RandomKVStoreFromAccount(RandomStrKey(), acs[2]),
		RandomKVStoreFromAccount(RandomStrKey(), acs[3]),
		RandomKVStoreFromAccount(RandomStrKey(), acs[4]),
	}
	unmarshalers := []model.Account{
		RandomAccount(),
		RandomAccount(),
		RandomAccount(),
		RandomAccount(),
		RandomAccount(),
	}

	// First Upsert tree1 and tree2
	for i, node := range nodes {
		testUpsertFirst(t, tree1, node, unmarshalers[i])
		testUpsertFirst(t, tree2, node, unmarshalers[i])
		assert.Equal(t, MustHash(tree1), MustHash(tree2))
	}
	firstHash := MustHash(tree1)

	// Second Upsert tree1 and tree2
	for i, node := range nodes {
		testUpsertSecond(t, tree1, node, unmarshalers[i])
		testUpsertSecond(t, tree2, node, unmarshalers[i])
		assert.Equal(t, MustHash(tree1), MustHash(tree2))
	}
	lastHash := MustHash(tree1)

	// Compare first hash and last hash
	assert.NotEqual(t, firstHash, lastHash)

	// First Upsert 終了時に巻き戻し
	actFirstNode := tree1.Iterator()
	actFirstHash, err := actFirstNode.Hash()
	require.NoError(t, err)
	assert.Equal(t, firstHash, actFirstHash)
	for _, node := range nodes {
		v := RandomMarshaler()
		err := actFirstNode.Find(node, v)
		require.NoError(t, err)
		assert.Equal(t, node, v)
	}
	// Second Upsert part 2 tree1
	for i, node := range nodes {
		testUpsertSecond(t, tree1, node, nodes3[i])
	}
}
