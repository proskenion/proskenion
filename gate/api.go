package gate

import (
	"fmt"
	"github.com/inconshreveable/log15"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/repository"
)

var (
	ErrAPIGateWriteTxNil = fmt.Errorf("Failed APIGate Write Error Tx nil")
)

type APIGate struct {
	db     core.DB
	queue  core.ProposalTxQueue
	logger log15.Logger
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
	return nil, nil
}
