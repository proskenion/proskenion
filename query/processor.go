package query

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"sort"
)

type QueryProcessor struct {
	fc   model.ModelFactory
	conf *config.Config
}

func NewQueryProcessor(fc model.ModelFactory, conf *config.Config) core.QueryProcessor {
	return &QueryProcessor{fc, conf}
}

func (q *QueryProcessor) Query(wsv model.ObjectFinder, query model.Query) (model.QueryResponse, error) {
	id := model.MustAddress(query.GetPayload().GetFromId())
	var object model.Object
	if id.Type() == model.WallettAddressType || query.GetPayload().GetRequestCode() != model.ListObjectCode {
		switch id.Storage() {
		case model.AccountStorageName:
			ac, err := q.accountObjectQuery(query.GetPayload(), wsv)
			if err != nil {
				return nil, err
			}
			object = q.selectAccount(ac, query)
		case model.PeerStorageName:
			peer, err := q.peerObjectQuery(query.GetPayload(), wsv)
			if err != nil {
				return nil, err
			}
			object = q.selectPeer(peer, query)
		default:
			storage, err := q.storageObjectQuery(query.GetPayload(), wsv)
			if err != nil {
				return nil, err
			}
			object = q.selectStorage(storage, query)
		}
	} else { // Range 検索
		obs := make([]model.Object, 0)
		switch id.Storage() {
		case model.AccountStorageName:
			acs, err := q.accountObjectQueryRange(query.GetPayload(), wsv)
			if err != nil {
				return nil, err
			}
			acs = q.limitAccounts(q.orderAccounts(q.whereAccounts(acs, query), query), query)
			for _, ac := range acs {
				obs = append(obs, q.selectAccount(ac, query))
			}
		case model.PeerStorageName:
			peers, err := q.peerObjectQueryRange(query.GetPayload(), wsv)
			if err != nil {
				return nil, err
			}
			peers = q.limitPeers(q.orderPeers(q.wherePeers(peers, query), query), query)
			for _, peer := range peers {
				obs = append(obs, q.selectPeer(peer, query))
			}
		default:
			storages, err := q.storageObjectQueryRange(query.GetPayload(), wsv)
			if err != nil {
				return nil, err
			}
			storages = q.limitStorages(q.orderStorages(q.whereStorages(storages, query), query), query)
			for _, storage := range storages {
				obs = append(obs, q.selectStorage(storage, query))
			}
		}
		object = q.fc.NewObjectBuilder().List(obs)
	}
	ret := q.fc.NewQueryResponseBuilder().Object(object).Build()
	if err := q.signedResponse(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (q *QueryProcessor) accountObjectQuery(qp model.QueryPayload, wsv model.ObjectFinder) (model.Account, error) {
	ac := q.fc.NewEmptyAccount()
	err := wsv.Query(model.MustAddress(qp.GetFromId()), ac)
	if err != nil {
		return nil, errors.Wrap(core.ErrQueryProcessorNotFound, err.Error())
	}
	return ac, nil
}

func (q *QueryProcessor) peerObjectQuery(qp model.QueryPayload, wsv model.ObjectFinder) (model.Peer, error) {
	peer := q.fc.NewEmptyPeer()
	err := wsv.Query(model.MustAddress(qp.GetFromId()), peer)
	if err != nil {
		return nil, errors.Wrap(core.ErrQueryProcessorNotFound, err.Error())
	}
	return peer, nil
}

func (q *QueryProcessor) storageObjectQuery(qp model.QueryPayload, wsv model.ObjectFinder) (model.Storage, error) {
	storage := q.fc.NewEmptyStorage()
	err := wsv.Query(model.MustAddress(qp.GetFromId()), storage)
	if err != nil {
		return nil, errors.Wrap(core.ErrQueryProcessorNotFound, err.Error())
	}
	return storage, nil
}

func (q *QueryProcessor) selectAccount(ac model.Account, query model.Query) model.Object {
	builder := q.fc.NewObjectBuilder()
	if query.GetPayload().GetSelect() != "*" {
		ret := ac.GetFromKey(query.GetPayload().GetSelect())
		if ret.GetType() != model.AnythingObjectCode {
			return ret
		}
	}
	return builder.Account(ac)
}

func (q *QueryProcessor) selectPeer(peer model.Peer, query model.Query) model.Object {
	builder := q.fc.NewObjectBuilder()
	if query.GetPayload().GetSelect() != "*" {
		ret := peer.GetFromKey(query.GetPayload().GetSelect())
		if ret.GetType() != model.AnythingObjectCode {
			return ret
		}
	}
	return builder.Peer(peer)
}

func (q *QueryProcessor) selectStorage(storage model.Storage, query model.Query) model.Object {
	builder := q.fc.NewObjectBuilder()
	if ret, ok := storage.GetObject()[query.GetPayload().GetSelect()]; ok {
		return ret
	}
	return builder.Storage(storage)
}

func (q *QueryProcessor) accountObjectQueryRange(qp model.QueryPayload, wsv model.ObjectFinder) ([]model.Account, error) {
	acs := make([]model.Account, 0)
	res, err := wsv.QueryAll(model.MustAddress(qp.GetFromId()), model.NewAccountUnmarshalerFactory(q.fc))
	if err != nil {
		return nil, errors.Wrap(core.ErrQueryProcessorNotFound, err.Error())
	}
	for _, r := range res {
		acs = append(acs, r.(model.Account))
	}
	return acs, nil
}

func (q *QueryProcessor) peerObjectQueryRange(qp model.QueryPayload, wsv model.ObjectFinder) ([]model.Peer, error) {
	peers := make([]model.Peer, 0)
	res, err := wsv.QueryAll(model.MustAddress(qp.GetFromId()), model.NewPeerUnmarshalerFactory(q.fc))
	if err != nil {
		return nil, errors.Wrap(core.ErrQueryProcessorNotFound, err.Error())
	}
	for _, r := range res {
		peers = append(peers, r.(model.Peer))
	}
	return peers, nil
}

func (q *QueryProcessor) storageObjectQueryRange(qp model.QueryPayload, wsv model.ObjectFinder) ([]model.Storage, error) {
	storages := make([]model.Storage, 0)
	res, err := wsv.QueryAll(model.MustAddress(qp.GetFromId()), model.NewStorageUnmarshalerFactory(q.fc))
	if err != nil {
		return nil, errors.Wrap(core.ErrQueryProcessorNotFound, err.Error())
	}
	for _, r := range res {
		storages = append(storages, r.(model.Storage))
	}
	return storages, nil
}

func (q *QueryProcessor) whereAccounts(acs []model.Account, query model.Query) []model.Account {
	return acs
}

func (q *QueryProcessor) wherePeers(peers []model.Peer, query model.Query) []model.Peer {
	return peers
}

func (q *QueryProcessor) whereStorages(storages []model.Storage, query model.Query) []model.Storage {
	return storages
}

type Accounts struct {
	acs []model.Account
	key string
}

func (a *Accounts) Len() int {
	return len(a.acs)
}

func (a *Accounts) Less(i, j int) bool {
	if !bytes.Equal(a.acs[i].GetFromKey(a.key).Hash(), a.acs[j].GetFromKey(a.key).Hash()) {
		switch a.key {
		case "id":
			return a.acs[i].GetAccountId() < a.acs[j].GetAccountId()
		case "name":
			return a.acs[i].GetAccountName() < a.acs[j].GetAccountName()
		case "balance":
			return a.acs[i].GetBalance() < a.acs[j].GetBalance()
		case "quorum":
			return a.acs[i].GetQuorum() < a.acs[j].GetQuorum()
		case "peer":
			return a.acs[i].GetDelegatePeerId() < a.acs[j].GetDelegatePeerId()
		}
	}
	return a.acs[i].GetAccountId() < a.acs[j].GetAccountId()
}

func (a *Accounts) Swap(i, j int) {
	a.acs[i], a.acs[j] = a.acs[j], a.acs[i]
}

func (q *QueryProcessor) orderAccounts(acs []model.Account, query model.Query) []model.Account {
	as := &Accounts{acs, query.GetPayload().GetOrderBy().GetKey()}
	sort.Sort(as)
	if query.GetPayload().GetOrderBy().GetOrder() == model.DESC {
		for i, j := 0, len(as.acs)-1; i < j; i, j = i+1, j-1 {
			as.acs[i], as.acs[j] = as.acs[j], as.acs[i]
		}
	}
	return as.acs
}

type Peers struct {
	peers []model.Peer
	key   string
}

func (a *Peers) Len() int {
	return len(a.peers)
}

func (a *Peers) Less(i, j int) bool {
	if !bytes.Equal(a.peers[i].GetFromKey(a.key).Hash(), a.peers[j].GetFromKey(a.key).Hash()) {
		switch a.key {
		case "id":
			return a.peers[i].GetPeerId() < a.peers[j].GetPeerId()
		case "address":
			return a.peers[i].GetAddress() < a.peers[j].GetAddress()
		}
	}
	return a.peers[i].GetPeerId() < a.peers[j].GetPeerId()
}

func (a *Peers) Swap(i, j int) {
	a.peers[i], a.peers[j] = a.peers[j], a.peers[i]
}

func (q *QueryProcessor) orderPeers(peers []model.Peer, query model.Query) []model.Peer {
	as := &Peers{peers, query.GetPayload().GetOrderBy().GetKey()}
	sort.Sort(as)
	if query.GetPayload().GetOrderBy().GetOrder() == model.DESC {
		for i, j := 0, len(as.peers)-1; i < j; i, j = i+1, j-1 {
			as.peers[i], as.peers[j] = as.peers[j], as.peers[i]
		}
	}
	return as.peers
}

type Storages struct {
	storages []model.Storage
	key      string
}

func (a *Storages) Len() int {
	return len(a.storages)
}

func (a *Storages) Less(i, j int) bool {
	if v1, ok := a.storages[i].GetObject()[a.key]; ok {
		if v2, ok := a.storages[j].GetObject()[a.key]; ok {
			return model.ObjectLess(v1, v2)
		}
	}
	return model.HasherLess(a.storages[i], a.storages[j])
}

func (a *Storages) Swap(i, j int) {
	a.storages[i], a.storages[j] = a.storages[j], a.storages[i]
}

func (q *QueryProcessor) orderStorages(storages []model.Storage, query model.Query) []model.Storage {
	as := &Storages{storages, query.GetPayload().GetOrderBy().GetKey()}
	sort.Sort(as)
	if query.GetPayload().GetOrderBy().GetOrder() == model.DESC {
		for i, j := 0, len(as.storages)-1; i < j; i, j = i+1, j-1 {
			as.storages[i], as.storages[j] = as.storages[j], as.storages[i]
		}
	}
	return as.storages
}

func (q *QueryProcessor) limitAccounts(acs []model.Account, query model.Query) []model.Account {
	if int32(len(acs)) < query.GetPayload().GetLimit() {
		return acs
	}
	return acs[:query.GetPayload().GetLimit()]
}

func (q *QueryProcessor) limitPeers(peers []model.Peer, query model.Query) []model.Peer {
	if int32(len(peers)) < query.GetPayload().GetLimit() {
		return peers
	}
	return peers[:query.GetPayload().GetLimit()]
}

func (q *QueryProcessor) limitStorages(storages []model.Storage, query model.Query) []model.Storage {
	if int32(len(storages)) < query.GetPayload().GetLimit() {
		return storages
	}
	return storages[:query.GetPayload().GetLimit()]
}

func (q *QueryProcessor) signedResponse(res model.QueryResponse) error {
	return res.Sign(q.conf.Peer.PublicKeyBytes(), q.conf.Peer.PrivateKeyBytes())
}
