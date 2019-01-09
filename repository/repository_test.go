package repository_test

import (
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/repository"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRepository_Commit(t *testing.T) {
	dba := RandomDBA()
	cryptor := RandomCryptor()
	fc := NewTestFactory()

	rp := NewRepository(dba, cryptor, fc)

	txList := RandomTxList()
	txList.Push(
		fc.NewTxBuilder().
			CreateAccount("authorizer@com", "authorizer@com", []model.PublicKey{}, 0).
			Build())
	require.NoError(t, rp.GenesisCommit(txList))

	top, ok := rp.Top()
	require.True(t, ok)

	block, txList := RandomCommitableBlock(t, top, rp)
	assert.NoError(t, rp.Commit(block, txList))

	topBlock, ok := rp.Top()
	require.True(t, ok)
	assert.Equal(t, block, topBlock)

	rtx, err := rp.Begin()
	require.NoError(t, err)
	bc, err := rtx.Blockchain(MustHash(topBlock))
	require.NoError(t, err)

	topBlock2, err := bc.Get(MustHash(topBlock))
	require.NoError(t, err)
	assert.Equal(t, MustHash(block), MustHash(topBlock2))
}
