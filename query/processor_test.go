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
	txList := repository.NewTxList(RandomCryptor(), RandomFactory())
	require.NoError(t, txList.Push(
		RandomFactory().NewTxBuilder().
			CreateAccount("root@/root", authorizer.AccountId, []model.PublicKey{authorizer.Pubkey}, 1).
			CreateAccount("root@/root", "target0@com", []model.PublicKey{}, 0).
			CreateAccount("root@/root", "target1@com", []model.PublicKey{}, 0).
			CreateAccount("root@/root", "target2@com", []model.PublicKey{}, 0).
			CreateAccount("root@/root", "target3@com", []model.PublicKey{}, 0).
			CreateAccount("root@/root", "target4@com", []model.PublicKey{}, 0).
			CreateAccount("root@/root", "targeta@pr", []model.PublicKey{}, 0).
			CreateAccount("root@/root", "targetb@pr", []model.PublicKey{}, 0).
			CreateAccount("root@/root", "targetc@pr", []model.PublicKey{}, 0).
			CreatedTime(0).
			Build()))
	require.NoError(t, rp.GenesisCommit(txList))
}

// TODO 不十分
func TestQueryProcessor_Query(t *testing.T) {
	fc := RandomFactory()
	rp := repository.NewRepository(RandomDBA(), RandomCryptor(), fc, RandomConfig())

	// GenesisCommit
	authorizer := NewAccountWithPri("authorizer@com/account")
	genesisCommit(t, rp, authorizer)

	qp := NewQueryProcessor( fc, RandomConfig())

	wsv, err := rp.TopWSV()
	require.NoError(t, err)
	defer wsv.Commit()
	query := GetAccountQuery(t, authorizer, "target0@com/account")
	res, err := qp.Query(wsv, query)
	require.NoError(t, err)
	ac := res.GetObject().GetAccount()
	assert.Equal(t, "target0@com", ac.GetAccountId())

	q2 := GetAccountQuery(t, authorizer, "targetx@com/account")
	_, err = qp.Query(wsv,q2)
	assert.EqualError(t, errors.Cause(err), core.ErrQueryProcessorNotFound.Error())

	q3 := GetAccountListQuery(t, authorizer, "com/account", "id", model.ASC, 100)
	res, err = qp.Query(wsv,q3)
	expctedIds := []string{
		"authorizer@com",
		"target0@com",
		"target1@com",
		"target2@com",
		"target3@com",
		"target4@com",
	}
	for i, l := range res.GetObject().GetList() {
		assert.Equal(t, expctedIds[i], l.GetAccount().GetAccountId())
	}

	q4 := GetAccountListQuery(t, authorizer, "com/account", "id", model.DESC, 100)
	res, err = qp.Query(wsv,q4)
	for i, l := range res.GetObject().GetList() {
		assert.Equal(t, expctedIds[len(expctedIds)-i-1], l.GetAccount().GetAccountId())
	}

	q5 := GetAccountListQuery(t, authorizer, "pr/account", "id", model.ASC, 100)
	res, err = qp.Query(wsv,q5)
	expctedIds2 := []string{
		"targeta@pr",
		"targetb@pr",
		"targetc@pr",
	}
	for i, l := range res.GetObject().GetList() {
		assert.Equal(t, expctedIds2[i], l.GetAccount().GetAccountId())
	}
}
