package test_utils

import (
	"github.com/proskenion/proskenion/command"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/repository"
	"github.com/stretchr/testify/require"
	"testing"
)

func RandomCommandExecutor() core.CommandExecutor {
	fc := NewTestFactory()
	ex := command.NewCommandExecutor()
	ex.SetFactory(fc)
	return ex
}

func RandomCommandValidator() core.CommandValidator {
	return command.NewCommandValidator()
}

type AccountWithPri struct {
	AccountId string
	Pubkey    model.PublicKey
	Prikey    model.PrivateKey
}

func GenesisCommitFromAccounts(t *testing.T, rp core.Repository, acs []*AccountWithPri) {
	txList := repository.NewTxList(RandomCryptor())

	builder := NewTestFactory().NewTxBuilder()
	for _, ac := range acs {
		builder = builder.CreateAccount("root", ac.AccountId).
			AddPublicKey("root", ac.AccountId, ac.Pubkey)
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
