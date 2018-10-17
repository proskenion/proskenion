package controller_test

import (
	"github.com/inconshreveable/log15"
	"github.com/proskenion/proskenion/commit"
	. "github.com/proskenion/proskenion/controller"
	"github.com/proskenion/proskenion/convertor"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/gate"
	"github.com/proskenion/proskenion/proto"
	"github.com/proskenion/proskenion/query"
	"github.com/proskenion/proskenion/repository"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"testing"
)

func initializeAPIGate(t *testing.T) ([]*AccountWithPri, proskenion.APIGateServer) {
	fc := NewTestFactory()
	rp := repository.NewRepository(RandomDBA(), RandomCryptor(), fc)
	queue := repository.NewProposalTxQueueOnMemory(NewTestConfig())
	logger := log15.New(context.TODO())
	qp := query.NewQueryProcessor(rp, fc)
	api := gate.NewAPIGate(queue, logger, qp)

	server := NewAPIGateServer(fc, api, logger)

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
		CreateAccountTx(t, acs[0], "target3@com"),
		CreateAccountTx(t, acs[0], "target4@com"),
		CreateAccountTx(t, acs[0], "target5@com"),
		CreateAccountTx(t, &AccountWithPri{acs[0].AccountId, acs[1].Pubkey, acs[1].Prikey}, "target6@com"),
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

	return acs, server
}

func statusCheck(t *testing.T, err error, code codes.Code) {
	assert.Equalf(t, status.Code(err), code, err.Error())
}

func TestAPIGateServer_Write(t *testing.T) {
	acs, server := initializeAPIGate(t)

	for _, c := range []struct {
		name    string
		query   model.Query
		pubkeys []model.PublicKey
		code    codes.Code
	}{
		{
			"case 1 ok",
			GetAccountQuery(t, acs[0], "target1@com"),
			[]model.PublicKey{acs[1].Pubkey},
			codes.OK,
		},
		{
			"case 2 ok",
			GetAccountQuery(t, acs[0], "target2@com"),
			[]model.PublicKey{acs[2].Pubkey},
			codes.OK,
		},
		{
			"case 3 ok",
			GetAccountQuery(t, acs[0], "target3@com"),
			[]model.PublicKey{acs[3].Pubkey},
			codes.OK,
		},
		{
			"case 4 ok",
			GetAccountQuery(t, acs[0], "target4@com"),
			[]model.PublicKey{},
			codes.OK,
		},
		{
			"case 5 ok",
			GetAccountQuery(t, acs[0], "target5@com"),
			[]model.PublicKey{},
			codes.OK,
		},
		{
			"case 6 not found",
			GetAccountQuery(t, acs[0], "target6@com"),
			[]model.PublicKey{},
			codes.NotFound,
		},
		{//TODO this is invalidArguments
			"case 7 invalid",
			GetAccountQuery(t, &AccountWithPri{acs[0].AccountId, acs[1].Pubkey, acs[1].Prikey}, "target1@com"),
			[]model.PublicKey{},
			codes.Internal,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			res, err := server.Read(context.TODO(), c.query.(*convertor.Query).Query)
			if c.code != codes.OK {
				statusCheck(t, err, c.code)
			} else {
				require.NoError(t, err)
				resq := NewTestFactory().NewEmptyQueryResponse()
				resq.(*convertor.QueryResponse).QueryResponse = res

				assert.Equal(t, c.query.GetPayload().GetTargetId(), resq.GetPayload().GetAccount().GetAccountId())
				assert.Equal(t, c.pubkeys, resq.GetPayload().GetAccount().GetPublicKeys())
				assert.Equal(t, int64(0), resq.GetPayload().GetAccount().GetAmount())
			}
		})
	}
}
