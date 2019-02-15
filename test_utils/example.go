package test_utils

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/client"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/stretchr/testify/assert"
	"reflect"
)

type SenderManager struct {
	Client     core.APIClient
	Authorizer *AccountWithPri
	fc         model.ModelFactory
	conf       *config.Config
}

func RequireNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func validateEqualArgs(expected, actual interface{}) error {
	if isFunction(expected) || isFunction(actual) {
		return errors.New("cannot take func type as argument")
	}
	return nil
}

func isFunction(arg interface{}) bool {
	if arg == nil {
		return false
	}
	return reflect.TypeOf(arg).Kind() == reflect.Func
}

func formatUnequalValues(expected, actual interface{}) (e string, a string) {
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		return fmt.Sprintf("%T(%#v)", expected, expected),
			fmt.Sprintf("%T(%#v)", actual, actual)
	}

	return fmt.Sprintf("%#v", expected),
		fmt.Sprintf("%#v", actual)
}

func AssertEqual(expected interface{}, actual interface{}) {
	if err := validateEqualArgs(expected, actual); err != nil {
		panic(fmt.Sprintf("Invalid operation: %#v == %#v (%s)", expected, actual, err))
	}
	if !assert.ObjectsAreEqual(expected, actual) {
		expected, actual = formatUnequalValues(expected, actual)
		panic(fmt.Sprintf("Not equal: \n"+
			"expected: %s\n"+
			"actual  : %s\n", expected, actual))
	}
}

func ForceVerify(pubkey model.PublicKey, hasher model.Hasher, sig []byte) {
	RequireNoError(RandomCryptor().Verify(pubkey, hasher, sig))
}

func ForceSign(fc model.ModelFactory, pubkey model.PublicKey, prikey model.PrivateKey, hasher model.Hasher) model.Signature {
	signature, err := RandomCryptor().Sign(hasher, prikey)
	RequireNoError(err)
	return fc.NewSignature(pubkey, signature)
}

func NewSenderManager(authorizer *AccountWithPri, server model.Peer, fc model.ModelFactory, conf *config.Config) *SenderManager {
	c, err := client.NewAPIClient(server, fc)
	RequireNoError(err)
	return &SenderManager{
		c,
		authorizer,
		fc,
		conf,
	}
}

func (am *SenderManager) SetAuthorizer() {
	tx := am.fc.NewTxBuilder().
		AddPublicKeys(am.Authorizer.AccountId, am.Authorizer.AccountId, []model.PublicKey{am.Authorizer.Pubkey}).
		SetQuorum(am.Authorizer.AccountId, am.Authorizer.AccountId, 1).
		Build()
	fmt.Println(color.CyanString("SetAuthorizer: %+v", tx))
	RequireNoError(tx.Sign(am.Authorizer.Pubkey, am.Authorizer.Prikey))
	RequireNoError(am.Client.Write(tx))
}

func (am *SenderManager) CreateAccount(ac *AccountWithPri) {
	tx := am.fc.NewTxBuilder().
		CreateAccount(am.Authorizer.AccountId, ac.AccountId, []model.PublicKey{ac.Pubkey}, 1).
		Consign(am.Authorizer.AccountId, ac.AccountId, "root@peer").
		Build()
	fmt.Println(color.CyanString("CreateAccount: %+v", tx))
	RequireNoError(tx.Sign(am.Authorizer.Pubkey, am.Authorizer.Prikey))
	RequireNoError(am.Client.Write(tx))
}

func (am *SenderManager) AddPeer(peer model.Peer) {
	tx := am.fc.NewTxBuilder().
		AddPeer(am.Authorizer.AccountId, peer.GetPeerId(), peer.GetAddress(), peer.GetPublicKey()).
		Build()
	fmt.Println(color.CyanString("AddPeer: %+v", tx))
	RequireNoError(tx.Sign(am.Authorizer.Pubkey, am.Authorizer.Prikey))
	RequireNoError(am.Client.Write(tx))
}

func (am *SenderManager) Consign(ac *AccountWithPri, peer model.Peer) {
	tx := am.fc.NewTxBuilder().
		Consign(am.Authorizer.AccountId, ac.AccountId, peer.GetPeerId()).
		Build()
	fmt.Println(color.CyanString("Consign: %+v", tx))
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
	fmt.Println(color.CyanString("CreateEdgeStorage: %+v", tx))
	RequireNoError(tx.Sign(am.Authorizer.Pubkey, am.Authorizer.Prikey))
	RequireNoError(am.Client.Write(tx))
}

func (am *SenderManager) ProposeNewAlgorithm(pType string, prosl []byte) { //consensus []byte, incentive []byte) {
	var proslId string
	switch pType {
	case core.IncentiveKey:
		proslId = MakeIncentiveWalletId(am.Authorizer).Id()
	case core.ConsensusKey:
		proslId = MakeConsensusWalletId(am.Authorizer).Id()
	default:
		panic(fmt.Sprintf("Error pType: %s", pType))
	}
	tx := am.fc.NewTxBuilder().
		CreateStorage(am.Authorizer.AccountId, proslId).
		UpdateObject(am.Authorizer.AccountId, proslId,
			core.ProslTypeKey, am.fc.NewObjectBuilder().Str(pType)).
		UpdateObject(am.Authorizer.AccountId, proslId,
			core.ProslKey, am.fc.NewObjectBuilder().Data(prosl)).
		Build()
	fmt.Println(color.CyanString("ProposeNewAlgorithm: %+v", tx))
	RequireNoError(tx.Sign(am.Authorizer.Pubkey, am.Authorizer.Prikey))
	RequireNoError(am.Client.Write(tx))
}

func (am *SenderManager) CreateProslSignStorage() {
	tx := am.fc.NewTxBuilder().
		CreateStorage(am.Authorizer.AccountId, MakeIncentiveSigsId(am.Authorizer).Id()).
		CreateStorage(am.Authorizer.AccountId, MakeConsensusSigsId(am.Authorizer).Id()).
		Build()
	fmt.Println(color.CyanString("CreateProslSignStorage: %+v", tx))
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
	signature := ForceSign(am.fc, am.Authorizer.Pubkey, am.Authorizer.Prikey, prosl)
	tx := am.fc.NewTxBuilder().
		CreateStorage(am.Authorizer.AccountId, srcWalletId.Id()).
		AddObject(am.Authorizer.AccountId, srcWalletId.Id(), ProSignKey,
			am.fc.NewObjectBuilder().Sig(signature)).
		TransferObject(am.Authorizer.AccountId, srcWalletId.Id(), destWalletId.Id(),
			ProSignKey, am.fc.NewObjectBuilder().Sig(signature)).
		Build()
	fmt.Println(color.CyanString("VoteNewConsensus: %+v", tx))
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
	fmt.Println(color.CyanString("CheckAndCommit: %+v", tx))
	RequireNoError(tx.Sign(am.Authorizer.Pubkey, am.Authorizer.Prikey))
	RequireNoError(am.Client.Write(tx))
}

func (am *SenderManager) CheckAndCommitInc() {
	incId := MakeIncentiveWalletId(am.Authorizer)
	tx := am.fc.NewTxBuilder().
		CheckAndCommitProsl(am.Authorizer.AccountId, incId.Id(),
			map[string]model.Object{
				"account_id": am.fc.NewObjectBuilder().Address(fmt.Sprintf("%s@%s", incId.Account(), incId.Domain())),
			}).
		Build()
	fmt.Println(color.CyanString("CheckAndCommitInc: %+v", tx))
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

	AssertEqual(retAc.GetAccountId(), ac.AccountId)
	//assert.Contains(t, retAc.GetPublicKeys(), ac.Pubkey)
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
	res := am.queryRangeAccounts(fromId, 10)
	AssertEqual(len(res), len(acs))
	amap := make(map[string]struct{})
	for _, ac := range res {
		amap[ac.GetAccountId()] = struct{}{}
	}
	for _, ac := range acs {
		_, ok := amap[ac.AccountId]
		AssertEqual(true, ok)
	}
}

func (am *SenderManager) QueryAccountDegradedPassed(ac *AccountWithPri, peer model.Peer) {
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
	AssertEqual(retAc.GetAccountId(), ac.AccountId)
	AssertEqual(retAc.GetDelegatePeerId(), peer.GetPeerId())
}

func (am *SenderManager) QueryPeersStatePassed(peers []model.Peer) {
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

	AssertEqual(len(res.GetObject().GetList()), len(peers))
	pactive := make(map[string]bool)
	for _, o := range res.GetObject().GetList() {
		p := o.GetPeer()
		pactive[p.GetPeerId()] = p.GetActive()
		AssertEqual(true, p.GetActive())
	}
	for _, p := range peers {
		_, ok := pactive[p.GetPeerId()]
		AssertEqual(true, ok)
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
			panic(fmt.Sprintf("not exist hash: %x", a.Hash()))
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
	resSt := am.QueryStorage(fromId)
	equalList(resSt.GetFromKey(FollowEdge).GetList(), os)
}

func (am *SenderManager) QueryProslPassed(pType string, prosl []byte) {
	var proslId string
	switch pType {
	case core.IncentiveKey:
		proslId = MakeIncentiveWalletId(am.Authorizer).Id()
	case core.ConsensusKey:
		proslId = MakeConsensusWalletId(am.Authorizer).Id()
	default:
		panic(fmt.Sprintf("Error pType: %s", pType))
	}
	res := am.QueryStorage(proslId)
	AssertEqual(res.GetFromKey(core.ProslTypeKey).GetStr(), pType)
	AssertEqual(res.GetFromKey(core.ProslKey).GetData(), prosl)
}

func (am *SenderManager) QueryCollectSigsPassed(pType string, prosl model.Object, num int) {
	var sigsId string
	switch pType {
	case core.IncentiveKey:
		sigsId = MakeIncentiveSigsId(am.Authorizer).Id()
	case core.ConsensusKey:
		sigsId = MakeConsensusSigsId(am.Authorizer).Id()
	default:
		panic(fmt.Sprintf("Error pType: %s", pType))
	}
	res := am.QueryStorage(sigsId).GetFromKey(ProSignKey)
	AssertEqual(num, len(res.GetList()))
	for _, o := range res.GetList() {
		sig := o.GetSig()
		ForceVerify(sig.GetPublicKey(), prosl, sig.GetSignature())
	}
}

func (am *SenderManager) QueryRootProslPassed(prosl model.Storage) {
	var proslId string
	pType := prosl.GetFromKey(core.ProslTypeKey).GetStr()
	switch pType {
	case core.IncentiveKey:
		proslId = am.conf.Prosl.Incentive.Id
	case core.ConsensusKey:
		proslId = am.conf.Prosl.Consensus.Id
	default:
		panic(fmt.Sprintf("Error pType: %s", pType))
	}
	res := am.QueryStorage(proslId)
	AssertEqual(prosl.Hash(), res.Hash())
}

func (am *SenderManager) QueryAccountsBalances() {
	acs := am.queryRangeAccounts("creator.pr/account", 100)

	for _, ac := range acs {
		fmt.Println(color.YellowString("id: %s, balance: %d", ac.GetAccountId(), ac.GetBalance()))
	}

	pcs := am.queryRangeAccounts("peer/account", 100)
	for _, ac := range pcs {
		fmt.Println(color.YellowString("id: %s, balance: %d", ac.GetAccountId(), ac.GetBalance()))
	}
}
