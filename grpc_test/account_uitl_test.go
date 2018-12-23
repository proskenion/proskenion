package grpc_test

import (
	"github.com/proskenion/proskenion/client"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type AccountManager struct {
	client     core.APIGateClient
	authorizer *AccountWithPri
	fc         model.ModelFactory
}

func andNewAccountManager(t *testing.T, server model.Peer) *AccountManager {
	fc := NewTestFactory()
	c, err := client.NewAPIGateClient(server, fc)
	require.NoError(t, err)
	return &AccountManager{
		c,
		NewAccountWithPri("authorizer@com"),
		fc,
	}
}

func (am *AccountManager) SetAuthorizer(t *testing.T) {
	tx := am.fc.NewTxBuilder().
		AddPublicKey(am.authorizer.AccountId, am.authorizer.AccountId, am.authorizer.Pubkey).
		Build()
	require.NoError(t, tx.Sign(am.authorizer.Pubkey, am.authorizer.Prikey))
	require.NoError(t, am.client.Write(tx))
}

func (am *AccountManager) CreateAccount(t *testing.T, ac *AccountWithPri) {
	tx := am.fc.NewTxBuilder().
		CreateAccount(am.authorizer.AccountId, ac.AccountId).
		AddPublicKey(am.authorizer.AccountId, ac.AccountId, ac.Pubkey).
		Build()
	require.NoError(t, tx.Sign(am.authorizer.Pubkey, am.authorizer.Prikey))
	require.NoError(t, am.client.Write(tx))
}

func (am *AccountManager) QueryAccountPassed(t *testing.T, ac *AccountWithPri) {
	query := am.fc.NewQueryBuilder().
		AuthorizerId(ac.AccountId).
		TargetId(ac.AccountId).
		CreatedTime(RandomNow()).
		RequestCode(model.AccountObjectCode).
		Build()
	require.NoError(t, query.Sign(ac.Pubkey, ac.Prikey))

	res, err := am.client.Read(query)
	require.NoError(t, err)

	assert.NoError(t, res.Verify())
	retAc := res.GetPayload().GetAccount()
	assert.Equal(t, retAc.GetAccountId(), ac.AccountId)
	assert.Equal(t, len(retAc.GetPublicKeys()), 1)
	assert.Contains(t, retAc.GetPublicKeys(), ac.Pubkey)
}
