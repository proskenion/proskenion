package repository

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type Blockchain struct {
	tx      core.DBATx
	factory model.ModelFactory
}

func NewBlockchain(tx core.DBATx, factory model.ModelFactory) core.Blockchain {
	return &Blockchain{tx, factory}
}

func (b *Blockchain) Get(blockHash model.Hash) (model.Block, bool) {
	retBlock := b.factory.NewEmptyBlock()
	err := b.tx.Load(blockHash, retBlock)
	if err != nil {
		return nil, false
	}
	return retBlock, true
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
func (b *Blockchain) Append(block model.Block) (err error) {
	hash, err := block.Hash()
	return b.tx.Store(hash, block)
}
