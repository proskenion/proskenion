package main

import (
	"fmt"
	"github.com/proskenion/proskenion/Client"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type SenderManager struct {
	Client     core.APIClient
	Authorizer *AccountWithPri
	fc         model.ModelFactory
}

func RequireNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func NewSenderManager(authorizer *AccountWithPri, server model.Peer) *SenderManager {
	fc := RandomFactory()
	c, err := client.NewAPIClient(server, fc)
	RequireNoError(err)
	return &SenderManager{
		c,
		authorizer,
		fc,
	}
}

func (am *SenderManager) SetAuthorizer() {
	tx := am.fc.NewTxBuilder().
		AddPublicKeys(am.Authorizer.AccountId, am.Authorizer.AccountId, []model.PublicKey{am.Authorizer.Pubkey}).
		SetQuorum(am.Authorizer.AccountId, am.Authorizer.AccountId, 1).
		Build()
	RequireNoError(tx.Sign(am.Authorizer.Pubkey, am.Authorizer.Prikey))
	RequireNoError(am.Client.Write(tx))
}

func (am *SenderManager) CreateAccount(ac *AccountWithPri) {
	tx := am.fc.NewTxBuilder().
		CreateAccount(am.Authorizer.AccountId, ac.AccountId, []model.PublicKey{ac.Pubkey}, 1).
		Consign(am.Authorizer.AccountId, ac.AccountId, "root@peer").
		Build()
	RequireNoError(tx.Sign(am.Authorizer.Pubkey, am.Authorizer.Prikey))
	RequireNoError(am.Client.Write(tx))
}

func (am *SenderManager) AddPeer(peer model.Peer) {
	tx := am.fc.NewTxBuilder().
		CreateAccount(am.Authorizer.AccountId, peer.GetPeerId(), []model.PublicKey{peer.GetPublicKey()}, 1).
		AddPeer(am.Authorizer.AccountId, peer.GetPeerId(), peer.GetAddress(), peer.GetPublicKey()).
		Build()
	RequireNoError(tx.Sign(am.Authorizer.Pubkey, am.Authorizer.Prikey))
	RequireNoError(am.Client.Write(tx))
}

func (am *SenderManager) Consign(ac *AccountWithPri, peer model.Peer) {
	tx := am.fc.NewTxBuilder().
		Consign(am.Authorizer.AccountId, ac.AccountId, peer.GetPeerId()).
		Build()
	RequireNoError(tx.Sign(am.Authorizer.Pubkey, am.Authorizer.Prikey))
	RequireNoError(am.Client.Write(tx))
}

const (
	FollowStorage  = "follow"
	FollowEdge     = "to"
	ProslStorage   = "prosl"
	ProSignStorage = "prsign"
	ProSignKey     = "sigs"
)

func MakeConsensusWalletId(ac *AccountWithPri) model.Address {
	id := model.MustAddress(ac.AccountId)
	return model.MustAddress(fmt.Sprintf("%s@%s.%s/%s", id.Account(), core.ConsensusKey, id.Domain(), ProslStorage))
}

func MakeIncentiveWalletId(ac *AccountWithPri) model.Address {
	id := model.MustAddress(ac.AccountId)
	return model.MustAddress(fmt.Sprintf("%s@%s.%s/%s", id.Account(), core.IncentiveKey, id.Domain(), ProslStorage))
}

func MakeConsensusSigsId(ac *AccountWithPri) model.Address {
	id := model.MustAddress(ac.AccountId)
	return model.MustAddress(fmt.Sprintf("%s@%s.%s/%s", id.Account(), core.ConsensusKey, id.Domain(), ProSignStorage))
}

func MakeIncentiveSigsId(ac *AccountWithPri) model.Address {
	id := model.MustAddress(ac.AccountId)
	return model.MustAddress(fmt.Sprintf("%s@%s.%s/%s", id.Account(), core.IncentiveKey, id.Domain(), ProSignStorage))
}

func (am *SenderManager) AddEdge(ac *AccountWithPri, to *AccountWithPri) {
	obj := am.fc.NewObjectBuilder().Address(to.AccountId)
	tx := am.fc.NewTxBuilder().
		AddObject(am.Authorizer.AccountId, fmt.Sprintf("%s/%s", ac.AccountId, FollowStorage), FollowEdge, obj).
		Build()
	RequireNoError(tx.Sign(am.Authorizer.Pubkey, am.Authorizer.Prikey))
	RequireNoError(am.Client.Write(tx))
}

func (am *SenderManager) CreateEdgeStorage(ac *AccountWithPri) {
	tx := am.fc.NewTxBuilder().
		CreateStorage(am.Authorizer.AccountId, fmt.Sprintf("%s/%s", am.Authorizer.AccountId, FollowStorage)).
		Build()
	RequireNoError(tx.Sign(am.Authorizer.Pubkey, am.Authorizer.Prikey))
	RequireNoError(am.Client.Write(tx))
}

func (am *SenderManager) ProposeNewConsensus(consensus []byte, incentive []byte) {
	IncentiveId := MakeIncentiveWalletId(am.Authorizer).Id()
	ConsensusId := MakeConsensusWalletId(am.Authorizer).Id()
	tx := am.fc.NewTxBuilder().
		CreateStorage(am.Authorizer.AccountId, IncentiveId).
		CreateStorage(am.Authorizer.AccountId, ConsensusId).
		UpdateObject(am.Authorizer.AccountId, IncentiveId,
			core.ProslTypeKey, am.fc.NewObjectBuilder().Str(core.IncentiveKey)).
		UpdateObject(am.Authorizer.AccountId, ConsensusId,
			core.ProslTypeKey, am.fc.NewObjectBuilder().Str(core.ConsensusKey)).
		UpdateObject(am.Authorizer.AccountId, IncentiveId,
			core.ProslKey, am.fc.NewObjectBuilder().Data(incentive)).
		UpdateObject(am.Authorizer.AccountId, ConsensusId,
			core.ProslKey, am.fc.NewObjectBuilder().Data(consensus)).
		Build()
	RequireNoError(tx.Sign(am.Authorizer.Pubkey, am.Authorizer.Prikey))
	RequireNoError(am.Client.Write(tx))
}

func (am *SenderManager) CreateProslSignStorage() {
	tx := am.fc.NewTxBuilder().
		CreateStorage(am.Authorizer.AccountId, MakeIncentiveSigsId(am.Authorizer).Id()).
		CreateStorage(am.Authorizer.AccountId, MakeConsensusSigsId(am.Authorizer).Id()).
		Build()
	RequireNoError(tx.Sign(am.Authorizer.Pubkey, am.Authorizer.Prikey))
	RequireNoError(am.Client.Write(tx))
}

func (am *SenderManager) VoteNewConsensus(dest *AccountWithPri, key string, prosl model.Object) {
	var destWalletId model.Address
	var srcWalletId model.Address
	switch key {
	case core.IncentiveKey:
		destWalletId = MakeIncentiveSigsId(dest)
		srcWalletId = MakeIncentiveSigsId(am.Authorizer)
	case core.ConsensusKey:
		destWalletId = MakeConsensusSigsId(dest)
		srcWalletId = MakeConsensusSigsId(am.Authorizer)
	default:
		panic(fmt.Sprintf("Error pType: %s", key))
	}
	sign, err := RandomCryptor().Sign(prosl, am.Authorizer.Prikey)
	RequireNoError(err)
	signature := am.fc.NewSignature(am.Authorizer.Pubkey, sign)
	tx := am.fc.NewTxBuilder().
		CreateStorage(am.Authorizer.AccountId, srcWalletId.Id()).
		AddObject(am.Authorizer.AccountId, srcWalletId.Id(), ProSignKey,
			am.fc.NewObjectBuilder().Sig(signature)).
		TransferObject(am.Authorizer.AccountId, srcWalletId.Id(), destWalletId.Id(),
			ProSignKey, am.fc.NewObjectBuilder().Sig(signature)).
		Build()
	RequireNoError(tx.Sign(am.Authorizer.Pubkey, am.Authorizer.Prikey))
	RequireNoError(am.Client.Write(tx))
}

func (am *SenderManager) CheckAndCommit() {
	incId := MakeIncentiveWalletId(am.Authorizer)
	conId := MakeConsensusWalletId(am.Authorizer)
	tx := am.fc.NewTxBuilder().
		CheckAndCommitProsl(am.Authorizer.AccountId, incId.Id(),
			map[string]model.Object{
				"account_id": am.fc.NewObjectBuilder().Address(fmt.Sprintf("%s@%s", incId.Account(), incId.Domain())),
			}).
		CheckAndCommitProsl(am.Authorizer.AccountId, conId.Id(),
			map[string]model.Object{
				"account_id": am.fc.NewObjectBuilder().Address(fmt.Sprintf("%s@%s", conId.Account(), conId.Domain())),
			}).
		Build()
	RequireNoError(tx.Sign(am.Authorizer.Pubkey, am.Authorizer.Prikey))
	RequireNoError(am.Client.Write(tx))
}

func (am *SenderManager) QueryAccountPassed(ac *AccountWithPri) {
	query := am.fc.NewQueryBuilder().
		AuthorizerId(am.Authorizer.AccountId).
		FromId(model.MustAddress(ac.AccountId).AccountId()).
		CreatedTime(RandomNow()).
		RequestCode(model.AccountObjectCode).
		Build()
	RequireNoError(query.Sign(am.Authorizer.Pubkey, am.Authorizer.Prikey))

	res, err := am.Client.Read(query)
	RequireNoError(err)

	RequireNoError(res.Verify())
	retAc := res.GetObject().GetAccount()
}

func (am *SenderManager) queryRangeAccounts(fromId string, limit int32) []model.Account {
	query := am.fc.NewQueryBuilder().
		AuthorizerId(am.Authorizer.AccountId).
		FromId(fromId).
		CreatedTime(RandomNow()).
		RequestCode(model.ListObjectCode).
		Limit(limit).
		Build()
	RequireNoError(query.Sign(am.Authorizer.Pubkey, am.Authorizer.Prikey))

	res, err := am.Client.Read(query)
	RequireNoError(err)

	ret := make([]model.Account, 0)
	for _, o := range res.GetObject().GetList() {
		ret = append(ret, o.GetAccount())
	}
	return ret
}

func (am *SenderManager) QueryRangeAccountsPassed(fromId string, acs []*AccountWithPri) {
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

func (am *SenderManager) QueryAccountDegradedPassed(ac *AccountWithPri, peer model.PeerWithPriKey) {
	query := am.fc.NewQueryBuilder().
		AuthorizerId(am.Authorizer.AccountId).
		FromId(model.MustAddress(ac.AccountId).AccountId()).
		CreatedTime(RandomNow()).
		RequestCode(model.AccountObjectCode).
		Build()
	RequireNoError(query.Sign(am.Authorizer.Pubkey, am.Authorizer.Prikey))

	res, err := am.Client.Read(query)
	RequireNoError(err)

	RequireNoError(res.Verify())
	retAc := res.GetObject().GetAccount()
	assert.Equal(t, retAc.GetAccountId(), ac.AccountId)
	assert.Equal(t, retAc.GetDelegatePeerId(), peer.GetPeerId())
}

func (am *SenderManager) QueryPeersStatePassed(peers []model.PeerWithPriKey) {
	query := am.fc.NewQueryBuilder().
		AuthorizerId(am.Authorizer.AccountId).
		FromId("/" + model.PeerStorageName).
		CreatedTime(RandomNow()).
		RequestCode(model.ListObjectCode).
		Limit(10).
		Build()
	RequireNoError(query.Sign(am.Authorizer.Pubkey, am.Authorizer.Prikey))

	res, err := am.Client.Read(query)
	RequireNoError(err)

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

func equalList(os []model.Object, as []model.Object) {
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

func (am *SenderManager) QueryStorage(fromId string) model.Storage {
	query := am.fc.NewQueryBuilder().
		AuthorizerId(am.Authorizer.AccountId).
		FromId(fromId).
		CreatedTime(RandomNow()).
		RequestCode(model.StorageObjectCode).
		Build()
	RequireNoError(query.Sign(am.Authorizer.Pubkey, am.Authorizer.Prikey))

	res, err := am.Client.Read(query)
	RequireNoError(err)
	return res.GetObject().GetStorage()
}

func (am *SenderManager) QueryStorageEdgesPassed(fromId string, os []model.Object) {
	resSt := am.QueryStorage(t, fromId)
	equalList(t, resSt.GetFromKey(FollowEdge).GetList(), os)
}

func (am *SenderManager) QueryProslPassed(pType string, prosl []byte) {
	var proslId string
	switch pType {
	case core.IncentiveKey:
		proslId = MakeIncentiveWalletId(am.Authorizer).Id()
	case core.ConsensusKey:
		proslId = MakeConsensusWalletId(am.Authorizer).Id()
	default:
		require.Failf(t, "Error pType: %s", pType)
	}
	res := am.QueryStorage(t, proslId)
	assert.Equal(t, res.GetFromKey(core.ProslTypeKey).GetStr(), pType)
	assert.Equal(t, res.GetFromKey(core.ProslKey).GetData(), prosl)
}

func (am *SenderManager) QueryCollectSigsPassed(pType string, prosl model.Object, num int) {
	var sigsId string
	switch pType {
	case core.IncentiveKey:
		sigsId = MakeIncentiveSigsId(am.Authorizer).Id()
	case core.ConsensusKey:
		sigsId = MakeConsensusSigsId(am.Authorizer).Id()
	default:
		require.Failf(t, "Error pType: %s", pType)
	}
	res := am.QueryStorage(t, sigsId).GetFromKey(ProSignKey)
	assert.Equal(t, num, len(res.GetList()))
	for _, o := range res.GetList() {
		sig := o.GetSig()
		RandomVerify(t, sig.GetPublicKey(), prosl, sig.GetSignature())
	}
}

func (am *SenderManager) QueryRootProslPassed(prosl model.Storage) {
	var proslId string
	pType := prosl.GetFromKey(core.ProslTypeKey).GetStr()
	switch pType {
	case core.IncentiveKey:
		proslId = RandomConfig().Prosl.Incentive.Id
	case core.ConsensusKey:
		proslId = RandomConfig().Prosl.Consensus.Id
	default:
		require.Failf(t, "Error pType: %s", pType)
	}
	res := am.QueryStorage(t, proslId)
	assert.Equal(t, prosl.Hash(), res.Hash())
	fmt.Println("QueryRootProsl:", res)
}
