package query

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type QueryProcessor struct {
	rp core.Repository
	fc model.ModelFactory
}

func NewQueryProcessor(rp core.Repository, fc model.ModelFactory) core.QueryProcessor {
	return &QueryProcessor{rp, fc}
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
	err := wsv.Query(qp.GetTargetId(), ac)
	if err != nil {
		return nil, err
	}

	qr := q.fc.NewQueryResponseBuilder().
		Account(ac).
		Build()
	return qr, nil
}

func (q *QueryProcessor) peerObjectQuery(qp model.QueryPayload, wsv core.WSV) (model.QueryResponse, error) {
	peer := q.fc.NewEmptyPeer()
	err := wsv.Query(qp.GetTargetId(), peer)
	if err != nil {
		return nil, err
	}

	qr := q.fc.NewQueryResponseBuilder().
		Peer(peer).
		Build()
	return qr, nil
}
