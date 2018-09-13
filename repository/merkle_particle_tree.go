package repository

import (
	"bytes"
	"encoding/gob"
	"github.com/pkg/errors"
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
	childs   []model.Hash // child node of tree. (must be alphabet prefix)
	dataHash model.Hash   // data access key
	hash     model.Hash   // Hash of this node
}

type keyMarshaler struct {
	b []byte
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

func (t *MerkleParticleNodeIterator) createMerkleParticleLeafNode(node core.KVNode) (*MerkleParticleNode, error) {
	keyHash, err := t.cryptor.Hash(node.Value())
	err = t.dba.Store(keyHash, node.Value())
	if err != nil {
		return nil, err
	}
	return &MerkleParticleNode{
		childs:   nil,
		dataHash: keyHash,
	}, nil
}

func (k *keyMarshaler) Marshal() ([]byte, error) {
	return k.b, nil
}

func (t *MerkleParticleNodeIterator) Data(unmarshaler core.Unmarshaler) error {
	return t.dba.Load(t.node.dataHash, unmarshaler)
}

// key で参照した先の iterator を取得
func (t *MerkleParticleNodeIterator) Find(key []byte) (core.MerkleParticleNodeIterator, error) {
	if len(key) == 0 {
		return t, nil
	}
	nextKey := t.node.childs[key[0]] // Check : out of range
	newIt := NewMerkleParticleNodeIterator(t.dba, t.cryptor)
	err := t.dba.Load(nextKey, newIt)
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

func (t *MerkleParticleNodeIterator) getLeafNode() (core.MerkleParticleNodeIterator, error) {
	nextIt := NewMerkleParticleNodeIterator(t.dba, t.cryptor)
	err := t.dba.Load(t.node.dataHash, nextIt)
	if err != nil {
		return nil, err
	}
	return nextIt, nil
}

func (t *MerkleParticleNodeIterator) getChildNode(node core.KVNode) (core.MerkleParticleNodeIterator, bool) {
	nextId := NewMerkleParticleNodeIterator(t.dba, t.cryptor)
	if len(t.node.childs[0]) == 0 {
		return nil, false
	}
	err := t.dba.Load(t.node.childs[0], nextId)
	if err != nil {
		return nil, false
	}
	return nextId, true
}

func (t *MerkleParticleNodeIterator) createNewLeafItereator(node core.KVNode) (core.MerkleParticleNodeIterator, error) {
	newLeaf, err := t.createMerkleParticleLeafNode(node)
	if err != nil {
		return nil, err
	}
	newIt := t.createMerkleParticleNodeIterator(newLeaf)
	return newIt, nil
}

func (t *MerkleParticleNodeIterator) newInnerNodeIterator(cnt int, childs ...core.MerkleParticleNodeIterator) (core.MerkleParticleNodeIterator, error) {
	// 子ノードを更新
	var newDataHash model.Hash = nil
	newChilds := make([]model.Hash, 26)
	newKey := t.node.key[:cnt]
	if len(t.Key()) == cnt { // 自身が分裂していない場合は自分の情報を受け継ぐ
		newDataHash = t.DataHash()
		newChilds = t.Childs()
	} else { // 分裂するときは分裂後の子を生成する
		it := t.createMerkleParticleNodeIterator(
			&MerkleParticleNode{
				key:      t.Key()[cnt:],
				childs:   t.Childs(),
				dataHash: t.DataHash(),
			},
		)
		childs = append(childs, it)
	}

	for _, child := range childs {
		hash, err := child.Hash()
		if err != nil {
			return nil, err
		}
		newChilds[child.Key()[0]] = hash
	}

	return t.createMerkleParticleNodeIterator(
		&MerkleParticleNode{
			key:      newKey,
			childs:   newChilds,
			dataHash: newDataHash,
		}), nil
}

// Upsert したあとの Iterator を生成して取得
func (t *MerkleParticleNodeIterator) Upsert(node core.KVNode) (core.MerkleParticleNodeIterator, error) {
	if t.Leaf() {
		return t.Append(node.Value())
	}
	cnt, ok := CountPrefixBytes(t.node.key, node.Key())
	node.Next(cnt)

	// key と prefix が完全一致して且つ子ノードが存在する
	if len(t.node.key) == cnt {
		if it, ok := t.getChildNode(node); ok {
			newIt, err := it.Upsert(node)
			if err != nil {
				return nil, err
			}
			return t.newInnerNodeIterator(cnt, newIt)
		}
	}
	if ok { // Perfect Match
		// key と完全一致したので dataHash の中身を更新
		it, err := t.getLeafNode()
		if err != nil {
			return nil, err
		}
		newIt, err := it.Upsert(node)
		if err != nil {
			return nil, err
		}
		return t.newInnerNodeIterator(cnt, newIt)
	} else if len(node.Key()) == 0 {
		// Insert される Node が -> Leaf Node になる
		// TODO
	} else {
		// Insert される側が Node が InnterNode -> Leaf Node になる
		newIt, err := t.createNewLeafItereator(node)
		if err != nil {
			return nil, err
		}
		return t.newInnerNodeIterator(cnt, newIt)
	}
	return nil, nil
}

// 現在参照しているノードに値を追加
func (t *MerkleParticleNodeIterator) Append(value core.Marshaler) (core.MerkleParticleNodeIterator, error) {
	hash, err := t.cryptor.Hash(value)
	if err != nil {
		return nil, err
	}
	err = t.dba.Store(hash, value)
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
			dataHash: hash,
			prevHash: thisHash,
		},
	)
	newItHash, err := newIt.Hash()
	if err != nil {
		return nil, err
	}
	return newIt, t.dba.Store(newItHash, newIt)
}

func (t *MerkleParticleNodeIterator) Leaf() bool {
	return len(t.node.childs) == 0
}

func (t *MerkleParticleNodeIterator) Prev() core.MerkleParticleNodeIterator {
	it := NewMerkleParticleNodeIterator(t.dba, t.cryptor)
	err := t.dba.Load(t.node.prevHash, it)
	if err != nil {
		return nil
	}
	return it
}

func (t *MerkleParticleNodeIterator) Hash() (model.Hash, error) {
	hash := t.cryptor.ConcatHash(t.node.childs...)
	return t.cryptor.ConcatHash(hash, t.node.dataHash), nil
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

func (t *MerkleParticleNodeIterator) Key() []byte {
	return t.node.key
}

func (t *MerkleParticleNodeIterator) Childs() []model.Hash {
	return t.node.childs
}

func (t *MerkleParticleNodeIterator) DataHash() model.Hash {
	return t.node.dataHash
}
