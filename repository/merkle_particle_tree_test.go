package repository_test

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func testUpsertFirst(t *testing.T, tree core.MerkleParticleTree, node *RandomMockMarshaler) {
	v := RandomMarshaler()
	err := tree.Find(node, v)
	assert.EqualError(t, errors.Cause(err), core.ErrMerkleParticleTreeNotFoundKey.Error())
	err = tree.Upsert(node, node)
	require.NoError(t, err)
	err = tree.Find(node, v)
	assert.Equal(t, node, v)
}

func testUpsertSecond(t *testing.T, tree core.MerkleParticleTree, key *RandomMockMarshaler, value *RandomMockMarshaler) {
	v := RandomMarshaler()
	err := tree.Find(key, v)
	require.NoError(t, err)
	assert.Equal(t, key, v)

	err = tree.Upsert(key, value)
	assert.NoError(t, err)
	err = tree.Find(key, v)
	assert.Equal(t, value, v)
}

func testMerkleParticleTree(t *testing.T, tree1 core.MerkleParticleTree, tree2 core.MerkleParticleTree) {
	nodes := []*RandomMockMarshaler{
		RandomMarshaler(),
		RandomMarshaler(),
		RandomMarshaler(),
		RandomMarshaler(),
		RandomMarshaler(),
	}
	nodes2 := []*RandomMockMarshaler{
		RandomMarshaler(),
		RandomMarshaler(),
		RandomMarshaler(),
		RandomMarshaler(),
		RandomMarshaler(),
	}
	nodes3 := []*RandomMockMarshaler{
		RandomMarshaler(),
		RandomMarshaler(),
		RandomMarshaler(),
		RandomMarshaler(),
		RandomMarshaler(),
	}
	// First Upsert tree1 and tree2
	for _, node := range nodes {
		testUpsertFirst(t, tree1, node)
		testUpsertFirst(t, tree2, node)
		assert.Equal(t, MustHash(tree1), MustHash(tree2))
	}
	firstHash := MustHash(tree1)

	// Second Upsert tree1 and tree2
	for i, node := range nodes {
		testUpsertSecond(t, tree1, node, nodes2[i])
		testUpsertSecond(t, tree2, node, nodes2[i])
		assert.Equal(t, MustHash(tree1), MustHash(tree2))
	}
	lastHash := MustHash(tree1)

	// Compare first hash and last hash
	assert.NotEqual(t, firstHash, lastHash)

	// First Upsert 終了時に巻き戻し
	actFirstNode := tree1.Root().Iterator().Prev()
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
