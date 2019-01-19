package test_utils

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/repository"
	"github.com/stretchr/testify/require"
	"testing"
)

func RandomCommandExecutor() core.CommandExecutor {
	_, ex, _, _, _, _, _ := NewTestFactories()
	return ex
}

func RandomCommandValidator() core.CommandValidator {
	_, _, vl, _, _, _, _ := NewTestFactories()
	return vl
}

type AccountWithPri struct {
	AccountId string
	Pubkey    model.PublicKey
	Prikey    model.PrivateKey
}

func GenesisCommitFromAccounts(t *testing.T, rp core.Repository, acs []*AccountWithPri) {
	txList, err := repository.NewTxListFromConf(RandomCryptor(), RandomProsl(), RandomConfig())
	require.NoError(t, err)

	builder := RandomFactory().NewTxBuilder()
	for _, ac := range acs {
		builder = builder.CreateAccount("root@com", ac.AccountId, []model.PublicKey{ac.Pubkey}, 1)
	}
	tx := builder.Build()
	require.NoError(t, txList.Push(tx))
	require.NoError(t, rp.GenesisCommit(txList))
}

func NewAccountWithPri(name string) *AccountWithPri {
	pub, pri := RandomKeyPairs()
	return &AccountWithPri{
		name,
		pub,
		pri,
	}
}

func RandomLandStorage(address string, id string, value int64, list []model.Object) model.Storage {
	return RandomFactory().NewStorageBuilder().
		Str("address", address).
		Address("owner", id).
		Int64("value", value).
		List("list", list).
		Build()
}
