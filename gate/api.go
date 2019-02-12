package gate

import (
	"fmt"
	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/repository"
)

type API struct {
	rp     core.Repository
	queue  core.ProposalTxQueue
	logger log15.Logger
	qp     core.QueryProcessor
	qv     core.QueryValidator
	gs     core.Gossip
}

func NewAPI(rp core.Repository, queue core.ProposalTxQueue, qp core.QueryProcessor, qv core.QueryValidator, gs core.Gossip, logger log15.Logger) core.API {
	return &API{rp, queue, logger, qp, qv, gs}
}

func (a *API) Write(tx model.Transaction) error {
	if err := tx.Verify(); err != nil {
		return errors.Wrap(core.ErrAPIWriteVerifyError, err.Error())
	}
	if err := a.queue.Push(tx); err != nil {
		if errors.Cause(err) == core.ErrProposalQueueAlreadyExist {
			return errors.Wrap(core.ErrAPIWriteTxAlreadyExist, err.Error())
		}
		return errors.Wrapf(repository.ErrProposalTxQueuePush, err.Error())
	}
	if err := a.gs.GossipTx(tx); err != nil {
		return errors.Wrap(core.ErrAPIWriteGossipTxError, err.Error())
	}
	return nil
}

func (a *API) Read(query model.Query) (model.QueryResponse, error) {
	if err := query.Verify(); err != nil {
		return nil, errors.Wrap(core.ErrAPIQueryVerifyError, err.Error())
	}
	wsv, err := a.rp.TopWSV()
	if err != nil {
		return nil, fmt.Errorf("Failed APIGate Read, error top WSV: %s", err.Error())
	}
	defer wsv.Commit()
	if err := a.qv.Validate(wsv, query); err != nil {
		return nil, errors.Wrap(core.ErrAPIQueryValidateError, err.Error())
	}
	res, err := a.qp.Query(wsv, query)
	if err != nil {
		if errors.Cause(err) == core.ErrQueryProcessorNotFound {
			return nil, errors.Wrap(core.ErrAPIQueryNotFound, err.Error())
		}
		return nil, err
	}
	return res, nil
}
