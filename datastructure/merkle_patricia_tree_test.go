package datastructure_test

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/datastructure"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func testUpsertFirst(t *testing.T, tree core.MerklePatriciaTree, node core.KVNode, unmarshaler model.Unmarshaler) {
	_, err := tree.Find(node.Key())
	assert.EqualError(t, errors.Cause(err), core.ErrMerklePatriciaTreeNotFoundKey.Error())

	_, err = tree.Upsert(node)
	require.NoError(t, err)

	it, err := tree.Find(node.Key())
	assert.NoError(t, err)
	assert.True(t, it.Leaf())
	err = it.Data(unmarshaler)
	assert.NoError(t, err)
	// use value is model.Account
	assert.Equal(t, node.Value().(model.Account).Hash(), unmarshaler.(model.Account).Hash())

	// search key
	iti, err := tree.Search(node.Key())
	assert.Equal(t, iti.DataHash(), it.Hash())
}

func testUpsertSecond(t *testing.T, tree core.MerklePatriciaTree, node core.KVNode, unmarshaler model.Unmarshaler) {
	it, err := tree.Find(node.Key())
	assert.NoError(t, err)

	_, err = tree.Upsert(node)
	require.NoError(t, err)
	it, err = tree.Find(node.Key())
	assert.True(t, it.Leaf())
	err = it.Data(unmarshaler)
	assert.NoError(t, err)
	// use value is model.Account
	assert.Equal(t, MustHash(node.Value().(model.Account)), MustHash(unmarshaler.(model.Account)))

	// search key
	iti, err := tree.Search(node.Key())
	assert.Equal(t, iti.DataHash(), it.Hash())
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
		assert.Equal(t, MustHash(acs[i]), MustHash(ac))
	}
	// Check SubTree1
	it := tree1.Iterator()
	leafs, err := it.SubLeafs()
	require.NoError(t, err)
	exAcs := acs
	assert.Equal(t, 5, len(leafs))
	assertAcs := func(ac model.Account) {
		for i, exAc := range exAcs {
			if bytes.Equal(MustHash(exAc), MustHash(ac)) {
				exAcs = append(exAcs[:i], exAcs[i+1:]...)
				return
			}
		}
		assert.Failf(t, "assert accounts.", "%x is not found.", MustHash(ac))
	}
	for i, leaf := range leafs {
		ac := RandomAccount()
		require.NoError(t, leaf.Data(ac))
		assertAcs(ac)
		assert.Equal(t, 5-i-1, len(exAcs))
	}
	assert.Equal(t, 0, len(exAcs))

	// Second Upsert part 2 tree1
	err = tree1.Set(lastHash)
	require.NoError(t, err)
	for i, key := range keys {
		ac := RandomAccount()
		it, err := tree1.Find(key)
		require.NoError(t, err)
		err = it.Data(ac)
		assert.NoError(t, err)
		assert.Equal(t, MustHash(acs2[i]), MustHash(ac))
	}

	// Check SubTree2
	it = tree1.Iterator()
	leafs, err = it.SubLeafs()
	require.NoError(t, err)
	exAcs = acs2
	assert.Equal(t, 5, len(leafs))
	for i, leaf := range leafs {
		ac := RandomAccount()
		require.NoError(t, leaf.Data(ac))
		assertAcs(ac)
		assert.Equal(t, 5-i-1, len(exAcs))
	}
	assert.Equal(t, 0, len(exAcs))
}

func TestMerklePatriciaTree(t *testing.T) {
	cryptor := RandomCryptor()
	tree1, err := NewMerklePatriciaTree(RandomDBA(), cryptor, model.Hash(nil), MOCK_ROOT_KEY)
	require.NoError(t, err)
	tree2, err := NewMerklePatriciaTree(RandomDBA(), cryptor, model.Hash(nil), MOCK_ROOT_KEY)
	require.NoError(t, err)
	testMerklePatriciaTree(t, tree1, tree2)
}
