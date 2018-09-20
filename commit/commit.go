package commit

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/repository"
)

type CommitSystem struct {
	dba     core.DBA
	bc      core.Blockchain
	factory model.ModelFactory
	cryptor core.Cryptor
}

var (
	ErrCommitLoadPreBlock  = errors.Errorf("Failed Commit Load PreBlock")
	ErrCommitLoadWSV       = errors.Errorf("Failed Commit Load WSV")
	ErrCommitLoadTxHistory = errors.Errorf("Failed Commit Load TxHistory")
)

func NewCommitSystem(bc core.Blockchain) core.Commit {
	return &CommitSystem{bc}
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
	preBlock, ok := c.bc.Get(block.GetPayload().GetPreBlockHash())
	if !ok {
		return errors.Wrap(ErrCommitLoadPreBlock,
			errors.Errorf("not found hash: %x", block.GetPayload().GetPreBlockHash()).Error())
	}

	dtx, err := c.dba.Begin()
	if err != nil {
		return err
	}
	wsv, err := repository.NewWSV(dtx, c.cryptor, preBlock.GetPayload().GetWSVHash())
	if err != nil {
		return errors.Wrap(ErrCommitLoadWSV, err.Error())
	}
	txHistory, err := repository.NewTxHistory(dtx, c.factory, c.cryptor, preBlock.GetPayload().GetTxHistoryHash())
	if err != nil {
		return errors.Wrap(ErrCommitLoadTxHistory, err.Error())
	}
	for _, tx := range txList.List() {
		for _, cmd := range tx.GetPayload().GetCommands() {
			cmd.GetTransfer().Execute()
		}
	}
	return nil
}
