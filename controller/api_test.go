package controller_test

import (
	"github.com/inconshreveable/log15"
	. "github.com/proskenion/proskenion/controller"
	"github.com/proskenion/proskenion/convertor"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/gate"
	"github.com/proskenion/proskenion/p2p"
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

func initializeAPI(t *testing.T) ([]*AccountWithPri, core.ProposalTxQueue, proskenion.APIServer) {
	fc := RandomFactory()
	conf := RandomConfig()
	rp := repository.NewRepository(RandomDBA(), RandomCryptor(), fc, conf)
	queue := repository.NewProposalTxQueueOnMemory(RandomConfig())
	logger := log15.New(context.TODO())
	qp := query.NewQueryProcessor( fc, RandomConfig())
	qv := query.NewQueryValidator( fc, conf)
	api := gate.NewAPI(rp, queue, qp, qv, &p2p.MockGossip{}, logger)

	server := NewAPIServer(fc, api, logger)

	// genesis Commit
	acs := []*AccountWithPri{
		NewAccountWithPri("authoirzer@com"),
		NewAccountWithPri("target1@com"),
		NewAccountWithPri("target2@com"),
		NewAccountWithPri("target3@com"),
	}
	GenesisCommitFromAccounts(t, rp, acs)
	return acs, queue, server
}

func statusCheck(t *testing.T, err error, code codes.Code) {
	require.Error(t, err)
	assert.Equalf(t, status.Code(err), code, "expected: %s, atctual: %s, error: %s", status.Code(err).String(), code.String(), err.Error())
}

func TestAPIServer_Write(t *testing.T) {
	acs, queue, server := initializeAPI(t)
	for _, c := range []struct {
		name string
		tx   model.Transaction
		code codes.Code
	}{
		{
			"case 1 ok",
			CreateAccountTx(t, acs[0], "target3@com"),
			codes.OK,
		},
		{
			"case 2 ok",
			CreateAccountTx(t, acs[0], "target4@com"),
			codes.OK,
		},
		{
			"case 3 ok",
			CreateAccountTx(t, acs[0], "target5@com"),
			codes.OK,
		},
		{
			"case 4 invalid",
			RandomInvalidTx(t),
			codes.InvalidArgument,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			_, err := server.Write(context.TODO(), c.tx.(*convertor.Transaction).Transaction)
			if c.code != codes.OK {
				statusCheck(t, err, c.code)
				return
			}
			require.NoError(t, err)

			_, err = server.Write(context.TODO(), c.tx.(*convertor.Transaction).Transaction)
			statusCheck(t, err, codes.AlreadyExists)

			actTx, ok := queue.Pop()
			require.True(t, ok)
			assert.Equal(t, MustHash(c.tx), actTx.Hash())
		})
	}
}

func TestAPIServer_Query(t *testing.T) {
	acs, _, server := initializeAPI(t)

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
			"case 6 not found",
			GetAccountQuery(t, acs[0], "target6@com"),
			[]model.PublicKey{},
			codes.NotFound,
		},
		{
			"case 7 invalid",
			GetAccountQuery(t, &AccountWithPri{acs[0].AccountId, acs[1].Pubkey, acs[1].Prikey}, "target1@com"),
			[]model.PublicKey{},
			codes.InvalidArgument,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			res, err := server.Read(context.TODO(), c.query.(*convertor.Query).Query)
			if c.code != codes.OK {
				statusCheck(t, err, c.code)
			} else {
				require.NoError(t, err)
				resq := RandomFactory().NewEmptyQueryResponse()
				resq.(*convertor.QueryResponse).QueryResponse = res

				assert.Equal(t, model.MustAddress(c.query.GetPayload().GetFromId()).Account(), resq.GetObject().GetAccount().GetAccountName())
				assert.Equal(t, c.pubkeys, resq.GetObject().GetAccount().GetPublicKeys())
				assert.Equal(t, int64(0), resq.GetObject().GetAccount().GetBalance())
			}
		})
	}
}
