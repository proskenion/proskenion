package repository

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type TxHistory struct {
	tx      core.DBATx
	factory model.ModelFactory
	tree    core.MerklePatriciaTree
}

var TX_HISTORY_ROOT_KEY byte = 1

func NewTxHistory(tx core.DBATx, factory model.ModelFactory, cryptor core.Cryptor, rootHash model.Hash) (core.TxHistory, error) {
	tree, err := NewMerklePatriciaTree(tx, cryptor, rootHash, TX_HISTORY_ROOT_KEY)
	if err != nil {
		return nil, err
	}
	return &TxHistory{tx, factory, tree}, nil
}

func (w *TxHistory) Hash() model.Hash {
	return w.tree.Hash()
}

func TxHashToKey(txHash model.Hash) []byte {
	return append([]byte{TX_HISTORY_ROOT_KEY}, txHash...)
}

// Query gets value from targetId
func (w *TxHistory) Query(txHash model.Hash) (model.Transaction, error) {
	txHash = TxHashToKey(txHash)
	it, err := w.tree.Find(txHash)
	if err != nil {
		if errors.Cause(err) == core.ErrMerklePatriciaTreeNotFoundKey {
			return nil, errors.Wrap(core.ErrTxHistoryNotFound, err.Error())
		}
		return nil, err
	}
	retTx := w.factory.NewEmptyTx()
	if err = it.Data(retTx); err != nil {
		return nil, errors.Wrap(core.ErrTxHistoryQueryUnmarshal, err.Error())
	}
	return retTx, nil
}

// Append [targetId] = value
func (w *TxHistory) Append(tx model.Transaction) error {
	txHash := tx.Hash()
	txHash = TxHashToKey(txHash)
	_, err := w.tree.Upsert(&KVNode{txHash, tx})
	return err
}

// Commit appenging nodes
func (w *TxHistory) Commit() error {
	if err := w.tx.Commit(); err != nil {
		if err := w.Rollback(); err != nil {
			return err
		}
		return err
	}
	return nil
}

// RollBack
func (w *TxHistory) Rollback() error {
	return w.tx.Rollback()
}
