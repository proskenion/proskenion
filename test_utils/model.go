package test_utils

import (
	"github.com/proskenion/proskenion/command"
	"github.com/proskenion/proskenion/convertor"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/crypto"
	"github.com/proskenion/proskenion/query"
	"github.com/stretchr/testify/require"
	"math/rand"
	"strconv"
	"testing"
)

func NewTestFactory() model.ModelFactory {
	return convertor.NewModelFactory(
		crypto.NewEd25519Sha256Cryptor(),
		command.NewCommandExecutor(),
		command.NewCommandValidator(),
		query.NewQueryValidator(),
	)
}

func RandomStr() string {
	return strconv.FormatUint(rand.Uint64(), 36)
}

func RandomByte() []byte {
	b, _ := RandomKeyPairs()
	return b
}

func RandomInvalidSig() model.Signature {
	pub, sig := RandomKeyPairs()
	return NewTestFactory().NewSignature(pub, sig)
}

func RandomTx() model.Transaction {
	tx := NewTestFactory().NewTxBuilder().
		CreatedTime(rand.Int63()).
		Transfer(RandomStr(), RandomStr(), rand.Int63()).
		Build()
	return tx
}

func RandomValidTx(t *testing.T) model.Transaction {
	validPub, validPriv := RandomKeyPairs()
	tx := RandomTx()
	err := tx.Sign(validPub, validPriv)
	require.NoError(t, err)
	return tx
}

func RandomInvalidTx(t *testing.T) model.Transaction {
	pub, _ := RandomKeyPairs()
	_, pri := RandomKeyPairs()
	tx := RandomTx()
	err := tx.Sign(pub, pri)
	require.NoError(t, err)
	return tx
}

func RandomValidTxs(t *testing.T) []model.Transaction {
	txs := make([]model.Transaction, 30)
	for id, _ := range txs {
		txs[id] = RandomValidTx(t)
	}
	return txs
}

func RandomInvalidTxs(t *testing.T) []model.Transaction {
	txs := make([]model.Transaction, 30)
	for id, _ := range txs {
		txs[id] = RandomInvalidTx(t)
	}
	return txs
}

func RandomTxs(t *testing.T) []model.Transaction {
	return RandomValidTxs(t)
}

func RandomAccount() model.Account {
	return NewTestFactory().NewAccount(RandomStr(), RandomStr(), []model.PublicKey{RandomByte()}, rand.Int63())
}

func RandomPeer() model.Peer {
	pub, _ := RandomKeyPairs()
	return NewTestFactory().NewPeer(RandomStr(), pub)
}

/*




func GetHash(t *testing.T, hasher model.Hasher) []byte {
	hash, err := hasher.GetHash()
	require.NoError(t, err)
	return hash
}

func RandomValidBlock(t *testing.T) model.Block {
	block, err := convertor.NewModelFactory().NewBlock(rand.Int63(), RandomByte(), rand.Int63(), RandomValidTxs(t))
	require.NoError(t, err)
	return block
}

func RandomInvalidBlock(t *testing.T) model.Block {
	block, err := convertor.NewModelFactory().NewBlock(rand.Int63(), RandomByte(), rand.Int63(), RandomInvalidTxs(t))
	require.NoError(t, err)
	return block
}

func RandomBlock(t *testing.T) model.Block {
	return RandomValidBlock(t)
}

func ValidSignedBlock(t *testing.T) model.Block {
	validPub, validPri := convertor.NewKeyPair()
	block := RandomValidBlock(t)

	err := block.Sign(validPub, validPri)
	require.NoError(t, err)
	require.NoError(t, block.Verify())
	return block
}

func InvalidSingedBlock(t *testing.T) model.Block {
	validPub, validPri := convertor.NewKeyPair()
	block := RandomInvalidBlock(t)

	err := block.Sign(validPub, validPri)
	require.NoError(t, err)
	require.NoError(t, block.Verify())
	return block
}

func InvalidErrSignedBlock(t *testing.T) model.Block {
	inValidPub := RandomByte()
	inValidPriv := RandomByte()
	block := RandomInvalidBlock(t)

	err := block.Sign(inValidPub, inValidPriv)
	require.Error(t, err)
	require.Error(t, block.Verify())
	return block
}

func ValidErrSignedBlock(t *testing.T) model.Block {
	inValidPub := RandomByte()
	inValidPriv := RandomByte()
	block := RandomInvalidBlock(t)

	err := block.Sign(inValidPub, inValidPriv)
	require.Error(t, err)
	require.Error(t, block.Verify())
	return block
}

func RandomPeer() model.Peer {
	return convertor.NewModelFactory().NewPeer(RandomStr(), RandomByte())
}
*/
