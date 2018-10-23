package query

import (
	"bytes"
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

func containsPublicKey(keys []model.PublicKey, pub model.PublicKey) bool {
	for _, key := range keys {
		if bytes.Equal(key, pub) {
			return true
		}
	}
	return false
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

	// 署名チェック
	ac := q.fc.NewEmptyAccount()
	err = wsv.Query(query.GetPayload().GetAuthorizerId(), ac)
	if err != nil {
		return nil, errors.Wrapf(core.ErrQueryProcessorNotExistAuthoirizer,
			"authorizer : %s", query.GetPayload().GetAuthorizerId())
	}
	if !containsPublicKey(ac.GetPublicKeys(), query.GetSignature().GetPublicKey()) {
		return nil, errors.Wrapf(core.ErrQueryProcessorNotSignedAuthorizer,
			"authorizer : %s, expect key : %x",
			query.GetPayload().GetAuthorizerId(), query.GetSignature().GetPublicKey())
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
	if err := q.signedResponse(qr); err != nil {
		return nil, err
	}
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
	if err := q.signedResponse(qr); err != nil {
		return nil, err
	}
	return qr, nil
}

func (q *QueryProcessor) signedResponse(res model.QueryResponse) error {
	return res.Sign(q.conf.Peer.PublicKeyBytes(), q.conf.Peer.PrivateKeyBytes())
}
