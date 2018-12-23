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

func genesisCommit(t *testing.T, rp core.Repository, authorizer *AccountWithPri) {
	txList := repository.NewTxList(RandomCryptor())
	require.NoError(t, txList.Push(
		NewTestFactory().NewTxBuilder().
			CreateAccount("root", authorizer.AccountId, []model.PublicKey{authorizer.Pubkey}, 1).
			CreateAccount("root", "target@com", []model.PublicKey{}, 0).
			CreatedTime(0).
			Build()))
	require.NoError(t, rp.GenesisCommit(txList))
}

// TODO 不十分
func TestQueryProcessor_Query(t *testing.T) {
	fc := NewTestFactory()
	rp := repository.NewRepository(RandomDBA(), RandomCryptor(), fc)

	// GenesisCommit
	authorizer := NewAccountWithPri("authorizer@com")
	genesisCommit(t, rp, authorizer)

	qp := NewQueryProcessor(rp, fc, NewTestConfig())

	query := GetAccountQuery(t, authorizer, "target@com")
	res, err := qp.Query(query)
	require.NoError(t, err)
	ac := res.GetPayload().GetAccount()
	assert.Equal(t, "target@com", ac.GetAccountId())

	q2 := GetAccountQuery(t, authorizer, "target1@com")
	_, err = qp.Query(q2)
	assert.EqualError(t, errors.Cause(err), core.ErrQueryProcessorNotFound.Error())

	tmpub, tmpri := RandomKeyPairs()
	q3 := GetAccountQuery(t, &AccountWithPri{authorizer.AccountId, tmpub, tmpri}, "target@com")
	_, err = qp.Query(q3)
	assert.EqualError(t, errors.Cause(err), core.ErrQueryProcessorNotSignedAuthorizer.Error())

	q4 := GetAccountQuery(t, &AccountWithPri{"authorizer1@com", tmpub, tmpri}, "target@com")
	_, err = qp.Query(q4)
	assert.EqualError(t, errors.Cause(err), core.ErrQueryProcessorNotExistAuthoirizer.Error())
}
