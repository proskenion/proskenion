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
func (t *MerkleParticleTree) Upsert(node []core.KVNode) (core.MerkleParticleNodeIterator, error) {
	return t.Iterator().Upsert(node)
}

// 現在参照しているノードに値を追加
func (t *MerkleParticleTree) Append(value core.Marshaler) error {
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
	depth    uint32 // depth of tree, root tree is 0
	height   uint64 // height of merkle Node, Merkle Particle Node like blockChain
	childs   []model.Hash
	dataKey  core.Marshaler
	prevHash model.Hash
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

func (k *keyMarshaler) Marshal() ([]byte, error) {
	return k.b, nil
}

func (t *MerkleParticleNodeIterator) Data(unmarshaler core.Unmarshaler) error {
	return t.dba.Load(t.node.dataKey, unmarshaler)
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

// Upsert したあとの Iterator を生成して取得
func (t *MerkleParticleNodeIterator) Upsert(kvNodes []core.KVNode) (core.MerkleParticleNodeIterator, error) {
	// TODO
	return nil, nil
}

// 現在参照しているノードに値を追加
func (t *MerkleParticleNodeIterator) Append(value core.Marshaler) error {
	hash, err := t.cryptor.Hash(value)
	if err != nil {
		return err
	}
	keyMarshal := createMarshaler(hash)
	err = t.dba.Store(keyMarshal, value)
	if err != nil {
		return err
	}
	thisHash, err := t.Hash()
	if err != nil {
		return err
	}
	newIt := &MerkleParticleNodeIterator{
		node: &MerkleParticleNode{
			depth:    t.node.depth,
			height:   t.node.height + 1,
			childs:   make([]model.Hash, 256),
			dataKey:  keyMarshal,
			prevHash: thisHash,
		},
	}
	newItHash, err := newIt.Hash()
	if err != nil {
		return err
	}
	return t.dba.Store(createMarshaler(newItHash), newIt)
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
