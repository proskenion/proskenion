package repository

import (
	"bytes"
	"encoding/gob"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"strconv"
)

// World State の管理
type MerkleParticleTree struct {
	dba     core.KeyValueStore
	cryptor core.Cryptor
	root    core.MerkleParticleNodeIterator
}

func NewMerkleParticleTree(kvStore core.KeyValueStore, cryptor core.Cryptor, hash model.Hash, rootKey byte) (core.MerkleParticleTree, error) {
	newInternal := &MerkleParticleNodeIterator{
		dba:     kvStore,
		cryptor: cryptor,
		node:    &MerkleParticleInternalNode{},
	}
	err := kvStore.Load(hash, newInternal)
	if err != nil {
		if err != core.ErrDBANotFoundLoad {
			return nil, err
		}

		// ROOT Internal Noed
		newInternal = &MerkleParticleNodeIterator{
			dba:     kvStore,
			cryptor: cryptor,
			node: &MerkleParticleInternalNode{
				Key_:      []byte{rootKey},
				Childs_:   make([]model.Hash, core.MERKLE_PARTICLE_CHILD_EDGES),
				DataHash_: model.Hash(nil),
			},
		}
		hash, err := newInternal.Hash()
		if err != nil {
			return nil, err
		}
		// saved
		err = kvStore.Store(hash, newInternal)
		if err != nil {
			return nil, err
		}
	}

	return &MerkleParticleTree{
		dba:     kvStore,
		cryptor: cryptor,
		root:    newInternal,
	}, nil
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

func (t *MerkleParticleTree) Hash() (model.Hash, error) {
	return t.Iterator().Hash()
}

func (t *MerkleParticleTree) Marshal() ([]byte, error) {
	return t.Iterator().Marshal()
}

func (t *MerkleParticleTree) Unmarshal(b []byte) error {
	return t.Iterator().Unmarshal(b)
}

type MerkleParticleNode interface {
	Leaf() bool

	// Internal
	Key() []byte
	Childs() []model.Hash
	DataHash() model.Hash

	// Leaf
	Height() int64
	PrevHash() model.Hash
	DataObject() []byte
}

type MerkleParticleInternalNode struct {
	Key_      []byte
	Childs_   []model.Hash // child node of tree. (must be alphabet prefix)
	DataHash_ model.Hash   // data access key (must be leaf node)
}

func (n *MerkleParticleInternalNode) Leaf() bool {
	return false
}

func (n *MerkleParticleInternalNode) Key() []byte {
	return n.Key_
}

func (n *MerkleParticleInternalNode) Childs() []model.Hash {
	return n.Childs_
}

func (n *MerkleParticleInternalNode) DataHash() model.Hash {
	return n.DataHash_
}

func (n *MerkleParticleInternalNode) Height() int64 {
	return 0
}

func (n *MerkleParticleInternalNode) PrevHash() model.Hash {
	return model.Hash{nil}
}

func (n *MerkleParticleInternalNode) DataObject() []byte {
	return nil
}

type MerkleParticleLeafNode struct {
	Height_     int64
	PrevHash_   model.Hash
	DataObject_ []byte // Unmarshaled data object
}

func (n *MerkleParticleLeafNode) Leaf() bool {
	return false
}

func (n *MerkleParticleLeafNode) Key() []byte {
	return nil
}

func (n *MerkleParticleLeafNode) Childs() []model.Hash {
	return nil
}

func (n *MerkleParticleLeafNode) DataHash() model.Hash {
	return nil
}

func (n *MerkleParticleLeafNode) Height() int64 {
	return n.Height_
}

func (n *MerkleParticleLeafNode) PrevHash() model.Hash {
	return n.PrevHash_
}

func (n *MerkleParticleLeafNode) DataObject() []byte {
	return n.DataObject_
}

type MerkleParticleNodeIterator struct {
	dba     core.KeyValueStore
	cryptor core.Cryptor
	node    MerkleParticleNode
}

func NewMerkleParticleNodeIterator(dba core.KeyValueStore, cryptor core.Cryptor) core.MerkleParticleNodeIterator {
	return &MerkleParticleNodeIterator{
		dba:     dba,
		cryptor: cryptor,
	}
}

// new** は単に型の生成、データの保存は行わない
func (t *MerkleParticleNodeIterator) newEmptyLeafIterator() core.MerkleParticleNodeIterator {
	return &MerkleParticleNodeIterator{
		dba:     t.dba,
		cryptor: t.cryptor,
		node:    &MerkleParticleLeafNode{},
	}
}

func (t *MerkleParticleNodeIterator) newEmptyInternalIterator() core.MerkleParticleNodeIterator {
	return &MerkleParticleNodeIterator{
		dba:     t.dba,
		cryptor: t.cryptor,
		node:    &MerkleParticleInternalNode{},
	}
}

func (t *MerkleParticleNodeIterator) getChild(key byte) (core.MerkleParticleNodeIterator, error) {
	unmarshaler := t.newEmptyInternalIterator()
	if key < 0 || len(t.Childs()) <= int(key) {
		return nil, errors.Errorf("Childs Out of Range")
	}
	err := t.dba.Load(t.Childs()[key], unmarshaler)
	if err != nil {
		return nil, err
	}
	return unmarshaler, nil
}

func (t *MerkleParticleNodeIterator) getLeaf() (core.MerkleParticleNodeIterator, error) {
	unmarshaler := t.newEmptyLeafIterator()
	err := t.dba.Load(t.DataHash(), unmarshaler)
	if err != nil {
		return nil, err
	}
	return unmarshaler, nil
}

func (t *MerkleParticleNodeIterator) getObject(unmarshaler core.Unmarshaler) error {
	return unmarshaler.Unmarshal(t.node.DataObject())
}

// create ** はデータを保存する
func (t *MerkleParticleNodeIterator) createMerkleParticleNodeIterator(node MerkleParticleNode) (core.MerkleParticleNodeIterator, error) {
	it := &MerkleParticleNodeIterator{
		dba:     t.dba,
		cryptor: t.cryptor,
		node:    node,
	}
	hash, err := it.Hash()
	if err != nil {
		return nil, err
	}
	// saved
	err = t.dba.Store(hash, it)
	if err != nil {
		return nil, err
	}
	return it, nil
}

// 追加したい node 情報から新しい葉ノード(or node にまだ key が残っているなら interanl node)のイテレータを保存して返す
func (t *MerkleParticleNodeIterator) createLeafIterator(node core.KVNode) (core.MerkleParticleNodeIterator, error) {
	object_, err := node.Value().Marshal()
	if err != nil {
		return nil, err
	}
	newLeafIt, err := t.createMerkleParticleNodeIterator(&MerkleParticleLeafNode{
		Height_:     0,
		DataObject_: object_,
		PrevHash_:   model.Hash(nil),
	})
	if err != nil {
		return nil, err
	}

	if len(node.Key()) > 0 {
		hash, err := newLeafIt.Hash()
		if err != nil {
			return nil, err
		}
		newInternalIt, err := t.createMerkleParticleNodeIterator(
			&MerkleParticleInternalNode{
				Key_:      node.Key(),
				DataHash_: hash,
				Childs_:   make([]model.Hash, core.MERKLE_PARTICLE_CHILD_EDGES),
			},
		)
		return newInternalIt, nil
	}
	return newLeafIt, nil
}

// key で参照した先の iterator を取得
func (t *MerkleParticleNodeIterator) Find(key []byte) (core.MerkleParticleNodeIterator, error) {
	if t.Leaf() {
		return t, nil
	}
	if len(t.Key()) < len(key) {
		return nil, core.ErrMerkleParticleTreeNotFoundKey
	}
	nextChild, err := t.getChild(key[0])
	if err != nil {
		return nil, err
	}
	return nextChild.Find(key[len(t.Key()):])
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

// node の Key に沿って子を返す
func (t *MerkleParticleNodeIterator) getChildFromNode(node core.KVNode) (core.MerkleParticleNodeIterator, bool) {
	it, err := t.getChild(node.Key()[0])
	if err != nil {
		return nil, false
	}
	return it, true
}

// ノード t を cnt 番目で分割、新たに child を加えた時の 中間ノードの生成(場合により child にも変更を加える)
func (t *MerkleParticleNodeIterator) createInternalIterator(cnt int, child core.MerkleParticleNodeIterator) (core.MerkleParticleNodeIterator, error) {
	// 子ノードを更新
	var newDataHash model.Hash = nil
	newChilds := make([]model.Hash, core.MERKLE_PARTICLE_CHILD_EDGES)
	newKey := t.Key()[:cnt]
	childs := []core.MerkleParticleNodeIterator{child}
	if len(t.Key()) == cnt { // 自身が分裂していない場合は自分の情報を受け継ぐ
		newDataHash = t.DataHash()
		newChilds = t.Childs()
	} else { // 分裂するときは分裂後の子を生成する
		// 分岐後の自分(child)
		it, err := t.createMerkleParticleNodeIterator(
			&MerkleParticleInternalNode{
				Key_:      t.Key()[cnt:],
				Childs_:   t.Childs(),
				DataHash_: t.DataHash(),
			},
		)
		if err != nil {
			return nil, err
		}

		if child.Leaf() {
			// 追加する子が葉であるとき、葉は DataHash に、自分のみを子にする
			newDataHash, err = child.Hash()
			if err != nil {
				return nil, err
			}
			childs = []core.MerkleParticleNodeIterator{it}
		} else {
			// そうでなければ自分の分身を新しい子の集合に加える
			childs = append(childs, it)
		}
	}

	for _, child := range childs {
		hash, err := child.Hash()
		if err != nil {
			return nil, err
		}
		newChilds[child.Key()[0]] = hash
	}

	return t.createMerkleParticleNodeIterator(
		&MerkleParticleInternalNode{
			Key_:      newKey,
			Childs_:   newChilds,
			DataHash_: newDataHash,
		})
}

// Upsert したあとの Iterator を生成して取得
func (t *MerkleParticleNodeIterator) Upsert(node core.KVNode) (core.MerkleParticleNodeIterator, error) {
	if t.Leaf() {
		return t.Append(node.Value())
	}
	cnt, ok := CountPrefixBytes(t.Key(), node.Key())
	node.Next(cnt)

	// key と prefix が完全一致して且つ子ノードが存在する
	if len(t.Key()) == cnt {
		if it, ok := t.getChildFromNode(node); ok {
			newIt, err := it.Upsert(node)
			if err != nil {
				return nil, err
			}
			return t.createInternalIterator(cnt, newIt)
		}
	}
	if ok { // Perfect Match
		// key と完全一致したので dataHash の中身を更新
		it, err := t.getLeaf()
		if err != nil {
			return nil, err
		}
		newIt, err := it.Upsert(node)
		if err != nil {
			return nil, err
		}
		return t.createInternalIterator(cnt, newIt)
	} else {
		// 現在のノードを分割して中身を更新
		newIt, err := t.createLeafIterator(node)
		if err != nil {
			return nil, err
		}
		return t.createInternalIterator(cnt, newIt)
	}
	return nil, nil
}

// 現在参照しているノードに値を追加
func (t *MerkleParticleNodeIterator) Append(value core.Marshaler) (core.MerkleParticleNodeIterator, error) {
	hash, err := t.cryptor.Hash(value)
	if err != nil {
		return nil, err
	}
	thisHash, err := t.Hash()
	if err != nil {
		return nil, err
	}
	newIt, err := t.createMerkleParticleNodeIterator(
		&MerkleParticleLeafNode{
			Height_:     t.node.Height() + 1,
			DataObject_: hash,
			PrevHash_:   thisHash,
		},
	)
	if err != nil {
		return nil, err
	}
	return newIt, nil
}

func (t *MerkleParticleNodeIterator) Leaf() bool {
	return t.node.Leaf()
}

func (t *MerkleParticleNodeIterator) Prev() core.MerkleParticleNodeIterator {
	it := t.newEmptyLeafIterator()
	err := t.dba.Load(t.node.PrevHash(), it)
	if err != nil {
		return nil
	}
	return it
}

func (t *MerkleParticleNodeIterator) Hash() (model.Hash, error) {
	if t.Leaf() {
		return t.cryptor.ConcatHash(t.node.PrevHash(),
			t.node.DataObject(),
			[]byte(strconv.FormatInt(t.node.Height(), 10))), nil
	} else {
		hash := t.cryptor.ConcatHash(t.node.Childs()...)
		return t.cryptor.ConcatHash(hash, t.node.DataHash(), t.Key()), nil
	}
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
	return t.node.Key()
}

func (t *MerkleParticleNodeIterator) Childs() []model.Hash {
	return t.node.Childs()
}

func (t *MerkleParticleNodeIterator) DataHash() model.Hash {
	return t.node.DataHash()
}
