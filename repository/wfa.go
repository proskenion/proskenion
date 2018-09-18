package repository

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type WFA struct {
	tx   core.DBATx
	tree core.MerklePatriciaTree
}

var WFA_ROOT_KEY byte = 0

func BeginWFA(dba core.DBA, cryptor core.Cryptor, rootHash model.Hash) (core.WFA, error) {
	tx, err := dba.Begin()
	if err != nil {
		return nil, err
	}
	tree, err := NewMerklePatriciaTree(tx, cryptor, rootHash, WFA_ROOT_KEY)
	if err != nil {
		return nil, err
	}
	return &WFA{tx, tree}, nil
}

func (w *WFA) Hash() (model.Hash, error) {
	return w.tree.Hash()
}

func TargetIdToKey(id string) []byte {
	ret := make([]byte, 1)
	ret[0] = WFA_ROOT_KEY
	for _, c := range id {
		ret = append(ret, byte(c-'a'))
	}
	fmt.Println(ret)
	return ret
}

// Query gets value from targetId
func (w *WFA) Query(targetId string, value model.Unmarshaler) error {
	it, err := w.tree.Find(TargetIdToKey(targetId))
	if err != nil {
		if errors.Cause(err) == core.ErrMerklePatriciaTreeNotFoundKey {
			return errors.Wrap(core.ErrWFANotFound, err.Error())
		}
		return err
	}
	if err = it.Data(value); err != nil {
		return errors.Wrap(core.ErrWFAQueryUnmarshal, err.Error())
	}
	return nil
}

type KVNode struct {
	key       []byte
	marshaler model.Marshaler
}

func (kv *KVNode) Key() []byte {
	return kv.key
}

func (kv *KVNode) Value() model.Marshaler {
	return kv.marshaler
}

func (kv *KVNode) Next(cnt int) core.KVNode {
	return &KVNode{
		kv.key[cnt:],
		kv.marshaler,
	}
}

// Append [targetId] = value
func (w *WFA) Append(targetId string, value model.Marshaler) error {
	_, err := w.tree.Upsert(&KVNode{TargetIdToKey(targetId), value})
	return err
}

// Commit appenging nodes
func (w *WFA) Commit() error {
	return w.tx.Commit()
}

// RollBack
func (w *WFA) Rollback() error {
	return w.tx.Rollback()
}
