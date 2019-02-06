package grpc_test

import (
	"fmt"
	"github.com/proskenion/proskenion/client"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type AccountManager struct {
	client     core.APIClient
	authorizer *AccountWithPri
	fc         model.ModelFactory
}

func NewAccountManager(t *testing.T, server model.Peer) *AccountManager {
	fc := RandomFactory()
	c, err := client.NewAPIClient(server, fc)
	require.NoError(t, err)
	return &AccountManager{
		c,
		NewAccountWithPri("authorizer@pr"),
		fc,
	}
}

func (am *AccountManager) SetAuthorizer(t *testing.T) {
	tx := am.fc.NewTxBuilder().
		AddPublicKeys(am.authorizer.AccountId, am.authorizer.AccountId, []model.PublicKey{am.authorizer.Pubkey}).
		SetQuorum(am.authorizer.AccountId, am.authorizer.AccountId, 1).
		Build()
	require.NoError(t, tx.Sign(am.authorizer.Pubkey, am.authorizer.Prikey))
	require.NoError(t, am.client.Write(tx))
}

func (am *AccountManager) CreateAccount(t *testing.T, ac *AccountWithPri) {
	tx := am.fc.NewTxBuilder().
		CreateAccount(am.authorizer.AccountId, ac.AccountId, []model.PublicKey{ac.Pubkey}, 1).
		Build()
	require.NoError(t, tx.Sign(am.authorizer.Pubkey, am.authorizer.Prikey))
	require.NoError(t, am.client.Write(tx))
}

func (am *AccountManager) AddPeer(t *testing.T, peer model.Peer) {
	tx := am.fc.NewTxBuilder().
		CreateAccount(am.authorizer.AccountId, peer.GetPeerId(), []model.PublicKey{peer.GetPublicKey()}, 1).
		AddPeer(am.authorizer.AccountId, peer.GetPeerId(), peer.GetAddress(), peer.GetPublicKey()).
		Build()
	require.NoError(t, tx.Sign(am.authorizer.Pubkey, am.authorizer.Prikey))
	require.NoError(t, am.client.Write(tx))
}

func (am *AccountManager) QueryAccountPassed(t *testing.T, ac *AccountWithPri) {
	query := am.fc.NewQueryBuilder().
		AuthorizerId(ac.AccountId).
		FromId(model.MustAddress(ac.AccountId).AccountId()).
		CreatedTime(RandomNow()).
		RequestCode(model.AccountObjectCode).
		Build()
	assert.NoError(t, query.Sign(ac.Pubkey, ac.Prikey))

	res, err := am.client.Read(query)
	assert.NoError(t, err)

	assert.NoError(t, res.Verify())
	retAc := res.GetObject().GetAccount()
	assert.Equal(t, retAc.GetAccountId(), ac.AccountId)
	assert.Equal(t, len(retAc.GetPublicKeys()), 1)
	assert.Contains(t, retAc.GetPublicKeys(), ac.Pubkey)
}

func (am *AccountManager) QueryPeersState(t *testing.T, peers []model.PeerWithPriKey) {
	query := am.fc.NewQueryBuilder().
		AuthorizerId(am.authorizer.AccountId).
		FromId("/peer").
		CreatedTime(RandomNow()).
		RequestCode(model.ListObjectCode).
		Limit(10).
		Build()
	assert.NoError(t, query.Sign(am.authorizer.Pubkey, am.authorizer.Prikey))

	res, err := am.client.Read(query)
	require.NoError(t, err)

	fmt.Println("object:", res)
	assert.Equal(t, len(res.GetObject().GetList()), len(peers))
	pactive := make(map[string]bool)
	for _, o := range res.GetObject().GetList() {
		p := o.GetPeer()
		pactive[p.GetPeerId()] = p.GetActive()
		fmt.Println(p)
		assert.True(t, p.GetActive())
	}
	for _, p := range peers {
		_, ok := pactive[p.GetPeerId()]
		assert.True(t, ok)
	}
}
