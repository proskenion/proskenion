package gate

import (
	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/repository"
)

type ConsensusGate struct {
	cc     core.CommitSystem
	queue  core.ProposalTxQueue
	logger log15.Logger
}

func NewConsensusGate(cc core.CommitSystem, queue core.ProposalTxQueue, logger log15.Logger) core.ConsensusGate {
	return &ConsensusGate{cc, queue, logger}
}

func (c *ConsensusGate) PropagateTx(tx model.Transaction) error {
	if err := tx.Verify(); err != nil {
		return errors.Wrap(core.ErrConsensusGatePropagateTxVerifyError, err.Error())
	}
	if err := c.queue.Push(tx); err != nil {
		if errors.Cause(err) == repository.ErrProposalQueuePush {
			return nil // already exist no error.
		}
		return errors.Wrapf(repository.ErrProposalTxQueuePush, err.Error())
	}
	return nil
}

func (c *ConsensusGate) PropagateBlock(block model.Block) error {
	return nil
}

// chan が Stream の返り値
func (c *ConsensusGate) CollectTx(blockHash model.Hash, txChan chan model.Transaction) error {
	return nil
}
