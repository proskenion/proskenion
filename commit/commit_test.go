package commit_test

import (
	"fmt"
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
	rp := repository.NewRepository(dba, cryptor, fc)

	cs := NewCommitSystem(dba, fc, cryptor, queue, cconf, rp)
	block, txList, err := cs.CreateBlock()
	require.NoError(t, err)
	assert.NoError(t, cs.VerifyCommit(block, txList))
	err = cs.Commit(block, txList)
	assert.Error(t, err)

	dba2 := RandomDBA()
	queue2 := RandomQueue()
	rp2 := repository.NewRepository(dba2, cryptor, fc)
	cs2 := NewCommitSystem(dba2, fc, cryptor, queue2, cconf, rp2)

	assert.NoError(t, cs2.VerifyCommit(block, txList))
	assert.NoError(t, cs2.Commit(block, txList))

	fmt.Println("blockHash: ", MustHash(block))
	rtx, err := rp.Begin()
	require.NoError(t, err)

	bc, err := rtx.Blockchain(MustHash(block))
	require.NoError(t, err)

	rtx2, err := rp2.Begin()
	require.NoError(t, err)
	bc2, err := rtx2.Blockchain(MustHash(block))

	b1, err := bc.Get(MustHash(block))
	require.NoError(t, err)
	b2, err := bc2.Get(MustHash(block))
	require.NoError(t, err)

	assert.Equal(t, MustHash(b1), MustHash(b2))
}
