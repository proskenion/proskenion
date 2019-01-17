package repository

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type WSV struct {
	tx     core.DBATx
	tree   core.MerklePatriciaTree
	fc     model.ObjectFactory
	ps     core.PeerService
	psHash model.Hash
}

var WSV_ROOT_KEY byte = 0

func NewWSV(tx core.DBATx, cryptor core.Cryptor, fc model.ObjectFactory, rootHash model.Hash) (core.WSV, error) {
	tree, err := NewMerklePatriciaTree(tx, cryptor, rootHash, WSV_ROOT_KEY)
	if err != nil {
		return nil, err
	}
	return &WSV{
		tx:   tx,
		tree: tree,
		fc:   fc,
	}, nil
}

func (w *WSV) Hash() model.Hash {
	return w.tree.Hash()
}

// targetId を MerklePatriciaTree の key バイト列に変換
func makeWSVId(id model.Address) []byte {
	ret := make([]byte, 1)
	ret[0] = WSV_ROOT_KEY
	return append(ret, id.GetBytes()...)
}

// Query gets value from targetId
func (w *WSV) Query(targetId model.Address, value model.Unmarshaler) error {
	it, err := w.tree.Find(makeWSVId(targetId))
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

func (w *WSV) QueryAll(fromId model.Address, ufc model.UnmarshalerFactory) ([]model.Unmarshaler, error) {
	it, err := w.tree.Search(makeWSVId(fromId))
	if err != nil {
		if errors.Cause(err) == core.ErrMerklePatriciaTreeNotSearchKey {
			return nil, errors.Wrap(core.ErrWSVNotFound, err.Error())
		}
		return nil, err
	}
	leafs, err := it.SubLeafs()
	rets := make([]model.Unmarshaler, 0, len(leafs))
	for _, leaf := range leafs {
		unm := ufc.CreateUnmarshaler()
		if err = leaf.Data(unm); err != nil {
			return nil, errors.Wrap(core.ErrWSVQueryUnmarshal, err.Error())
		}
		rets = append(rets, unm)
	}
	return rets, nil
}

// PeerService gets value from targetId
func (w *WSV) PeerService(peerRootId model.Address) (core.PeerService, error) {
	peerRoot, err := w.tree.Search(makeWSVId(peerRootId))
	if err != nil {
		return nil, err
	}
	peerRootHash := peerRoot.Hash()
	if len(w.psHash) != 0 {
		// キャッシュがあったら再利用
		if bytes.Equal(w.psHash, peerRootHash) {
			return w.ps, nil
		}
	}

	leafs, err := peerRoot.SubLeafs()
	if err != nil {
		return nil, err
	}

	peers := make([]model.Peer, 0, len(leafs))
	for _, leaf := range leafs {
		peer := w.fc.NewEmptyPeer()
		err := leaf.Data(peer)
		if err != nil {
			return nil, err
		}
		peers = append(peers, peer)
	}
	w.ps = NewPeerService(peers)
	w.psHash = peerRootHash
	return NewPeerService(peers), nil
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
func (w *WSV) Append(targetId model.Address, value model.Marshaler) error {
	_, err := w.tree.Upsert(&KVNode{makeWSVId(targetId), value})
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
