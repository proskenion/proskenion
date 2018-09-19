package repository

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type Blockchain struct {
	dba        core.DBA
	factory    model.ModelFactory
	top        model.Block
	mostHeight int64
}

func NewBlockchain(dba core.DBA, factory model.ModelFactory) core.Blockchain {
	return &Blockchain{dba, factory, nil, 0}
}

func (b *Blockchain) Top() (model.Block, bool) {
	if b.top == nil {
		return nil, false
	}
	return b.top, true
}

func rollBackTx(tx core.DBATx, mtErr error) error {
	if err := tx.Rollback(); err != nil {
		return errors.Wrap(err, mtErr.Error())
	}
	return mtErr
}

func commitTx(tx core.DBATx) error {
	if err := tx.Commit(); err != nil {
		return rollBackTx(tx, err)
	}
	return nil
}

// Commit is allowed only Commitable Block, ohterwise panic
func (b *Blockchain) Commit(block model.Block) (err error) {
	tx, err := b.dba.Begin()
	if err != nil {
		return err
	}
	hash, err := block.Hash()
	if err != nil {
		return rollBackTx(tx, err)
	}
	if err = tx.Store(hash, block); err != nil {
		return rollBackTx(tx, err)
	}
	return commitTx(tx)
}

// Commit 可能かどうかの判定
func (b *Blockchain) VerifyCommit(block model.Block) error {
	if err := block.Verify(); err != nil {
		return err
	}
	// preBlockHash と同値の状態が存在するかの判定
	preBlock := b.factory.NewEmptyBlock()
	err := b.dba.Load(block.GetPayload().GetPreBlockHash(), preBlock)
	preBlockHash, err := preBlock.Hash()
	if err != nil {
		return err
	}
	if !bytes.Equal(preBlockHash, block.GetPayload().GetPreBlockHash()) {
		return errors.Errorf("Failed Blockchain Verify Commit Not Matched Hash")
	}
	return nil
}
