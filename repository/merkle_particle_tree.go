package repository

import (
	"bytes"
	"encoding/gob"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

// World State の管理
type MerkleParticleTree struct {
	dba     core.KeyValueStore
	cryptor core.Cryptor
	root    core.MerkleParticleNodeIterator
}

func NewMerkleParticleTree(kvStore core.KeyValueStore, cryptor core.Cryptor, hash model.Hash) core.MerkleParticleTree {
	return &MerkleParticleTree{
		dba:     kvStore,
		cryptor: cryptor,
	}
}

func (t *MerkleParticleTree) Iterator() core.MerkleParticleNodeIterator {
	return t.root
}

// key で参照した先の iterator を取得
func (t *MerkleParticleTree) Find(key []byte) (core.MerkleParticleNodeIterator, error) {
	return t.Iterator().Find(key)
}

// Upsert したあとの新しい Iterator を生成して取得
func (t *MerkleParticleTree) Upsert(node core.KVNode) (core.MerkleParticleNodeIterator, error) {
	return t.Iterator().Upsert(node)
}

// 現在参照しているノードに値を追加
func (t *MerkleParticleTree) Append(value core.Marshaler) (core.MerkleParticleNodeIterator, error) {
	return t.Iterator().Append(value)
}

func (t *MerkleParticleTree) Hash() (model.Hash, error) {
	return t.Iterator().Hash()
}

func (t *MerkleParticleTree) Marshal() ([]byte, error) {
	return t.Iterator().Marshal()
}

func (t *MerkleParticleTree) Unmarshal(b []byte) error {
	return t.Iterator().Unmarshal(b)
}

type MerkleParticleNode struct {
	key      []byte
	height   uint64       // height of merkle Node, Merkle Particle Node like blockChain
	childs   []model.Hash // child node of tree. (must be alphabet prefix)
	dataKey  model.Hash   // data access key
	prevHash model.Hash   // previous node (like blockChain)
}

type keyMarshaler struct {
	b []byte
}

func createMarshaler(b []byte) core.Marshaler {
	return &keyMarshaler{b}
}

type MerkleParticleNodeIterator struct {
	dba     core.KeyValueStore
	cryptor core.Cryptor
	node    *MerkleParticleNode
}

func NewMerkleParticleNodeIterator(dba core.KeyValueStore, cryptor core.Cryptor) core.MerkleParticleNodeIterator {
	return &MerkleParticleNodeIterator{
		dba:     dba,
		cryptor: cryptor,
	}
}

func (t *MerkleParticleNodeIterator) createMerkleParticleNodeIterator(node *MerkleParticleNode) core.MerkleParticleNodeIterator {
	return &MerkleParticleNodeIterator{
		dba:     t.dba,
		cryptor: t.cryptor,
		node:    node,
	}
}

func (t *MerkleParticleNodeIterator) createMerkleParticleLeafNode(kvNode core.KVNode) (*MerkleParticleNode, error) {
	keyHash, err := t.cryptor.Hash(kvNode.Value())
	err = t.dba.Store(createMarshaler(keyHash), kvNode.Value())
	if err != nil {
		return nil, err
	}
	return &MerkleParticleNode{
		height:   0,
		childs:   nil,
		dataKey:  keyHash,
		prevHash: nil,
	}, nil
}

func (k *keyMarshaler) Marshal() ([]byte, error) {
	return k.b, nil
}

func (t *MerkleParticleNodeIterator) Data(unmarshaler core.Unmarshaler) error {
	return t.dba.Load(createMarshaler(t.node.dataKey), unmarshaler)
}

// key で参照した先の iterator を取得
func (t *MerkleParticleNodeIterator) Find(key []byte) (core.MerkleParticleNodeIterator, error) {
	if len(key) == 0 {
		return t, nil
	}
	nextKey := t.node.childs[key[0]] // Check : out of range
	newIt := NewMerkleParticleNodeIterator(t.dba, t.cryptor)
	err := t.dba.Load(createMarshaler(nextKey), newIt)
	if err != nil {
		return nil, err
	}
	return newIt.Find(key[1:])
}

// return ( number of match prefix bytes, Is perfect matche? )
func CountPrefixBytes(a []byte, b []byte) (int, bool) {
	cnt := 0
	for ; cnt < len(a) && cnt < len(b); cnt++ {
		if a[cnt] != b[cnt] {
			return cnt, false
		}
	}
	if len(a) == len(b) {
		return cnt, true
	}
	return cnt, false
}

func (t *MerkleParticleNodeIterator) getNode(hash model.Hash) (core.MerkleParticleNodeIterator, error) {
	nextIt := NewMerkleParticleNodeIterator(t.dba, t.cryptor)
	err := t.dba.Load(createMarshaler(hash), nextIt)
	if err != nil {
		return nil, err
	}
	return nextIt, nil
}

// Upsert したあとの Iterator を生成して取得
func (t *MerkleParticleNodeIterator) Upsert(kvNode core.KVNode) (core.MerkleParticleNodeIterator, error) {
	cnt, ok := CountPrefixBytes(t.node.key, kvNode.Key())
	if ok { // Perfect Match
		// key と完全一致したので dataKey の中身を更新
		it, err := t.getNode(t.node.dataKey)
		if err != nil {
			return nil, err
		}
		it, err = it.Append(kvNode.Value())
		if err != nil {
			return nil, err
		}
		return it, nil
	} else if cnt == len(t.node.key) {
		// node の Key と一致。 kvNode から生える。
	} else if cnt == len(kvNode.Key()) {
		// new node の Key と一致。 new node から生える。
	} else {
		// 両方 prefix を分割
		baseKey := t.node.key[:cnt]
		newIt := t.createMerkleParticleNodeIterator(
			&MerkleParticleNode{
				key:      baseKey,
				height:   0,
				childs:   make([]model.Hash, 26),
				dataKey:  nil,
				prevHash: nil,
			},
		)
		newItL := t.createMerkleParticleNodeIterator(
			&MerkleParticleNode{
				key:      t.node.key[cnt:],
				height:   t.node.height,
				childs:   t.node.childs,
				dataKey:  t.node.dataKey,
				prevHash: t.node.prevHash,
			},
		)
		newLeaf, err := t.createMerkleParticleLeafNode(kvNode)
		if err != nil {
			return nil, err
		}
		newItR := t.createMerkleParticleNodeIterator(newLeaf)
		var itL core.MerkleParticleNodeIterator
		if cnt < len(t.node.key) {
			createMerkleParticleNodeIterator(t.dba, t.cryptor,
				&MerkleParticleNode{
					height:   0,
					childs:   make([]model.Hash, 16),
					dataKey:  keyMarshal,
					prevHash: thisHash,
				})
		}
		var itR core.MerkleParticleNodeIterator
		if cnt < len(kvNode.Key()) {

		}

	}
	return nil, nil
}

// 現在参照しているノードに値を追加
func (t *MerkleParticleNodeIterator) Append(value core.Marshaler) (core.MerkleParticleNodeIterator, error) {
	hash, err := t.cryptor.Hash(value)
	if err != nil {
		return nil, err
	}
	keyMarshal := createMarshaler(hash)
	err = t.dba.Store(keyMarshal, value)
	if err != nil {
		return nil, err
	}
	thisHash, err := t.Hash()
	if err != nil {
		return nil, err
	}
	newIt := t.createMerkleParticleNodeIterator(
		&MerkleParticleNode{
			depth:    t.node.depth,
			height:   t.node.height + 1,
			childs:   make([]model.Hash, 16),
			dataKey:  keyMarshal,
			prevHash: thisHash,
		},
	)
	newItHash, err := newIt.Hash()
	if err != nil {
		return nil, err
	}
	return newIt, t.dba.Store(createMarshaler(newItHash), newIt)
}

func (t *MerkleParticleNodeIterator) Prev() core.MerkleParticleNodeIterator {
	it := NewMerkleParticleNodeIterator(t.dba, t.cryptor)
	err := t.dba.Load(createMarshaler(t.node.prevHash), it)
	if err != nil {
		return nil
	}
	return it
}

func (t *MerkleParticleNodeIterator) Hash() (model.Hash, error) {
	return t.cryptor.Hash(t)
}

func (t *MerkleParticleNodeIterator) Marshal() ([]byte, error) {
	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.
	err := enc.Encode(t.node)
	if err != nil {
		return nil, err
	}
	return network.Bytes(), nil
}

func (t *MerkleParticleNodeIterator) Unmarshal(b []byte) error {
	network := bytes.NewBuffer(b)
	dec := gob.NewDecoder(network) // Will read from network.
	err := dec.Decode(t.node)
	if err != nil {
		return err
	}
	return nil
}
