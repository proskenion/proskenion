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
		CreateAccount(RandomStr(), RandomStr(), []model.PublicKey{}, 1).
		Build()
	return tx
}

func RandomValidTx() model.Transaction {
	validPub, validPriv := RandomKeyPairs()
	tx := RandomTx()
	tx.Sign(validPub, validPriv)
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

func RandomValidTxs() []model.Transaction {
	txs := make([]model.Transaction, 30)
	for id, _ := range txs {
		txs[id] = RandomValidTx()
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

func RandomTxs() []model.Transaction {
	return RandomValidTxs()
}

func RandomAccount() model.Account {
	return NewTestFactory().NewAccount(RandomStr(), RandomStr(), []model.PublicKey{RandomByte()}, rand.Int31(), rand.Int63(), RandomStr())

}

func RandomPeer() model.Peer {
	pub, _ := RandomKeyPairs()
	return NewTestFactory().NewPeer(RandomStr(), RandomStr(), pub)
}

func RandomBlock() model.Block {
	return NewTestFactory().NewBlockBuilder().
		Height(rand.Int63()).
		Round(0).
		WSVHash(RandomByte()).
		TxHistoryHash(RandomByte()).
		PreBlockHash(RandomByte()).
		CreatedTime(rand.Int63()).
		TxsHash(RandomByte()).
		Build()
}

func TxSign(t *testing.T, tx model.Transaction, pub []model.PublicKey, pri []model.PrivateKey) model.Transaction {
	require.Equal(t, len(pub), len(pri))
	for i, _ := range pub {
		require.NoError(t, tx.Sign(pub[i], pri[i]))
	}
	return tx
}

func GetAccountQuery(t *testing.T, authorizer *AccountWithPri, target string) model.Query {
	q := NewTestFactory().NewQueryBuilder().
		AuthorizerId(authorizer.AccountId).
		FromId(target).
		RequestCode(model.AccountObjectCode).
		Build()
	require.NoError(t, q.Sign(authorizer.Pubkey, authorizer.Prikey))
	return q
}

func CreateAccountTx(t *testing.T, authorizer *AccountWithPri, target string) model.Transaction {
	tx := NewTestFactory().NewTxBuilder().
		CreateAccount(authorizer.AccountId, target, []model.PublicKey{}, 0).
		Build()
	require.NoError(t, tx.Sign(authorizer.Pubkey, authorizer.Prikey))
	return tx
}

/*

func GetHash(t *testing.T, hasher model.Hasher) []byte {
	hash := hasher.Ge.Hash()
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
