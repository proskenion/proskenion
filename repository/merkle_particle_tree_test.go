package repository_test

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/dba"
	"github.com/proskenion/proskenion/repository"
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

var MOCK_ROOT_KEY byte = 0

func RandomStrKey() []byte {
	str := RandomStr()
	ret := make([]byte, 1)
	ret[0] = MOCK_ROOT_KEY
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

	_, err = tree.Upsert(node)
	require.NoError(t, err)
	it, err = tree.Find(node.Key())
	err = it.Data(unmarshaler)
	assert.NoError(t, err)
	// use value is model.Account
	assert.Equal(t, node.Value().(model.Account), unmarshaler.(model.Account))
}

func testMerklePatriciaTree(t *testing.T, tree1 core.MerklePatriciaTree, tree2 core.MerklePatriciaTree) {
	acs := []model.Account{
		RandomAccount(),
		RandomAccount(),
		RandomAccount(),
		RandomAccount(),
		RandomAccount(),
	}
	acs2 := []model.Account{
		RandomAccount(),
		RandomAccount(),
		RandomAccount(),
		RandomAccount(),
		RandomAccount(),
	}
	keys := [][]byte{
		RandomStrKey(),
		RandomStrKey(),
		RandomStrKey(),
		RandomStrKey(),
		RandomStrKey(),
	}
	unmarshaler := RandomAccount()

	// First Upsert tree1 and tree2
	for i, key := range keys {
		kvNode := RandomKVStoreFromAccount(key, acs[i])

		testUpsertFirst(t, tree1, kvNode, unmarshaler)
		testUpsertFirst(t, tree2, kvNode, unmarshaler)
		assert.Equal(t, MustHash(tree1), MustHash(tree2))
	}

	firstHash := MustHash(tree1)

	// Second Upsert tree1 and tree2
	for i, key := range keys {
		kvNode := RandomKVStoreFromAccount(key, acs2[i])
		testUpsertSecond(t, tree1, kvNode, unmarshaler)
		testUpsertSecond(t, tree2, kvNode, unmarshaler)
		assert.Equal(t, MustHash(tree1), MustHash(tree2))
	}
	lastHash := MustHash(tree1)

	// Compare first hash and last hash
	assert.NotEqual(t, firstHash, lastHash)

	// First Upsert 終了時に巻き戻し
	err := tree1.Set(firstHash)
	require.NoError(t, err)
	for i, key := range keys {
		ac := RandomAccount()
		it, err := tree1.Find(key)
		require.NoError(t, err)
		err = it.Data(ac)
		assert.NoError(t, err)
		assert.Equal(t, acs[i], ac)
	}

	// Second Upsert part 2 tree1
	err = tree1.Set(lastHash)
	require.NoError(t, err)
	for i, key := range keys {
		ac := RandomAccount()
		it, err := tree1.Find(key)
		require.NoError(t, err)
		err = it.Data(ac)
		assert.NoError(t, err)
		assert.Equal(t, acs2[i], ac)
	}
}

func TestMerklePatriciaTree(t *testing.T) {
	cryptor := RandomCryptor()
	tree1, err := repository.NewMerklePatriciaTree(dba.NewDBAOnMemory(), cryptor, model.Hash(nil), MOCK_ROOT_KEY)
	require.NoError(t, err)
	tree2, err := repository.NewMerklePatriciaTree(dba.NewDBAOnMemory(), cryptor, model.Hash(nil), MOCK_ROOT_KEY)
	require.NoError(t, err)
	testMerklePatriciaTree(t, tree1, tree2)
}
