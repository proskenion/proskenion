package gate

import (
	"fmt"
	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/repository"
)

var (
	ErrAPIGateWriteTxNil         = fmt.Errorf("Failed APIGate Write Error Tx nil")
	ErrAPIGateQueryTxNil         = fmt.Errorf("Failed APIGate Read Error Query nil")
	ErrAPIGateQueryVerifyError   = fmt.Errorf("Failed APIGate Read query Verify Error")
	ErrAPIGateQueryValitateError = fmt.Errorf("Failed APIGate Read query Validate Error")
)

type APIGate struct {
	queue  core.ProposalTxQueue
	logger log15.Logger
	qp     core.QueryProcessor
}

func NewAPIGate(queue core.ProposalTxQueue, logger log15.Logger, qp core.QueryProcessor) core.APIGate {
	return &APIGate{queue, logger, qp}
}

func (a *APIGate) Write(tx model.Transaction) error {
	if tx == nil {
		return ErrAPIGateWriteTxNil
	}
	if err := a.queue.Push(tx); err != nil {
		if err == repository.ErrProposalTxQueueAlreadyExistTx {
			a.logger.Debug("Write Tx but already exists : ", err.Error())
			return nil
		}
		return repository.ErrProposalTxQueuePush
	}
	return nil
}

func (a *APIGate) Read(query model.Query) (model.QueryResponse, error) {
	if query == nil {
		return nil, ErrAPIGateQueryTxNil
	}
	if err := query.Verify(); err != nil {
		return nil, errors.Wrap(ErrAPIGateQueryVerifyError, err.Error())
	}
	if err := query.Validate(); err != nil {
		return nil, errors.Wrap(ErrAPIGateQueryValitateError, err.Error())
	}
	res, err := a.qp.Query(query)
	if err != nil {
		return nil, err
	}
	return res, nil
}
