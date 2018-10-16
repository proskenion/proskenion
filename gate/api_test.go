package gate_test

import (
	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/commit"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/gate"
	"github.com/proskenion/proskenion/query"
	"github.com/proskenion/proskenion/repository"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
)

func createAccount(t *testing.T, authorizer *AccountWithPri, target string) model.Transaction {
	tx := NewTestFactory().NewTxBuilder().
		CreateAccount(authorizer.AccountId, target).
		Build()
	require.NoError(t, tx.Sign(authorizer.Pubkey, authorizer.Prikey))
	return tx
}

func getAccountQuery(t *testing.T, authorizer *AccountWithPri, target string) model.Query {
	q := NewTestFactory().NewQueryBuilder().
		AuthorizerId(authorizer.AccountId).
		TargetId(target).
		RequestCode(model.AccountObjectCode).
		Build()
	require.NoError(t, q.Sign(authorizer.Pubkey, authorizer.Prikey))
	return q
}

func TestAPIGate_WriteAndRead(t *testing.T) {
	fc := NewTestFactory()
	rp := repository.NewRepository(RandomDBA(), RandomCryptor(), fc)
	queue := repository.NewProposalTxQueueOnMemory(NewTestConfig())
	logger := log15.New(context.TODO())
	qp := query.NewQueryProcessor(rp, fc)
	api := NewAPIGate(queue, logger, qp)
	cm := commit.NewCommitSystem(fc, RandomCryptor(), queue, RandomCommitProperty(), rp)

	// genesis Commit
	acs := []*AccountWithPri{
		NewAccountWithPri("authoirzer@com"),
		NewAccountWithPri("target1@com"),
		NewAccountWithPri("target2@com"),
		NewAccountWithPri("target3@com"),
	}
	GenesisCommitFromAccounts(t, rp, acs)

	txs := []model.Transaction{
		createAccount(t, acs[0], "target3@com"),
		createAccount(t, acs[0], "target4@com"),
		createAccount(t, acs[0], "target5@com"),
		createAccount(t, &AccountWithPri{acs[0].AccountId, acs[1].Pubkey, acs[1].Prikey}, "target6@com"),
	}
	for _, tx := range txs {
		require.NoError(t, api.Write(tx))
	}
	// Commit
	_, txList, err := cm.CreateBlock()
	require.NoError(t, err)

	assert.Equal(t, 2, txList.Size())
	assert.Equal(t, MustHash(txs[1]), MustHash(txList.List()[0]))
	assert.Equal(t, MustHash(txs[2]), MustHash(txList.List()[1]))

	for _, q := range []struct {
		query   model.Query
		pubkeys []model.PublicKey
		err     error
	}{
		{
			getAccountQuery(t, acs[0], "target1@com"),
			[]model.PublicKey{acs[1].Pubkey},
			nil,
		},
		{
			getAccountQuery(t, acs[0], "target2@com"),
			[]model.PublicKey{acs[2].Pubkey},
			nil,
		},
		{
			getAccountQuery(t, acs[0], "target3@com"),
			[]model.PublicKey{acs[3].Pubkey},
			nil,
		},
		{
			getAccountQuery(t, acs[0], "target4@com"),
			[]model.PublicKey{},
			nil,
		},
		{
			getAccountQuery(t, acs[0], "target5@com"),
			[]model.PublicKey{},
			nil,
		},
		{
			getAccountQuery(t, acs[0], "target6@com"),
			[]model.PublicKey{},
			core.ErrQueryProcessorNotFound,
		},
	} {
		res, err := api.Read(q.query)
		if q.err != nil {
			assert.EqualError(t, errors.Cause(err), q.err.Error())
		} else {
			require.NoError(t, err)
			assert.Equal(t, q.query.GetPayload().GetTargetId(), res.GetPayload().GetAccount().GetAccountId())
			assert.Equal(t, q.pubkeys, res.GetPayload().GetAccount().GetPublicKeys())
			assert.Equal(t, int64(0), res.GetPayload().GetAccount().GetAmount())
		}
	}
}
