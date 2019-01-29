package core

import (
	"github.com/pkg/errors"
	. "github.com/proskenion/proskenion/core/model"
)

var (
	ErrMerklePatriciaTreeNotSearchKey = errors.Errorf("Failed MerklePatriciaTree can not search key")
	ErrMerklePatriciaTreeNotFoundKey  = errors.Errorf("Failed MerklePatriciaTree Not Found key")
	ErrInvalidKVNodes                 = errors.Errorf("Failed Key Value Nodes Invalid")
)

// ProposalQueue
var (
	ErrProposalQueueLimits       = errors.Errorf("PropposalQueue run limit reached")
	ErrProposalQueueAlreadyExist = errors.Errorf("Failed Push Already Exist ")
	ErrProposalQueuePush         = errors.Errorf("Failed ProposalQueue Push")
	ErrProposalQueuePushNil      = errors.Errorf("Failed ProposalQueue Push nil hasher")
	ErrProposalQueueEraseUnexist = errors.Errorf("Faield Erase Unexist Hash")
)

// CacheMap
var (
	ErrCacheMapPush    = errors.Errorf("Failed CacheMap Push")
	ErrCacheMapPushNil = errors.Errorf("Failed CacheMap Push nil hasher")
)

type KeyValueStore interface {
	Load(key Hash, value Unmarshaler) error // value = Load(key)
	Store(key Hash, value Marshaler) error  // Duplicate Insert error
}

// ProposalQueue (Cached Queue + erase by hash)
type ProposalQueue interface {
	Push(hasher Hasher) error
	Erase(hash Hash) error
	Pop() (Hasher, bool)
}

// CacheMap
type CacheMap interface {
	Set(hasher Hasher) error
	Get(hash Hash) (Hasher, bool)
}

// Transaction 列の管理
type MerkleTree interface {
	Push(hash Hasher) error
	Hasher
}

type KVNode interface {
	// KVNode{key = Key()[cnt:], value=value}
	Next(cnt int) KVNode
	Key() []byte
	Value() Marshaler
}

// Merkle Patricia Tree に対する操作
type MerklePatriciaController interface {
	// key と prefix が一致している最も浅い internal iterator を取得
	Search(key []byte) (MerklePatriciaNodeIterator, error)
	// key で参照した先の leaf iterator を取得
	Find(key []byte) (MerklePatriciaNodeIterator, error)
	// Upsert したあとの Iterator を生成して取得
	Upsert(KVNode) (MerklePatriciaNodeIterator, error)
	// hash を Root にする
	Set(hash Hash) error
	Get(hash Hash) (MerklePatriciaNodeIterator, error)
	Hasher
	Marshaler
	Unmarshaler
}

// World State の管理 に使う(SubTree の管理にも使う)
type MerklePatriciaTree interface {
	Iterator() MerklePatriciaNodeIterator
	MerklePatriciaController
}

// Merkle Patricia Node を管理する Iterator
type MerklePatriciaNodeIterator interface {
	MerklePatriciaController
	Key() []byte
	Childs() map[byte]Hash
	DataHash() Hash
	Leaf() bool
	Data(unmarshaler Unmarshaler) error
	Prev() (MerklePatriciaNodeIterator, error)
	SubLeafs() ([]MerklePatriciaNodeIterator, error)
}
