package repository

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type WSV struct {
	tx   core.DBATx
	tree core.MerklePatriciaTree
}

var WSV_ROOT_KEY byte = 0

func NewWSV(tx core.DBATx, cryptor core.Cryptor, rootHash model.Hash) (core.WSV, error) {
	tree, err := NewMerklePatriciaTree(tx, cryptor, rootHash, WSV_ROOT_KEY)
	if err != nil {
		return nil, err
	}
	return &WSV{tx, tree}, nil
}

func (w *WSV) Hash() (model.Hash, error) {
	return w.tree.Hash()
}

// targetId を MerklePatriciaTree の key バイト列に変換
// WIP : @, ., # に対応
func TargetIdToKey(id string) []byte {
	ret := make([]byte, 1)
	ret[0] = WSV_ROOT_KEY
	for _, c := range id {
		ret = append(ret, byte(c-'a'))
	}
	return ret
}

// Query gets value from targetId
func (w *WSV) Query(targetId string, value model.Unmarshaler) error {
	it, err := w.tree.Find(TargetIdToKey(targetId))
	if err != nil {
		if errors.Cause(err) == core.ErrMerklePatriciaTreeNotFoundKey {
			return errors.Wrap(core.ErrWSVNotFound, err.Error())
		}
		return err
	}
	if err = it.Data(value); err != nil {
		return errors.Wrap(core.ErrWSVQueryUnmarshal, err.Error())
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
func (w *WSV) Append(targetId string, value model.Marshaler) error {
	_, err := w.tree.Upsert(&KVNode{TargetIdToKey(targetId), value})
	return err
}

// Commit appenging nodes
func (w *WSV) Commit() error {
	if err := w.tx.Commit(); err != nil {
		if err := w.Rollback(); err != nil {
			return err
		}
		return err
	}
	return nil
}

// RollBack
func (w *WSV) Rollback() error {
	return w.tx.Rollback()
}
