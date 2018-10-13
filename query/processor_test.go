package query_test

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/query"
	"github.com/proskenion/proskenion/repository"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func genesisCommit(t *testing.T, rp core.Repository) {
	txList := repository.NewTxList(RandomCryptor())
	require.NoError(t, txList.Push(
		NewTestFactory().NewTxBuilder().
			CreateAccount("root", "authorizer@com").
			CreateAccount("root", "target@com").
			CreatedTime(0).
			Build()))
	require.NoError(t, rp.GenesisCommit(txList))
}

// TODO 不十分
func TestQueryProcessor_Query(t *testing.T) {
	fc := NewTestFactory()
	rp := repository.NewRepository(RandomDBA(), RandomCryptor(), fc)

	// GenesisCommit
	genesisCommit(t, rp)

	qp := NewQueryProcessor(rp, fc)

	query := fc.NewQueryBuilder().
		AuthorizerId("authorizer@com").
		TargetId("target@com").
		CreatedTime(RandomNow()).
		RequestCode(model.AccountObjectCode).
		Build()

	res, err := qp.Query(query)
	require.NoError(t, err)
	ac := res.GetPayload().GetAccount()
	assert.Equal(t, "target@com", ac.GetAccountId())

	q2 := fc.NewQueryBuilder().
		AuthorizerId("authorizer@com").
		TargetId("target1@com").
		CreatedTime(RandomNow()).
		RequestCode(model.AccountObjectCode).
		Build()

	_, err = qp.Query(q2)
	assert.EqualError(t, errors.Cause(err), core.ErrQueryProcessorNotFound.Error())
}
