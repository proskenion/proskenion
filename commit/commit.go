package commit

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"time"
)

type CommitSystem struct {
	factory  model.ModelFactory
	cryptor  core.Cryptor
	queue    core.ProposalTxQueue
	property *CommitProperty
	rp       core.Repository
}

func NewCommitSystem(factory model.ModelFactory, cryptor core.Cryptor, queue core.ProposalTxQueue, property *CommitProperty, rp core.Repository) core.CommitSystem {
	return &CommitSystem{factory, cryptor, queue, property, rp}
}

func UnixTime(t time.Time) int64 {
	return t.UnixNano()
}

func Now() int64 {
	return UnixTime(time.Now())
}

// Stateless Validate
func (c *CommitSystem) VerifyCommit(block model.Block, txList core.TxList) error {
	if err := block.Verify(); err != nil {
		return err
	}
	if !bytes.Equal(block.GetPayload().GetTxsHash(), txList.Top()) {
		return errors.Errorf("Failed Verify Commit: not matched txsHash")
	}
	for _, tx := range txList.List() {
		if err := tx.Verify(); err != nil {
			return err
		}
	}
	return nil
}

func (c *CommitSystem) Commit(block model.Block, txList core.TxList) error {
	return c.rp.Commit(block, txList)
}

// CreateBlock
func (c *CommitSystem) CreateBlock() (model.Block, core.TxList, error) {
	return c.rp.CreateBlock(c.queue, Now())
}
