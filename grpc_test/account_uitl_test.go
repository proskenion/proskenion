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

func NewAccountManager(t *testing.T, authorizer *AccountWithPri, server model.Peer) *AccountManager {
	fc := RandomFactory()
	c, err := client.NewAPIClient(server, fc)
	require.NoError(t, err)
	return &AccountManager{
		c,
		authorizer,
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

func (am *AccountManager) Consign(t *testing.T, ac *AccountWithPri, peer model.Peer) {
	tx := am.fc.NewTxBuilder().
		Consign(am.authorizer.AccountId, ac.AccountId, peer.GetPeerId()).
		Build()
	require.NoError(t, tx.Sign(am.authorizer.Pubkey, am.authorizer.Prikey))
	require.NoError(t, am.client.Write(tx))
}

const (
	FollowStorage = "follow"
	FollowEdge    = "to"
	ProslStorage  = "prosl"
)

func (am *AccountManager) AddEdge(t *testing.T, ac *AccountWithPri, to *AccountWithPri) {
	obj := am.fc.NewObjectBuilder().Address(to.AccountId)
	tx := am.fc.NewTxBuilder().
		AddObject(am.authorizer.AccountId, fmt.Sprintf("%s/%s", ac.AccountId, FollowStorage), FollowEdge, obj).
		Build()
	require.NoError(t, tx.Sign(am.authorizer.Pubkey, am.authorizer.Prikey))
	require.NoError(t, am.client.Write(tx))
}

func (am *AccountManager) CreateEdgeStorage(t *testing.T, ac *AccountWithPri) {
	tx := am.fc.NewTxBuilder().
		CreateStorage(am.authorizer.AccountId, fmt.Sprintf("%s/%s", am.authorizer.AccountId, FollowStorage)).
		Build()
	require.NoError(t, tx.Sign(am.authorizer.Pubkey, am.authorizer.Prikey))
	require.NoError(t, am.client.Write(tx))
}

func (am *AccountManager) ProposeNewConsensus(t *testing.T, consensus []byte, incentive []byte) {
	id := model.MustAddress(am.authorizer.AccountId)
	IncentiveId := fmt.Sprintf("%s@%s.%s/%s", id.Account(), core.IncentiveKey, id.Domain(), ProslStorage)
	ConsensusId := fmt.Sprintf("%s@%s.%s/%s", id.Account(), core.ConsensusKey, id.Domain(), ProslStorage)
	tx := am.fc.NewTxBuilder().
		CreateStorage(am.authorizer.AccountId, IncentiveId).
		CreateStorage(am.authorizer.AccountId, ConsensusId).
		UpdateObject(am.authorizer.AccountId, IncentiveId,
			core.ProslTypeKey, am.fc.NewObjectBuilder().Str(core.IncentiveKey)).
		UpdateObject(am.authorizer.AccountId, ConsensusId,
			core.ProslTypeKey, am.fc.NewObjectBuilder().Str(core.ConsensusKey)).
		UpdateObject(am.authorizer.AccountId, IncentiveId,
			core.ProslKey, am.fc.NewObjectBuilder().Data(incentive)).
		UpdateObject(am.authorizer.AccountId, ConsensusId,
			core.ProslKey, am.fc.NewObjectBuilder().Data(consensus)).
		Build()
	require.NoError(t, tx.Sign(am.authorizer.Pubkey, am.authorizer.Prikey))
	require.NoError(t, am.client.Write(tx))
}

func (am *AccountManager) QueryAccountPassed(t *testing.T, ac *AccountWithPri) {
	query := am.fc.NewQueryBuilder().
		AuthorizerId(am.authorizer.AccountId).
		FromId(model.MustAddress(ac.AccountId).AccountId()).
		CreatedTime(RandomNow()).
		RequestCode(model.AccountObjectCode).
		Build()
	assert.NoError(t, query.Sign(am.authorizer.Pubkey, am.authorizer.Prikey))

	res, err := am.client.Read(query)
	assert.NoError(t, err)

	assert.NoError(t, res.Verify())
	retAc := res.GetObject().GetAccount()
	assert.Equal(t, retAc.GetAccountId(), ac.AccountId)
	assert.Equal(t, len(retAc.GetPublicKeys()), 1)
	assert.Contains(t, retAc.GetPublicKeys(), ac.Pubkey)
}

func (am *AccountManager) queryRangeAccounts(t *testing.T, fromId string, limit int32) []model.Account {
	query := am.fc.NewQueryBuilder().
		AuthorizerId(am.authorizer.AccountId).
		FromId(fromId).
		CreatedTime(RandomNow()).
		RequestCode(model.ListObjectCode).
		Limit(limit).
		Build()
	assert.NoError(t, query.Sign(am.authorizer.Pubkey, am.authorizer.Prikey))

	res, err := am.client.Read(query)
	require.NoError(t, err)

	ret := make([]model.Account, 0)
	for _, o := range res.GetObject().GetList() {
		ret = append(ret, o.GetAccount())
	}
	return ret
}

func (am *AccountManager) QueryRangeAccountsPassed(t *testing.T, fromId string, acs []*AccountWithPri) {
	res := am.queryRangeAccounts(t, fromId, 10)
	assert.Equal(t, len(res), len(acs))
	amap := make(map[string]struct{})
	for _, ac := range res {
		amap[ac.GetAccountId()] = struct{}{}
	}
	for _, ac := range acs {
		_, ok := amap[ac.AccountId]
		require.True(t, ok)
	}
}

func (am *AccountManager) QueryAccountDegradedPassed(t *testing.T, ac *AccountWithPri, peer model.PeerWithPriKey) {
	query := am.fc.NewQueryBuilder().
		AuthorizerId(am.authorizer.AccountId).
		FromId(model.MustAddress(ac.AccountId).AccountId()).
		CreatedTime(RandomNow()).
		RequestCode(model.AccountObjectCode).
		Build()
	assert.NoError(t, query.Sign(am.authorizer.Pubkey, am.authorizer.Prikey))

	res, err := am.client.Read(query)
	assert.NoError(t, err)

	assert.NoError(t, res.Verify())
	retAc := res.GetObject().GetAccount()
	assert.Equal(t, retAc.GetAccountId(), ac.AccountId)
	assert.Equal(t, retAc.GetDelegatePeerId(), peer.GetPeerId())
}

func (am *AccountManager) QueryPeersStatePassed(t *testing.T, peers []model.PeerWithPriKey) {
	query := am.fc.NewQueryBuilder().
		AuthorizerId(am.authorizer.AccountId).
		FromId("/" + model.PeerStorageName).
		CreatedTime(RandomNow()).
		RequestCode(model.ListObjectCode).
		Limit(10).
		Build()
	assert.NoError(t, query.Sign(am.authorizer.Pubkey, am.authorizer.Prikey))

	res, err := am.client.Read(query)
	require.NoError(t, err)

	assert.Equal(t, len(res.GetObject().GetList()), len(peers))
	pactive := make(map[string]bool)
	for _, o := range res.GetObject().GetList() {
		p := o.GetPeer()
		pactive[p.GetPeerId()] = p.GetActive()
		assert.True(t, p.GetActive())
	}
	for _, p := range peers {
		_, ok := pactive[p.GetPeerId()]
		assert.True(t, ok)
	}
}

func equalList(t *testing.T, os []model.Object, as []model.Object) {
	h := make(map[string]struct{})
	for _, o := range os {
		h[string(o.Hash())] = struct{}{}
	}
	for _, a := range as {
		_, ok := h[string(a.Hash())]
		if !ok {
			t.Fatalf("not exist hash: %x", a.Hash())
		}
	}
}

func (am *AccountManager) queryStorage(t *testing.T, fromId string) model.Storage {
	query := am.fc.NewQueryBuilder().
		AuthorizerId(am.authorizer.AccountId).
		FromId(fromId).
		CreatedTime(RandomNow()).
		RequestCode(model.StorageAddressType).
		Build()
	assert.NoError(t, query.Sign(am.authorizer.Pubkey, am.authorizer.Prikey))

	res, err := am.client.Read(query)
	require.NoError(t, err)
	return res.GetObject().GetStorage()
}

func (am *AccountManager) QueryStorageEdgesPassed(t *testing.T, fromId string, os []model.Object) {
	resSt := am.queryStorage(t, fromId)
	equalList(t, resSt.GetFromKey(FollowEdge).GetList(), os)
}

func (am *AccountManager) QueryProslPassed(t *testing.T, pType string, prosl []byte) {
	id := model.MustAddress(am.authorizer.AccountId)
	var proslId string
	switch pType {
	case core.IncentiveKey:
		proslId = fmt.Sprintf("%s@%s.%s/%s", id.Account(), core.IncentiveKey, id.Domain(), ProslStorage)
	case core.ConsensusKey:
		proslId = fmt.Sprintf("%s@%s.%s/%s", id.Account(), core.ConsensusKey, id.Domain(), ProslStorage)
	default:
		require.Failf(t, "Error pType: %s", pType)
	}
	res := am.queryStorage(t, proslId)
	assert.Equal(t, res.GetFromKey(core.ProslTypeKey).GetStr(), pType)
	assert.Equal(t, res.GetFromKey(core.ProslKey).GetData(), prosl)
}
