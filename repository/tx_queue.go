package repository

import (
	"fmt"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

var (
	ErrProposalTxQueuePush = fmt.Errorf("Failed ProposalTxQueue Push")
)

type ProposalTxQueueOnMemory struct {
	core.ProposalQueue
}

func NewProposalTxQueueOnMemory(conf *config.Config) core.ProposalTxQueue {
	return &ProposalTxQueueOnMemory{NewProposalQueueOnMemory(conf.Queue.TxsLimits)}
}

func (q *ProposalTxQueueOnMemory) Push(tx model.Transaction) error {
	return q.ProposalQueue.Push(tx)
}

func (q *ProposalTxQueueOnMemory) Pop() (model.Transaction, bool) {
	ret, ok := q.ProposalQueue.Pop()
	if !ok {
		return nil, false
	}
	tx, ok := ret.(model.Transaction)
	if !ok {
		return nil, false
	}
	return tx, true
}

func (q *ProposalTxQueueOnMemory) Erase(hash model.Hash) error {
	return q.ProposalQueue.Erase(hash)
}
