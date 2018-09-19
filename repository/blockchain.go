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
	wfa        WFA
	txHistory  core.TxHistory
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

// Commit is allowed only Commitable Block, ohterwise panic
func (b *Blockchain) Commit(block model.Block) error {
	return nil
}

// Commit 可能かどうかの判定
func (b *Blockchain) VerifyCommit(block model.Block) error {
	if err := block.Verify(); err != nil {
		return err
	}
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
