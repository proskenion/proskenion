package repository

import (
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/datastructure"
)

type ProposalBlockQueueOnMemory struct {
	core.ProposalQueue
	pushChan chan struct{}
}

func NewProposalBlockQueueOnMemory(conf *config.Config) core.ProposalBlockQueue {
	return &ProposalBlockQueueOnMemory{datastructure.NewProposalQueueOnMemory(conf.Queue.BlockLimits),
		make(chan struct{}, conf.Queue.BlockLimits)}
}

func (q *ProposalBlockQueueOnMemory) Push(block model.Block) error {
	if err := q.ProposalQueue.Push(block); err != nil {
		return err
	}
	q.pushChan <- struct{}{}
	return nil
}

func (q *ProposalBlockQueueOnMemory) Pop() (model.Block, bool) {
	ret, ok := q.ProposalQueue.Pop()
	if !ok {
		return nil, false
	}
	block, ok := ret.(model.Block)
	if !ok {
		return nil, false
	}
	return block, true
}

func (q *ProposalBlockQueueOnMemory) Erase(hash model.Hash) error {
	return q.ProposalQueue.Erase(hash)
}

func (q *ProposalBlockQueueOnMemory) WaitPush() struct{} {
	return <-q.pushChan
}
