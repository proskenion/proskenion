package repository

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
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
	queue  map[uint64]model.Transaction
	findTx map[string]uint64
	head   uint64
	tail   uint64
	middle uint64
}

func NewProposalTxQueueOnMemory(conf *config.Config) core.ProposalTxQueue {
	return &ProposalTxQueueOnMemory{
		new(sync.Mutex),
		conf.ProposalTxsLimits,
		make(map[uint64]model.Transaction),
		make(map[string]uint64),
		0,
		0,
		0,
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
	if q.tail-q.head-q.middle < uint64(q.limit) {
		q.findTx[string(hash)] = q.tail
		q.queue[q.tail] = tx
		q.tail++
	} else {
		//log.Print(ErrProposalTxQueueLimits, "queue's max length: %d", q.limit)
		return errors.Wrapf(ErrProposalTxQueueLimits, "queue's max length: %d", q.limit)
	}
	return nil
}

func (q *ProposalTxQueueOnMemory) Pop() (model.Transaction, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for len(q.findTx) != 0 {
		tx, ok := q.queue[q.head]
		if ok {
			txHash, err := tx.Hash()
			if err != nil {
				return nil, false
			}
			delete(q.findTx, string(txHash))
			delete(q.queue, q.head)
			q.head++
			return tx, true
		}
		q.middle--
		q.head++
	}
	return nil, false
}

func (q *ProposalTxQueueOnMemory) Erase(hash model.Hash) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if id, ok := q.findTx[string(hash)]; ok {
		delete(q.queue, id)
		delete(q.findTx, string(hash))
		q.middle++
	} else {
		return errors.Errorf("unexist tx's hash %x", hash)
	}
	return nil
}
