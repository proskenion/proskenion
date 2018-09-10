package repository

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/crypto"
	"github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testMerkleTree_PushAndTop(t *testing.T, tree core.MerkleTree) {
	txs := test_utils.RandomTxs(t)

	for _, tx := range txs {
		preHash := tree.Top()
		tree.Push(tx)
		assert.NotEqual(t, preHash, tree.Top())
	}
}

func TestAccumulateHash_PushAndTop(t *testing.T) {
	cryptor := crypto.NewEd25519Sha256Cryptor()
	accumulateHash := NewAccumulateHash(cryptor)
	testMerkleTree_PushAndTop(t, accumulateHash)
}
