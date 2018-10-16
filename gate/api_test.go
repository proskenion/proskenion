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

func createAccount(t *testing.T, authorizer string, target string, pub model.PublicKey, pri model.PrivateKey) model.Transaction {
	tx := NewTestFactory().NewTxBuilder().
		CreateAccount(authorizer, target).
		Build()
	require.NoError(t, tx.Sign(pub, pri))
	return tx
}

func getAccountQuery(authorizer string, target string) model.Query {
	return NewTestFactory().NewQueryBuilder().
		AuthorizerId(authorizer).
		TargetId(target).
		RequestCode(model.AccountObjectCode).
		Build()
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
		createAccount(t, acs[0].AccountId, "target3@com", acs[0].Pubkey, acs[0].prikey),
		createAccount(t, acs[0].AccountId, "target4@com", acs[0].Pubkey, acs[0].prikey),
		createAccount(t, acs[0].AccountId, "target5@com", acs[0].Pubkey, acs[0].prikey),
		createAccount(t, acs[0].AccountId, "target6@com", acs[1].Pubkey, acs[1].prikey),
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
		quqery  model.Query
		pubkeys []model.PublicKey
		err     error
	}{
		{
			getAccountQuery("authorizer@com", "target1@com"),
			[]model.PublicKey{acs[1].Pubkey},
			nil,
		},
		{
			getAccountQuery("authorizer@com", "target2@com"),
			[]model.PublicKey{acs[2].Pubkey},
			nil,
		},
		{
			getAccountQuery("authorizer@com", "target3@com"),
			[]model.PublicKey{acs[3].Pubkey},
			nil,
		},
		{
			getAccountQuery("authorizer@com", "target5@com"),
			[]model.PublicKey{},
			nil,
		},

		{
			getAccountQuery("authorizer@com", "target10@com"),
			[]model.PublicKey{},
			core.ErrQueryProcessorNotFound,
		},
	} {
		res, err := api.Read(q.quqery)
		if q.err != nil {
			assert.EqualError(t, errors.Cause(err), q.err.Error())
		} else {
			require.NoError(t, err)
			assert.Equal(t, q.quqery.GetPayload().GetTargetId(), res.GetPayload().GetAccount().GetAccountId())
			assert.Equal(t, q.pubkeys, res.GetPayload().GetAccount().GetPublicKeys())
			assert.Equal(t, 0, res.GetPayload().GetAccount().GetAmount())
		}
	}
}
