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
	factory model.ModelFactory
	cryptor core.Cryptor

	height int64
	top    model.Block
}

var (
	ErrCommitLoadPreBlock  = errors.Errorf("Failed Commit Load PreBlock")
	ErrCommitLoadWSV       = errors.Errorf("Failed Commit Load WSV")
	ErrCommitLoadTxHistory = errors.Errorf("Failed Commit Load TxHistory")
)

func NewCommitSystem(dba core.DBA, factory model.ModelFactory, cryptor core.Cryptor) core.CommitSystem {
	return &CommitSystem{dba, factory, cryptor, 0, nil}
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
	dtx, err := c.dba.Begin()
	if err != nil {
		return err
	}

	// load state
	bc := repository.NewBlockchain(dtx, c.factory)
	preBlock, ok := bc.Get(block.GetPayload().GetPreBlockHash())
	if !ok {
		return errors.Wrap(ErrCommitLoadPreBlock,
			errors.Errorf("not found hash: %x", block.GetPayload().GetPreBlockHash()).Error())
	}

	wsv, err := repository.NewWSV(dtx, c.cryptor, preBlock.GetPayload().GetWSVHash())
	if err != nil {
		return errors.Wrap(ErrCommitLoadWSV, err.Error())
	}
	txHistory, err := repository.NewTxHistory(dtx, c.factory, c.cryptor, preBlock.GetPayload().GetTxHistoryHash())
	if err != nil {
		return errors.Wrap(ErrCommitLoadTxHistory, err.Error())
	}

	// transactions execute
	for _, tx := range txList.List() {
		for _, cmd := range tx.GetPayload().GetCommands() {
			if err := cmd.Validate(wsv); err != nil {
				return rollBackTx(dtx, err)
			}
			if err := cmd.Execute(wsv); err != nil {
				return rollBackTx(dtx, err)
			}
		}
		if err := txHistory.Append(tx); err != nil {
			return rollBackTx(dtx, err)
		}
	}

	// hash check
	wsvHash, err := wsv.Hash()
	if err != nil {
		return rollBackTx(dtx, err)
	}
	if !bytes.Equal(block.GetPayload().GetWSVHash(), wsvHash) {
		return rollBackTx(dtx, errors.Errorf("not equaled wsv Hash : %x", wsvHash))
	}
	txHistoryHash, err := txHistory.Hash()
	if err != nil {
		return rollBackTx(dtx, err)
	}
	if !bytes.Equal(block.GetPayload().GetTxHistoryHash(), txHistoryHash) {
		return rollBackTx(dtx, errors.Errorf("not equaled txHistory Hash : %x", txHistoryHash))
	}

	// top ブロックを更新
	if c.height < block.GetPayload().GetHeight() {
		c.height = block.GetPayload().GetHeight()
		c.top = block
	}
	return commitTx(dtx)
}

// CreateBlock
func (c *CommitSystem) CreateBlock() (model.Block, core.TxList, error) {
	var err error
	wsvHash := model.Hash(nil)
	txHistoryHash := model.Hash(nil)
	if c.top != nil {
		wsvHash = c.top.GetPayload().GetWSVHash()
		txHistoryHash = c.top.GetPayload().GetTxHistoryHash()
	}

	dtx, err := c.dba.Begin()
	if err != nil {
		return nil, nil, err
	}

	// load state
	wsv, err := repository.NewWSV(dtx, c.cryptor, wsvHash)
	if err != nil {
		return nil, nil, errors.Wrap(ErrCommitLoadWSV, err.Error())
	}
	txHistory, err := repository.NewTxHistory(dtx, c.factory, c.cryptor, txHistoryHash)
	if err != nil {
		return nil, nil, errors.Wrap(ErrCommitLoadTxHistory, err.Error())
	}

	// ProposalTxQueue から valid な Tx をとってきて hoge る

}
