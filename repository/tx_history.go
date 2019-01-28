package repository

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/datastructure"
)

type TxHistory struct {
	tx          core.DBATx
	factory     model.ModelFactory
	c           core.Cryptor
	tree        core.MerklePatriciaTree
	txListCache core.TxListCache
}

type TxIndexed struct {
	TxListHash model.Hash
	Index      int
}

func (t *TxIndexed) Marshal() ([]byte, error) {
	return model.GobMarshal(t)
}

func (t *TxIndexed) Unmarshal(b []byte) error {
	return model.GobUnmarshal(b, t)
}

var (
	TxHistoryRootKey       byte = 1
	TxHistoryTxRootKey     byte = 0
	TxHistoryTxListRootKey byte = 1
)

func NewTxHistory(tx core.DBATx, factory model.ModelFactory, cryptor core.Cryptor, cache core.TxListCache, rootHash model.Hash) (core.TxHistory, error) {
	tree, err := datastructure.NewMerklePatriciaTree(tx, cryptor, rootHash, TxHistoryRootKey)
	if err != nil {
		return nil, err
	}
	return &TxHistory{tx, factory, cryptor, tree, cache}, nil
}

func (w *TxHistory) Hash() model.Hash {
	return w.tree.Hash()
}

func TxHashToKey(txHash model.Hash) []byte {
	return append([]byte{TxHistoryRootKey, TxHistoryTxRootKey}, txHash...)
}

func TxListHashToKey(txHash model.Hash) []byte {
	return append([]byte{TxHistoryRootKey, TxHistoryTxListRootKey}, txHash...)
}

// GetTxList gets txList from txHash
func (w *TxHistory) GetTxList(txListHash model.Hash) (core.TxList, error) {
	// get txList from cache.
	if cachetxList, ok := w.txListCache.Get(txListHash); ok {
		return cachetxList, nil
	}

	txListKey := TxListHashToKey(txListHash)
	it, err := w.tree.Find(txListKey)
	if err != nil {
		if errors.Cause(err) == core.ErrMerklePatriciaTreeNotFoundKey {
			return nil, errors.Wrap(core.ErrTxHistoryNotFound, err.Error())
		}
		return nil, err
	}
	retTxList := NewTxList(w.c,w.factory)
	if err = it.Data(retTxList); err != nil {
		return nil, errors.Wrap(core.ErrTxHistoryQueryUnmarshal, err.Error())
	}

	// set cache txList.
	if err := w.txListCache.Set(retTxList); err != nil {
		return nil, err
	}
	return retTxList, nil
}

// GetTxList gets
func (w *TxHistory) GetTx(txHash model.Hash) (model.Transaction, error) {
	txKey := TxHashToKey(txHash)
	it, err := w.tree.Find(txKey)
	if err != nil {
		if errors.Cause(err) == core.ErrMerklePatriciaTreeNotFoundKey {
			return nil, errors.Wrap(core.ErrTxHistoryNotFound, err.Error())
		}
		return nil, err
	}
	retTxIndexed := &TxIndexed{}
	if err = it.Data(retTxIndexed); err != nil {
		return nil, errors.Wrap(core.ErrTxHistoryQueryUnmarshal, err.Error())
	}

	txList, err := w.GetTxList(retTxIndexed.TxListHash)
	if err != nil {
		return nil, err
	}
	if len(txList.List()) >= retTxIndexed.Index {
		return nil,
			fmt.Errorf("Failed GetTx index out of range. len is %d, but index %d.", len(txList.List()), retTxIndexed.Index)
	}
	return txList.List()[retTxIndexed.Index], nil
}

// Append txList
func (w *TxHistory) Append(txList core.TxList) error {
	txListKey := TxListHashToKey(txList.Hash())
	if _, err := w.tree.Upsert(&KVNode{txListKey, txList}); err != nil {
		return err
	}
	for i, tx := range txList.List() {
		txKey := TxHashToKey(tx.Hash())
		if _, err := w.tree.Upsert(&KVNode{txKey, &TxIndexed{txList.Hash(), i}}); err != nil {
			return err
		}
	}
	return nil
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
