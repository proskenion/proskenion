package repository

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"log"
	"sync"
)

var (
	ErrProposalTxQueueLimits         = errors.Errorf("PropposalTxQueue run limit reached")
	ErrProposalTxQueueAlreadyExistTx = errors.Errorf("Failed Push Already Exist Tx")
	ErrProposalTxQueuePush           = errors.Errorf("Failed ProposalTxQueue Push")
)

type ProposalTxQueueOnMemory struct {
	mutex  *sync.Mutex
	limit  int
	queue  []model.Transaction
	findTx map[string]model.Transaction
}

func NewProposalTxQueueOnMemory(conf *config.Config) core.ProposalTxQueue {
	return &ProposalTxQueueOnMemory{
		new(sync.Mutex),
		conf.ProposalTxsLimits,
		make([]model.Transaction, 0, conf.ProposalTxsLimits),
		make(map[string]model.Transaction),
	}
}

func (q *ProposalTxQueueOnMemory) Push(tx model.Transaction) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if tx == nil {
		return errors.Wrapf(model.ErrInvalidTransaction, "push transaction is nil")
	}

	hash, err := tx.Hash()
	if err != nil {
		return errors.Wrapf(model.ErrTransactionHash, err.Error())
	}
	if _, ok := q.findTx[string(hash)]; ok {
		return errors.Wrapf(ErrProposalTxQueueAlreadyExistTx, "already tx : %x, push to proposal tx queue", hash)
	}
	if len(q.queue) < q.limit {
		q.findTx[string(hash)] = tx
		q.queue = append(q.queue, tx)
	} else {
		log.Print(ErrProposalTxQueueLimits, "queue's max length: ", q.limit)
		return errors.Wrapf(ErrProposalTxQueueLimits, "queue's max length: %d", q.limit)
	}
	return nil
}

func (q *ProposalTxQueueOnMemory) Pop() (model.Transaction, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.queue) == 0 {
		return nil, false
	}
	front := q.queue[0]
	frontHash, err := front.Hash()
	if err != nil {
		return nil, false
	}
	delete(q.findTx, string(frontHash))
	q.queue = q.queue[1:]
	return front, true
}
