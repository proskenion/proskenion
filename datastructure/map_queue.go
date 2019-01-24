package datastructure

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"sync"
)

type ProposalQueueOnMemory struct {
	mutex  *sync.Mutex
	limit  int
	queue  map[uint64]model.Hasher
	find   map[string]uint64
	head   uint64
	tail   uint64
	middle uint64
}

func NewProposalQueueOnMemory(limit int) core.ProposalQueue {
	return &ProposalQueueOnMemory{
		new(sync.Mutex),
		limit,
		make(map[uint64]model.Hasher),
		make(map[string]uint64),
		0,
		0,
		0,
	}
}

func (q *ProposalQueueOnMemory) Push(hasher model.Hasher) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if hasher == nil {
		return errors.Wrapf(core.ErrProposalQueuePushNil, "push object is nil")
	}

	hash := hasher.Hash()
	if _, ok := q.find[string(hash)]; ok {
		return errors.Wrapf(core.ErrProposalQueueAlreadyExist, "already hasher : %x, push to proposal hasher queue", hash)
	}
	if q.tail-q.head-q.middle < uint64(q.limit) {
		q.find[string(hash)] = q.tail
		q.queue[q.tail] = hasher
		q.tail++
	} else {
		//log.Print(ErrProposalQueueLimits, "queue's max length: %d", q.limit)
		return errors.Wrapf(core.ErrProposalQueueLimits, "queue's max length: %d", q.limit)
	}
	return nil
}

func (q *ProposalQueueOnMemory) Pop() (model.Hasher, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for len(q.find) != 0 {
		hasher, ok := q.queue[q.head]
		if ok {
			hasherHash := hasher.Hash()
			delete(q.find, string(hasherHash))
			delete(q.queue, q.head)
			q.head++
			return hasher, true
		}
		q.middle--
		q.head++
	}
	return nil, false
}

func (q *ProposalQueueOnMemory) Erase(hash model.Hash) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if id, ok := q.find[string(hash)]; ok {
		delete(q.queue, id)
		delete(q.find, string(hash))
		q.middle++
	} else {
		return errors.Wrapf(core.ErrProposalQueueEraseUnexist, "unexist hasher's hash %x", hash)
	}
	return nil
}
