package gate

import (
	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/repository"
)

type APIGate struct {
	queue  core.ProposalTxQueue
	logger log15.Logger
	qp     core.QueryProcessor
}

func NewAPIGate(queue core.ProposalTxQueue, qp core.QueryProcessor, logger log15.Logger) core.APIGate {
	return &APIGate{queue, logger, qp}
}

func (a *APIGate) Write(tx model.Transaction) error {
	if err := tx.Verify(); err != nil {
		return errors.Wrap(core.ErrAPIGateWriteVerifyError, err.Error())
	}
	if err := a.queue.Push(tx); err != nil {
		if errors.Cause(err) == repository.ErrProposalTxQueueAlreadyExistTx {
			return errors.Wrap(core.ErrAPIGateWriteTxAlreadyExist, err.Error())
		}
		return repository.ErrProposalTxQueuePush
	}
	return nil
}

func (a *APIGate) Read(query model.Query) (model.QueryResponse, error) {
	if err := query.Verify(); err != nil {
		return nil, errors.Wrap(core.ErrAPIGateQueryVerifyError, err.Error())
	}
	if err := query.Validate(); err != nil {
		return nil, errors.Wrap(core.ErrAPIGateQueryValidateError, err.Error())
	}
	res, err := a.qp.Query(query)
	if err != nil {
		if errors.Cause(err) == core.ErrQueryProcessorNotFound {
			return nil, errors.Wrap(core.ErrAPIGateQueryNotFound, err.Error())
		}
		return nil, err
	}
	return res, nil
}
