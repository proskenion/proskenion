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
	fc := RandomFactory()
	cryptor := RandomCryptor()
	queue := RandomQueue()
	conf := RandomConfig()
	rp := repository.NewRepository(RandomDBA(), cryptor, fc, RandomConfig())
	require.NoError(t,rp.GenesisCommit(RandomGenesisTxList(t)))

	cs := NewCommitSystem(fc, cryptor, queue,  rp, conf)
	block, txList, err := cs.CreateBlock(0)
	require.NoError(t, err)
	assert.NoError(t, cs.VerifyCommit(block, txList))
	err = cs.Commit(block, txList)
	assert.Error(t, err)

	queue2 := RandomQueue()
	rp2 := repository.NewRepository(RandomDBA(), cryptor, fc, RandomConfig())
	require.NoError(t,rp2.GenesisCommit(RandomGenesisTxList(t)))
	cs2 := NewCommitSystem(fc, cryptor, queue2,  rp2, conf)

	assert.NoError(t, cs2.VerifyCommit(block, txList))
	assert.NoError(t, cs2.Commit(block, txList))

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
