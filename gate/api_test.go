package gate_test

import (
	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/commit"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/gate"
	"github.com/proskenion/proskenion/p2p"
	"github.com/proskenion/proskenion/query"
	"github.com/proskenion/proskenion/repository"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
)

func TestAPI_WriteAndRead(t *testing.T) {
	fc := RandomFactory()
	conf := RandomConfig()
	rp := repository.NewRepository(RandomDBA(), RandomCryptor(), fc, RandomConfig())
	queue := repository.NewProposalTxQueueOnMemory(RandomConfig())
	logger := log15.New(context.TODO())
	qp := query.NewQueryProcessor(fc, RandomConfig())
	qv := query.NewQueryValidator(fc, RandomConfig())
	api := NewAPI(rp, queue, qp, qv, &p2p.MockGossip{}, logger)
	cm := commit.NewCommitSystem(fc, RandomCryptor(), queue, rp, conf)

	// genesis Commit
	acs := []*AccountWithPri{
		NewAccountWithPri("authoirzer@com"),
		NewAccountWithPri("target1@com"),
		NewAccountWithPri("target2@com"),
		NewAccountWithPri("target3@com"),
	}
	GenesisCommitFromAccounts(t, rp, acs)

	txs := []model.Transaction{
		CreateAccountTx(t, acs[0], "target3@com"),
		CreateAccountTx(t, acs[0], "target4@com"),
		CreateAccountTx(t, acs[0], "target5@com"),
		CreateAccountTx(t, &AccountWithPri{acs[0].AccountId, acs[1].Pubkey, acs[1].Prikey}, "target6@com"),
	}
	for _, tx := range txs {
		require.NoError(t, api.Write(tx))
	}
	// Commit
	_, txList, err := cm.CreateBlock(0)
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
			GetAccountQuery(t, acs[0], "target1@com"),
			[]model.PublicKey{acs[1].Pubkey},
			nil,
		},
		{
			GetAccountQuery(t, acs[0], "target2@com"),
			[]model.PublicKey{acs[2].Pubkey},
			nil,
		},
		{
			GetAccountQuery(t, acs[0], "target3@com"),
			[]model.PublicKey{acs[3].Pubkey},
			nil,
		},
		{
			GetAccountQuery(t, acs[0], "target4@com"),
			[]model.PublicKey{},
			nil,
		},
		{
			GetAccountQuery(t, acs[0], "target5@com"),
			[]model.PublicKey{},
			nil,
		},
		{
			GetAccountQuery(t, acs[0], "target6@com"),
			[]model.PublicKey{},
			core.ErrAPIQueryNotFound,
		},
		{
			GetAccountQuery(t, &AccountWithPri{acs[0].AccountId, acs[1].Pubkey, acs[1].Prikey}, "target1@com"),
			[]model.PublicKey{},
			core.ErrAPIQueryValidateError,
		},
		{
			GetAccountQuery(t, &AccountWithPri{"auth@com", acs[1].Pubkey, acs[1].Prikey}, "target1@com"),
			[]model.PublicKey{},
			core.ErrAPIQueryValidateError,
		},
	} {
		res, err := api.Read(q.query)
		if q.err != nil {
			assert.EqualError(t, errors.Cause(err), q.err.Error())
		} else {
			require.NoError(t, err)
			assert.Equal(t, model.MustAddress(q.query.GetPayload().GetFromId()).Account(), res.GetObject().GetAccount().GetAccountName())
			assert.Equal(t, q.pubkeys, res.GetObject().GetAccount().GetPublicKeys())
			assert.Equal(t, int64(0), res.GetObject().GetAccount().GetBalance())
		}
	}
}
