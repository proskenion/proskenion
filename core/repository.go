package core

import (
	"github.com/pkg/errors"
	. "github.com/proskenion/proskenion/core/model"
)

var (
	ErrMerkleParticleTreeNotFoundKey = errors.Errorf("Failed MerkleParticleTree Not Found key")
	ErrInvalidKVNodes                = errors.Errorf("Failed Key Value Nodes Invalid")
)

// Transaction 列の管理
type MerkleTree interface {
	Push(hash Hasher) error
	Top() Hash
}

type KVNode interface {
	// KVNode{key = Key()[cnt:], value=value}
	Next(cnt int) KVNode
	Key() []byte
	Value() Marshaler
}

// Merkle Particle Tree に対する操作
type MerkleParticleController interface {
	// key で参照した先の iterator を取得
	Find(key []byte) (MerkleParticleNodeIterator, error)
	// Upsert したあとの Iterator を生成して取得
	Upsert(KVNode) (MerkleParticleNodeIterator, error)
	Hasher
	Marshaler
	Unmarshaler
}

var MERKLE_PARTICLE_CHILD_EDGES = 26

// World State の管理 に使う(SubTree の管理にも使う)
type MerkleParticleTree interface {
	Iterator() MerkleParticleNodeIterator
	MerkleParticleController
}

// Merkle Particle Node を管理する Iterator
type MerkleParticleNodeIterator interface {
	MerkleParticleController
	Key() []byte
	Childs() []Hash
	DataHash() Hash
	Leaf() bool
}

// WFA
type WFA interface {
	Hash() (Hash, error)
	// Query gets value from targetId
	Query(targetId string, value Unmarshaler) error
	// Append [targetId] = value
	Append(targetId string, value Marshaler) error
	// Commit appenging nodes
	Commit() error
	// RollBack
	Rollback() error
}

// 全Tx履歴
type TxHistory interface {
	FindTx(hash Hash) (Transaction, error)
}

// BlockChain
type Blockchain interface {
	Top() (Block, bool)
	// Commit is allowed only Commitable Block, ohterwise panic
	Commit(block Block)
	VerifyCommit(block Block) error
}
