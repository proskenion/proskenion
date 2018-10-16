package commit

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/repository"
	"time"
)

type CommitSystem struct {
	factory  model.ModelFactory
	cryptor  core.Cryptor
	queue    core.ProposalTxQueue
	property *CommitProperty
	rp       core.Repository
}

var (
	ErrCommitLoadPreBlock  = errors.Errorf("Failed Commit Load PreBlock")
	ErrCommitLoadWSV       = errors.Errorf("Failed Commit Load WSV")
	ErrCommitLoadTxHistory = errors.Errorf("Failed Commit Load TxHistory")
)

func NewCommitSystem(factory model.ModelFactory, cryptor core.Cryptor, queue core.ProposalTxQueue, property *CommitProperty, rp core.Repository) core.CommitSystem {
	return &CommitSystem{factory, cryptor, queue, property, rp}
}

func UnixTime(t time.Time) int64 {
	return t.UnixNano()
}

func Now() int64 {
	return UnixTime(time.Now())
}

func rollBackTx(tx core.RepositoryTx, mtErr error) error {
	if err := tx.Rollback(); err != nil {
		return errors.Wrap(err, mtErr.Error())
	}
	return mtErr
}

func commitTx(tx core.RepositoryTx) error {
	if err := tx.Commit(); err != nil {
		return rollBackTx(tx, err)
	}
	return nil
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
	var err error
	wsvHash := model.Hash(nil)
	txHistoryHash := model.Hash(nil)
	topHash := model.Hash(nil)
	topHeight := int64(0)
	if top, ok := c.rp.Top(); ok {
		wsvHash = top.GetPayload().GetWSVHash()
		txHistoryHash = top.GetPayload().GetTxHistoryHash()
		topHash, err = top.Hash()
		topHeight = top.GetPayload().GetHeight()
		if err != nil {
			return nil, nil, err
		}
	}

	dtx, err := c.rp.Begin()
	if err != nil {
		return nil, nil, err
	}

	// load state
	bc, err := dtx.Blockchain(topHash)
	if err != nil {
		return nil, nil, errors.Wrap(ErrCommitLoadPreBlock, err.Error())
	}
	wsv, err := dtx.WSV(wsvHash)
	if err != nil {
		return nil, nil, errors.Wrap(ErrCommitLoadWSV, err.Error())
	}
	txHistory, err := dtx.TxHistory(txHistoryHash)
	if err != nil {
		return nil, nil, errors.Wrap(ErrCommitLoadTxHistory, err.Error())
	}

	txList := repository.NewTxList(c.cryptor)
	// ProposalTxQueue から valid な Tx をとってきて hoge る
	for txList.Size() < c.property.NumTxInBlock {
		tx, ok := c.queue.Pop()
		if !ok {
			break
		}

		// tx を構築
		if err := tx.Validate(wsv, txHistory); err != nil {
			goto txskip
		}
		for _, cmd := range tx.GetPayload().GetCommands() {
			if err := cmd.Validate(wsv); err != nil {
				goto txskip
			}
			if err := cmd.Execute(wsv); err != nil {
				goto txskip // WIP : 要考
				//return nil, nil, rollBackTx(dtx, err)
			}
		}
		if err := txHistory.Append(tx); err != nil {
			return nil, nil, rollBackTx(dtx, err)
		}
		txList.Push(tx)

	txskip:
	}

	newTxHistoryHash, err := txHistory.Hash()
	if err != nil {
		return nil, nil, rollBackTx(dtx, err)
	}
	newWSVHash, err := wsv.Hash()
	if err != nil {
		return nil, nil, rollBackTx(dtx, err)
	}

	newBlock := c.factory.NewBlockBuilder().
		Round(0).
		TxsHash(txList.Top()).
		TxHistoryHash(newTxHistoryHash).
		WSVHash(newWSVHash).
		CreatedTime(Now()).
		Height(topHeight + 1).
		PreBlockHash(topHash).
		Build()
	err = newBlock.Sign(c.property.PublicKey, c.property.PrivateKey)
	if err != nil {
		return nil, nil, rollBackTx(dtx, err)
	}
	if err := bc.Append(newBlock); err != nil {
		return nil, nil, rollBackTx(dtx, err)
	}
	return newBlock, txList, commitTx(dtx)
}
