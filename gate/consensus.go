package gate

import (
	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/repository"
)

type ConsensusGate struct {
	cc         core.CommitSystem
	txQueue    core.ProposalTxQueue
	blockQueue core.ProposalBlockQueue
	logger     log15.Logger
}

func NewConsensusGate(cc core.CommitSystem, txQueue core.ProposalTxQueue, blockQueue core.ProposalBlockQueue, logger log15.Logger) core.ConsensusGate {
	return &ConsensusGate{cc, txQueue, blockQueue, logger}
}

func (c *ConsensusGate) PropagateTx(tx model.Transaction) error {
	if err := tx.Verify(); err != nil {
		return errors.Wrap(core.ErrConsensusGatePropagateTxVerifyError, err.Error())
	}
	if err := c.txQueue.Push(tx); err != nil {
		if errors.Cause(err) == core.ErrProposalQueueAlreadyExist {
			return nil // already exist no error.
		}
		return errors.Wrapf(repository.ErrProposalTxQueuePush, err.Error())
	}
	return nil
}

func (c *ConsensusGate) PropagateBlock(block model.Block) error {
	if err := block.Verify(); err != nil {
		return errors.Wrap(core.ErrConsensusGatePropagateBlockVerifyError, err.Error())
	}
	if err := c.blockQueue.Push(block); err != nil {
		if errors.Cause(err) == core.ErrProposalQueueAlreadyExist {
			return errors.Wrapf(core.ErrConsensusGatePropagateBlockAlreadyExist, err.Error())
		}
		return errors.Wrapf(repository.ErrProposalBlockQueuePush, err.Error())
	}
	return nil
}

// chan が Stream の返り値
func (c *ConsensusGate) CollectTx(blockHash model.Hash, txChan chan model.Transaction, errChan chan error) error {
	for _, tx := range []model.Transaction{} { // TODO
		txChan <- tx
		if err := <- errChan; err != nil {
			return err
		}
	}
	return nil
}
