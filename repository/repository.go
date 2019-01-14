package repository

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

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

type Repository struct {
	dba     core.DBA
	cryptor core.Cryptor
	fc      model.ModelFactory

	top    model.Block
	height int64
}

func NewRepository(dba core.DBA, cryptor core.Cryptor, fc model.ModelFactory) core.Repository {
	return &Repository{dba, cryptor, fc, nil, 0}
}

func (r *Repository) Begin() (core.RepositoryTx, error) {
	tx, err := r.dba.Begin()
	if err != nil {
		return nil, err
	}
	return &RepositoryTx{tx, r.cryptor, r.fc}, nil
}

func (r *Repository) Top() (model.Block, bool) {
	if r.top == nil {
		return nil, false
	}
	return r.top, true
}

func (r *Repository) Commit(block model.Block, txList core.TxList) (err error) {
	dtx, err := r.Begin()
	if err != nil {
		return err
	}

	// load state
	var bc core.Blockchain
	preBlockHash := block.GetPayload().GetPreBlockHash()
	preBlock := r.fc.NewEmptyBlock()
	if !bytes.Equal(preBlockHash, model.Hash(nil)) {
		if bc, err = dtx.Blockchain(preBlockHash); err != nil {
			return errors.Wrap(core.ErrRepositoryCommitLoadPreBlock,
				errors.Errorf("not found hash: %x", block.GetPayload().GetPreBlockHash()).Error())
		}
		preBlock, err = bc.Get(preBlockHash)
		if err != nil {
			return errors.Wrap(core.ErrRepositoryCommitLoadPreBlock,
				errors.Errorf("not found hash: %x", block.GetPayload().GetPreBlockHash()).Error())
		}
	} else {
		if bc, err = dtx.Blockchain(nil); err != nil {
			return err
		}
	}

	wsv, err := dtx.WSV(preBlock.GetPayload().GetWSVHash())
	if err != nil {
		return errors.Wrap(core.ErrRepositoryCommitLoadWSV, err.Error())
	}
	txHistory, err := dtx.TxHistory(preBlock.GetPayload().GetTxHistoryHash())
	if err != nil {
		return errors.Wrap(core.ErrRepositoryCommitLoadTxHistory, err.Error())
	}

	// transactions execute
	for _, tx := range txList.List() {
		if err := tx.Validate(wsv, txHistory); err != nil {
			return rollBackTx(dtx, err)
		}
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
	wsvHash := wsv.Hash()
	if err != nil {
		return rollBackTx(dtx, err)
	}
	if !bytes.Equal(block.GetPayload().GetWSVHash(), wsvHash) {
		return rollBackTx(dtx, errors.Errorf("not equaled wsv Hash : %x", wsvHash))
	}
	txHistoryHash := txHistory.Hash()
	if err != nil {
		return rollBackTx(dtx, err)
	}
	if !bytes.Equal(block.GetPayload().GetTxHistoryHash(), txHistoryHash) {
		return rollBackTx(dtx, errors.Errorf("not equaled txHistory Hash : %x", txHistoryHash))
	}

	// block を追加・
	if err := bc.Append(block); err != nil {
		return err
	}
	// top ブロックを更新
	if r.height < block.GetPayload().GetHeight() {
		r.height = block.GetPayload().GetHeight()
		r.top = block
	}
	return commitTx(dtx)
}

func (r *Repository) GenesisCommit(txList core.TxList) (err error) {
	dtx, err := r.Begin()
	if err != nil {
		return err
	}

	// load state
	var bc core.Blockchain
	if bc, err = dtx.Blockchain(nil); err != nil {
		return err
	}
	wsv, err := dtx.WSV(nil)
	if err != nil {
		return errors.Wrap(core.ErrRepositoryCommitLoadWSV, err.Error())
	}
	txHistory, err := dtx.TxHistory(nil)
	if err != nil {
		return errors.Wrap(core.ErrRepositoryCommitLoadTxHistory, err.Error())
	}

	// transactions execute (no validate)
	for _, tx := range txList.List() {
		for _, cmd := range tx.GetPayload().GetCommands() {
			if err := cmd.Execute(wsv); err != nil {
				return rollBackTx(dtx, err)
			}
		}
		if err := txHistory.Append(tx); err != nil {
			return rollBackTx(dtx, err)
		}
	}

	// hash check and block 生成
	wsvHash := wsv.Hash()
	if err != nil {
		return rollBackTx(dtx, err)
	}
	txHistoryHash := txHistory.Hash()
	if err != nil {
		return rollBackTx(dtx, err)
	}
	genesisBlock := r.fc.NewBlockBuilder().
		CreatedTime(0).
		TxsHash(txList.Top()).
		PreBlockHash(nil).
		TxHistoryHash(txHistoryHash).
		WSVHash(wsvHash).
		Round(0).
		Height(0).
		Build()

	// block を追加・
	if err := bc.Append(genesisBlock); err != nil {
		return err
	}
	// top ブロックを更新
	r.height = genesisBlock.GetPayload().GetHeight()
	r.top = genesisBlock
	return commitTx(dtx)
}

type RepositoryTx struct {
	tx      core.DBATx
	cryptor core.Cryptor
	fc      model.ModelFactory
}

func (r *RepositoryTx) WSV(hash model.Hash) (core.WSV, error) {
	return NewWSV(r.tx, r.cryptor, r.fc, hash)
}

func (r *RepositoryTx) TxHistory(hash model.Hash) (core.TxHistory, error) {
	return NewTxHistory(r.tx, r.fc, r.cryptor, hash)
}

func (r *RepositoryTx) Blockchain(topBlockHash model.Hash) (core.Blockchain, error) {
	return NewBlockchainFromTopBlock(r.tx, r.fc, r.cryptor, topBlockHash)
}

func (r *RepositoryTx) Top() (model.Block, error) {
	return nil, nil
}

func (r *RepositoryTx) Commit() error {
	return r.tx.Commit()
}

func (r *RepositoryTx) Rollback() error {
	return r.tx.Rollback()
}
