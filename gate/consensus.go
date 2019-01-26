package gate

import (
	"bytes"
	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/repository"
	"io"
)

type ConsensusGate struct {
	fc          model.ModelFactory
	c           core.Cryptor
	txQueue     core.ProposalTxQueue
	txListCache core.TxListCache
	blockQueue  core.ProposalBlockQueue
	logger      log15.Logger
	conf        *config.Config
}

func NewConsensusGate(fc model.ModelFactory, c core.Cryptor, txQueue core.ProposalTxQueue, txListCache core.TxListCache, blockQueue core.ProposalBlockQueue, logger log15.Logger, conf *config.Config) core.ConsensusGate {
	return &ConsensusGate{fc, c, txQueue, txListCache, blockQueue, logger, conf}
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
		return errors.Wrapf(core.ErrProposalBlockQueuePush, err.Error())
	}
	return nil
}

func (c *ConsensusGate) PropagateBlockAck(block model.Block) (model.Signature, error) {
	if err := block.Verify(); err != nil {
		return nil, errors.Wrap(core.ErrConsensusGatePropagateBlockVerifyError, err.Error())
	}
	// TODO err ceheck : public key は存在する Peer のものか。
	signature, err := c.c.Sign(block, c.conf.Peer.PrivateKeyBytes())
	if err != nil {
		return nil, err
	}
	ret := c.fc.NewSignature(c.conf.Peer.PublicKeyBytes(), signature)
	return ret, nil
}

func (c *ConsensusGate) PropagateBlockStreamTx(block model.Block, txChan chan model.Transaction, errChan chan error) error {
	txList := repository.NewTxList(c.c)
	for {
		select {
		case tx := <-txChan:
			txList.Push(tx)
		case err := <-errChan:
			if err != nil && err != io.EOF {
				return err
			}
			goto afterFor
		}
	}
afterFor:

	if !bytes.Equal(block.GetPayload().GetTxsHash(), txList.Hash()) {
		return errors.Wrapf(core.ErrConsensusGatePropagateBlockDifferentHash,
			"txsHash: %x, txListHash: %x",
			block.GetPayload().GetTxsHash(), txList.Hash())
	}

	if err := c.txListCache.Set(txList); err != nil {
		return errors.Wrapf(core.ErrProposalTxListCacheSet, err.Error())
	}
	if err := c.blockQueue.Push(block); err != nil {
		return errors.Wrapf(core.ErrProposalBlockQueuePush, err.Error())
	}
	return nil
}

// chan が Stream の返り値
func (c *ConsensusGate) CollectTx(blockHash model.Hash, txChan chan model.Transaction, errChan chan error) error {
	for _, tx := range []model.Transaction{} { // TODO
		txChan <- tx
		if err := <-errChan; err != nil {
			return err
		}
	}
	return nil
}
