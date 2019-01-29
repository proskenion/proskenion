package repository_test

import (
	"github.com/pkg/errors"
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
	b3 := RandomBlock()
	b4 := RandomBlock()

	test_Blockchain_Upserts(t, bc, b1)
	require.NoError(t, tx.Commit())

	test_Blockchain_Upserts(t, bc, b2)
	require.NoError(t, tx.Commit())

	tx, err = dba.Begin()
	require.NoError(t, err)
	bcFirstHistory, err := NewBlockchainFromTopBlock(tx, RandomFactory(), RandomCryptor(), MustHash(b1))
	require.NoError(t, err)

	test_Blockchain_Upserts(t, bcFirstHistory, b3)
	require.NoError(t, err)
	require.NoError(t, tx.Commit())

	tx, err = dba.Begin()
	require.NoError(t, err)
	bcSecondHistory, err := NewBlockchainFromTopBlock(tx, RandomFactory(), RandomCryptor(), MustHash(b2))
	require.NoError(t, err)
	bcThirdHistory, err := NewBlockchainFromTopBlock(tx, RandomFactory(), RandomCryptor(), MustHash(b3))
	require.NoError(t, err)

	test_Blockchain_Upserts(t, bcSecondHistory, b4)

	// has b1 ok
	actb1, err := bcSecondHistory.Get(MustHash(b1))
	require.NoError(t, err)
	assert.Equal(t, MustHash(b1), MustHash(actb1))

	// has b1 ok
	actb1, err = bcThirdHistory.Get(MustHash(b1))
	require.NoError(t, err)
	assert.Equal(t, MustHash(b1), MustHash(actb1))

	// has b2 ok
	actb2, err := bcSecondHistory.Get(MustHash(b2))
	require.NoError(t, err)
	assert.Equal(t, MustHash(b2), MustHash(actb2))
	// has b2 ng
	actb2, err = bcThirdHistory.Get(MustHash(b2))
	assert.EqualError(t, errors.Cause(err), core.ErrBlockchainNotFound.Error())

	// has b3 ng
	actb3, err := bcSecondHistory.Get(MustHash(b3))
	assert.EqualError(t, errors.Cause(err), core.ErrBlockchainNotFound.Error())
	// has b3 ok
	actb3, err = bcThirdHistory.Get(MustHash(b3))
	require.NoError(t, err)
	assert.Equal(t, MustHash(b3), MustHash(actb3))
}

func TestBlockchain(t *testing.T) {
	dba := RandomDBA()
	test_Blockchain(t, dba)
}
