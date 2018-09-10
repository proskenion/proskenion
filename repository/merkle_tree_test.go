package repository_test

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/crypto"
	. "github.com/proskenion/proskenion/repository"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func testMerkleTree_PushAndTop(t *testing.T, tree core.MerkleTree) {
	hashers := []core.Hasher{
		RandomMarshalerFromStr("1a"),
		RandomMarshalerFromStr("2a"),
		RandomMarshalerFromStr("3a"),
		RandomMarshalerFromStr("4a"),
	}

	expTop := []model.Hash{
		DecodeMustString(t, "58f7b0780592032e4d8602a3e8690fb2c701b2e1dd546e703445aabd6469734d"),
		DecodeMustString(t, "46f23c8ea5d95f0b22ee8e18f00c0aeb04d9ad79e83684eca4940217cdb4afc1"),
		DecodeMustString(t, "4d56cecf438ff85bd2759a569a19e0c3a8b4428c89628a4214c2c680398d4c61"),
		DecodeMustString(t, "a31011ebbb39ad3ee003c575d4fbd2e475345d3417c6c814e3bc7c583f1b8049"),
	}

	assert.Equal(t, model.Hash(nil), tree.Top())
	for i, hasher := range hashers {
		err := tree.Push(hasher)
		require.NoError(t, err)
		assert.Equal(t, expTop[i], tree.Top())
	}
}

func TestAccumulateHash_PushAndTop(t *testing.T) {
	cryptor := crypto.NewEd25519Sha256Cryptor()
	accumulateHash := NewAccumulateHash(cryptor)
	testMerkleTree_PushAndTop(t, accumulateHash)
}
