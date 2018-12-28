package query

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type QueryProcessor struct {
	rp   core.Repository
	fc   model.ModelFactory
	conf *config.Config
}

func NewQueryProcessor(rp core.Repository, fc model.ModelFactory, conf *config.Config) core.QueryProcessor {
	return &QueryProcessor{rp, fc, conf}
}

func (q *QueryProcessor) Query(query model.Query) (model.QueryResponse, error) {
	top, ok := q.rp.Top()
	if !ok {
		return nil, core.ErrQueryProcessorQueryEmptyBlockchain
	}
	rtx, err := q.rp.Begin()
	if err != nil {
		return nil, err
	}
	wsv, err := rtx.WSV(top.GetPayload().GetWSVHash())
	if err != nil {
		return nil, err
	}

	var res model.QueryResponse
	code := query.GetPayload().GetRequestCode()
	switch code {
	case model.AccountObjectCode:
		res, err = q.accountObjectQuery(query.GetPayload(), wsv)
	case model.PeerObjectCode:
		res, err = q.peerObjectQuery(query.GetPayload(), wsv)
	case model.StorageObjectCode:
		res, err = q.storageObjectQuery(query.GetPayload(), wsv)
	case model.ListObjectCode:
		res, err = q.listObjectQuery(query.GetPayload(), wsv)
	default:
		err = core.ErrQueryProcessorQueryObjectCodeNotImplemented
	}
	if errors.Cause(err) == core.ErrWSVNotFound {
		return nil, errors.Wrap(core.ErrQueryProcessorNotFound, err.Error())
	}
	return res, err
}

func (q *QueryProcessor) accountObjectQuery(qp model.QueryPayload, wsv core.WSV) (model.QueryResponse, error) {
	ac := q.fc.NewEmptyAccount()
	err := wsv.Query(model.MustAddress(qp.GetFromId()), ac)
	if err != nil {
		return nil, err
	}

	qr := q.fc.NewQueryResponseBuilder().
		Account(ac).
		Build()
	if err := q.signedResponse(qr); err != nil {
		return nil, err
	}
	return qr, nil
}

func (q *QueryProcessor) peerObjectQuery(qp model.QueryPayload, wsv core.WSV) (model.QueryResponse, error) {
	peer := q.fc.NewEmptyPeer()
	err := wsv.Query(model.MustAddress(qp.GetFromId()), peer)
	if err != nil {
		return nil, err
	}

	qr := q.fc.NewQueryResponseBuilder().
		Peer(peer).
		Build()
	if err := q.signedResponse(qr); err != nil {
		return nil, err
	}
	return qr, nil
}

func (q *QueryProcessor) storageObjectQuery(qp model.QueryPayload, wsv core.WSV) (model.QueryResponse, error) {
	storage := q.fc.NewEmptyStorage()
	err := wsv.Query(model.MustAddress(qp.GetFromId()), storage)
	if err != nil {
		return nil, err
	}

	qr := q.fc.NewQueryResponseBuilder().
		Storage(storage).
		Build()
	if err := q.signedResponse(qr); err != nil {
		return nil, err
	}
	return qr, nil
}

type AccountUnmarshalerFactory struct {
	fc model.ModelFactory
}

func (f *AccountUnmarshalerFactory) CreateUnmarshaler() model.Unmarshaler {
	return f.fc.NewEmptyAccount()
}

func NewAccountUnmarshalerFactory(fc model.ModelFactory) model.UnmarshalerFactory {
	return &AccountUnmarshalerFactory{fc}
}

type PeerUnmarshalerFactory struct {
	fc model.ModelFactory
}

func (f *PeerUnmarshalerFactory) CreateUnmarshaler() model.Unmarshaler {
	return f.fc.NewEmptyPeer()
}

func NewPeerUnmarshalerFactory(fc model.ModelFactory) model.UnmarshalerFactory {
	return &PeerUnmarshalerFactory{fc}
}

type StorageUnmarshalerFactory struct {
	fc model.ModelFactory
}

func (f *StorageUnmarshalerFactory) CreateUnmarshaler() model.Unmarshaler {
	return f.fc.NewEmptyStorage()
}

func NewStorageUnmarshalerFactory(fc model.ModelFactory) model.UnmarshalerFactory {
	return &StorageUnmarshalerFactory{fc}
}

func (q *QueryProcessor) listObjectQuery(qp model.QueryPayload, wsv core.WSV) (model.QueryResponse, error) {
	address := model.MustAddress(qp.GetFromId())

	list := make([]model.Object, 0)
	switch address.Storage() {
	case core.AccountStorageName:
		res, err := wsv.QueryAll(model.MustAddress(qp.GetFromId()), NewAccountUnmarshalerFactory(q.fc))
		if err != nil {
			return nil, err
		}
		for _, r := range res {
			list = append(list, q.fc.NewObjectBuilder().Account(r.(model.Account)).Build())
		}
	case core.PeerStorageName:
		res, err := wsv.QueryAll(model.MustAddress(qp.GetFromId()), NewPeerUnmarshalerFactory(q.fc))
		if err != nil {
			return nil, err
		}
		for _, r := range res {
			list = append(list, q.fc.NewObjectBuilder().Peer(r.(model.Peer)).Build())
		}
	default:
		res, err := wsv.QueryAll(model.MustAddress(qp.GetFromId()), NewStorageUnmarshalerFactory(q.fc))
		if err != nil {
			return nil, err
		}
		for _, r := range res {
			list = append(list, q.fc.NewObjectBuilder().Storage(r.(model.Storage)).Build())
		}
	}

	qr := q.fc.NewQueryResponseBuilder().
		List(list).
		Build()
	if err := q.signedResponse(qr); err != nil {
		return nil, err
	}
	return qr, nil
}

func (q *QueryProcessor) signedResponse(res model.QueryResponse) error {
	return res.Sign(q.conf.Peer.PublicKeyBytes(), q.conf.Peer.PrivateKeyBytes())
}
