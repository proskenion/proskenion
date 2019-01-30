package repository_test

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/convertor"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/repository"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func test_Blockchain_Upserts(t *testing.T, blockchain core.Blockchain, block model.Block) {
	_, err := blockchain.Get(block.Hash())
	require.EqualError(t, errors.Cause(err), core.ErrBlockchainNotFound.Error())
	err = blockchain.Append(block)
	require.NoError(t, err)

	retBlock, err := blockchain.Get(block.Hash())
	require.NoError(t, err)
	assert.Equal(t, block.Hash(), retBlock.Hash())
}

func test_Blockchain(t *testing.T, dba core.DBA) {
	tx, err := dba.Begin()
	require.NoError(t, err)

	bc, err := NewBlockchainFromTopBlock(tx, RandomFactory(), RandomCryptor(), model.Hash(nil))
	require.NoError(t, err)

	b1 := RandomBlock()
	b2 := RandomBlock()
	b2.(*convertor.Block).Block.Payload.PreBlockHash = b1.Hash()
	b3 := RandomBlock()
	b3.(*convertor.Block).Block.Payload.PreBlockHash = b1.Hash()
	b4 := RandomBlock()
	b4.(*convertor.Block).Block.Payload.PreBlockHash = b2.Hash()

	t.Run("case 1 : scucess b1 -> b2", func(t *testing.T) {
		test_Blockchain_Upserts(t, bc, b1)
		require.NoError(t, tx.Commit())

		test_Blockchain_Upserts(t, bc, b2)
		require.NoError(t, tx.Commit())
	})

	t.Run("case 2 : success b1 -> b3", func(t *testing.T) {
		tx, err = dba.Begin()
		require.NoError(t, err)
		bcFirstHistory, err := NewBlockchainFromTopBlock(tx, RandomFactory(), RandomCryptor(), b1.Hash())
		require.NoError(t, err)

		test_Blockchain_Upserts(t, bcFirstHistory, b3)
		require.NoError(t, err)
		require.NoError(t, tx.Commit())
	})

	t.Run("case 3 : sucess b2 -> b4, b3 -> b4", func(t *testing.T) {
		tx, err = dba.Begin()
		require.NoError(t, err)
		bcSecondHistory, err := NewBlockchainFromTopBlock(tx, RandomFactory(), RandomCryptor(), b2.Hash())
		require.NoError(t, err)
		bcThirdHistory, err := NewBlockchainFromTopBlock(tx, RandomFactory(), RandomCryptor(), b3.Hash())
		require.NoError(t, err)

		// b2 -> b4
		test_Blockchain_Upserts(t, bcSecondHistory, b4)

		// b2 has b1 ok
		actb1, err := bcSecondHistory.Get(b1.Hash())
		require.NoError(t, err)
		assert.Equal(t, b1.Hash(), MustHash(actb1))

		// b3 has b1 ok
		actb1, err = bcThirdHistory.Get(b1.Hash())
		require.NoError(t, err)
		assert.Equal(t, b1.Hash(), MustHash(actb1))

		// b2 has b2 ok
		actb2, err := bcSecondHistory.Get(b2.Hash())
		require.NoError(t, err)
		assert.Equal(t, b2.Hash(), MustHash(actb2))
		// b3 has b2 ng
		actb2, err = bcThirdHistory.Get(b2.Hash())
		assert.EqualError(t, errors.Cause(err), core.ErrBlockchainNotFound.Error())

		// b3 has b3 ng
		actb3, err := bcSecondHistory.Get(b3.Hash())
		assert.EqualError(t, errors.Cause(err), core.ErrBlockchainNotFound.Error())
		// b3 has b3 ok
		actb3, err = bcThirdHistory.Get(b3.Hash())
		require.NoError(t, err)
		assert.Equal(t, b3.Hash(), MustHash(actb3))

		// b2 has b4 ok
		actb4, err := bcSecondHistory.Get(b4.Hash())
		require.NoError(t, err)
		assert.Equal(t, b4.Hash(), MustHash(actb4))
		// b3 has b4 ok
		_, err = bcThirdHistory.Get(b4.Hash())
		assert.Errorf(t, errors.Cause(err), core.ErrBlockchainNotFound.Error())
		assert.Equal(t, b4.Hash(), MustHash(actb4))

		// b2 has b1 next-> b2
		nxtb1, err := bcSecondHistory.Next(b1.Hash())
		require.NoError(t, err)
		assert.Equal(t, b2.Hash(), nxtb1.Hash())
		// b2 has b2 next-> b4
		nxtb2, err := bcSecondHistory.Next(b2.Hash())
		require.NoError(t, err)
		assert.Equal(t, b4.Hash(), nxtb2.Hash())
		// b3 has b1 next-> b3
		nxtb1, err = bcThirdHistory.Next(b1.Hash())
		require.NoError(t, err)
		assert.Equal(t, b3.Hash(), nxtb1.Hash())
	})

	t.Run("case 4 : 100 next", func(t *testing.T) {
		rp := RandomRepository()
		require.NoError(t, rp.GenesisCommit(RandomGenesisTxList(t)))
		blocks := make([]model.Block, 0)
		blocks = append(blocks, MusTop(rp))
		limits := 100
		for i := 0; i < limits; i++ {
			block, _ := RandomCommitableBlockAndTxList(t, rp)
			blocks = append(blocks, block)
		}
		rtx, err := rp.Begin()
		require.NoError(t, err)

		bc, err := rtx.Blockchain(MusTop(rp).Hash())
		for i := 0; i < limits; i++ {
			nxt, err := bc.Next(blocks[i].Hash())
			require.NoError(t, err)
			assert.Equal(t, blocks[i+1].Hash(), nxt.Hash())
		}
	})
}

func TestBlockchain(t *testing.T) {
	dba := RandomDBA()
	test_Blockchain(t, dba)
}
