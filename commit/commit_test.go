package commit_test

import (
	. "github.com/proskenion/proskenion/commit"
	"github.com/proskenion/proskenion/repository"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCommitSystem_CreateBlock_Commit(t *testing.T) {
	dba := RandomDBA()
	fc := NewTestFactory()
	cryptor := RandomCryptor()
	queue := RandomQueue()
	cconf := RandomCommitProperty()

	cs := NewCommitSystem(dba, fc, cryptor, queue, cconf)
	block, txList, err := cs.CreateBlock()
	require.NoError(t, err)

	dba2 := RandomDBA()
	queue2 := RandomQueue()
	cs2 := NewCommitSystem(dba2, fc, cryptor, queue2, cconf)

	assert.NoError(t, cs2.VerifyCommit(block, txList))
	assert.NoError(t, cs2.Commit(block, txList))

	tx, err := dba.Begin()
	require.NoError(t, err)
	bc := repository.NewBlockchain(tx, fc)

	tx2, err := dba2.Begin()
	require.NoError(t, err)
	bc2 := repository.NewBlockchain(tx2, fc)

	blockHash, err := block.Hash()
	require.NoError(t, err)
	b1, ok := bc.Get(blockHash)
	require.True(t, ok)
	b2, ok := bc2.Get(blockHash)
	require.True(t, ok)

	assert.Equal(t, MustHash(b1), MustHash(b2))
}
