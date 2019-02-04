package repository_test

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/repository"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func sameRepositoryTop(t *testing.T, rp core.Repository, block model.Block) {
	topBlock, ok := rp.Top()
	require.True(t, ok)
	assert.Equal(t, block, topBlock)

	rtx, err := rp.Begin()
	require.NoError(t, err)
	bc, err := rtx.Blockchain(topBlock.Hash())
	require.NoError(t, err)

	topBlock2, err := bc.Get(topBlock.Hash())
	require.NoError(t, err)
	assert.Equal(t, block.Hash(), topBlock2.Hash())
}

func TestRepository_Commit(t *testing.T) {
	rp := NewRepository(RandomDBA(), RandomCryptor(), RandomFactory(), RandomConfig())
	_, ok := rp.Top()
	assert.False(t, ok)
	require.NoError(t, rp.GenesisCommit(RandomGenesisTxList(t)))

	top, ok := rp.Top()
	assert.True(t, ok)

	tx := RandomFactory().NewTxBuilder().
		CreateAccount("authorizer@com", RandomStr()+"@com", []model.PublicKey{}, 0).
		CreatedTime(RandomNow()).Build()

	queue := RandomQueue()
	require.NoError(t, queue.Push(tx))
	newBlock, newTxList, err := rp.CreateBlock(queue, 0, RandomNow())
	require.NoError(t, err)
	require.Equal(t, tx.Hash(), newTxList.List()[0].Hash())
	sameRepositoryTop(t, rp, newBlock)

	// second == same result
	rp2 := NewRepository(RandomDBA(), RandomCryptor(), RandomFactory(), RandomConfig())
	require.NoError(t, rp2.GenesisCommit(RandomGenesisTxList(t)))
	top2, ok := rp2.Top()
	assert.True(t, ok)
	assert.Equal(t, top.Hash(), top2.Hash())
	require.NoError(t, rp2.Commit(newBlock, newTxList))
	sameRepositoryTop(t, rp2, newBlock)
}

func TestRepository_GetDelegatedAccounts(t *testing.T) {
	rp := NewRepository(RandomDBA(), RandomCryptor(), RandomFactory(), RandomConfig())
	require.NoError(t, rp.GenesisCommit(RandomGenesisTxList(t)))

	acs, err := rp.GetDelegatedAccounts()
	require.NoError(t, err)

	assert.Equal(t, 2, len(acs))
	assert.Equal(t, "root@peer", acs[0].GetDelegatePeerId())
	assert.Equal(t, "root@peer", acs[1].GetDelegatePeerId())
}
