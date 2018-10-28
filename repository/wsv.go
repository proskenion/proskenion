package repository

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/convertor"
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
var OBJECT_ROOT_KEY byte = 5
var PEER_ROOT_KEY byte = 6

func NewWSV(tx core.DBATx, cryptor core.Cryptor, rootHash model.Hash) (core.WSV, error) {
	tree, err := NewMerklePatriciaTree(tx, cryptor, rootHash, WSV_ROOT_KEY)
	if err != nil {
		return nil, err
	}
	return &WSV{
		tx:   tx,
		tree: tree,
		fc:   convertor.NewObjectFactory(cryptor),
	}, nil
}

func (w *WSV) Hash() (model.Hash, error) {
	return w.tree.Hash()
}

// targetId を MerklePatriciaTree の key バイト列に変換
// WIP : @, ., # に対応
func TargetIdToKey(id string) []byte {
	ret := make([]byte, 2)
	ret[0] = WSV_ROOT_KEY
	if id[0] == ':' {
		ret[1] = PEER_ROOT_KEY
	} else {
		ret[1] = OBJECT_ROOT_KEY
	}
	for _, c := range id {
		ret = append(ret, byte(c))
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

func (w *WSV) getPeerRoot() (core.MerklePatriciaNodeIterator, error) {
	if peerRootHash, ok := w.tree.Iterator().Childs()[PEER_ROOT_KEY]; ok {
		it, err := w.tree.Get(peerRootHash)
		if err != nil {
			return nil, err
		}
		fmt.Println("it Key()")
		fmt.Println(it.Key())
		return it, nil
	}
	return nil, errors.Errorf("not found peerservice")
}

// PeerService gets value from targetId
func (w *WSV) PeerService() (core.PeerService, error) {
	peerRoot, err := w.getPeerRoot()
	if err != nil {
		return nil, err
	}
	peerRootHash, err := peerRoot.Hash()
	if err != nil {
		return nil, err
	}
	if len(w.psHash) != 0 {
		// キャッシュがあったら再利用
		if bytes.Equal(w.psHash, peerRootHash) {
			return w.ps, nil
		}
	}

	fmt.Println("peerRoot Key()")
	fmt.Println(peerRoot.Key())
	leafs, err := peerRoot.SubLeafs()
	if err != nil {
		return nil, err
	}
	fmt.Println(len(leafs))
	peers := make([]model.Peer, 0, len(leafs))
	for _, leaf := range leafs {
		peer := w.fc.NewEmptyPeer()
		err := leaf.Data(peer)
		if err != nil {
			return nil, err
		}
		fmt.Println(peer.GetAddress())
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
