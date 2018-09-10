package core

import (
	"github.com/pkg/errors"
	. "github.com/proskenion/proskenion/core/model"
)

var (
	ErrMerkleParticleTreeNotFoundKey = errors.Errorf("Failed MerkleParticleTree Not Found key")
)

// Transaction 列の管理
type MerkleTree interface {
	Push(hash Hasher) error
	Top() Hash
}

// World State の管理
type MerkleParticleTree interface {
	Root() MerkleParticleIterator
	Find(key Marshaler, value Unmarshaler) error
	Upsert(key Marshaler, value Unmarshaler) error
	Hash() (Hash, error)
	Unmarshal() ([]byte, error)
}

type MerkleParticleIterator interface {
	Data(unmarshaler Unmarshaler) error
	Iterator() MerkleParticleNodeIterator
	Find(key Marshaler, value Unmarshaler) error
	Upsert(key Marshaler, value Unmarshaler) error
	Next() MerkleParticleIterator
	Prev() MerkleParticleIterator
	First() bool
	Last() bool
	Hash() (Hash, error)
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

type MerkleParticleNodeIterator interface {
	Data(unmarshaler Unmarshaler) error
	Find(key Marshaler, value Unmarshaler) error
	Upsert(key Marshaler, value Unmarshaler) error
	Next() MerkleParticleNodeIterator
	Prev() MerkleParticleNodeIterator
	First() bool
	Last() bool
	Hash() (Hash, error)
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

// WFA
type WFA interface {
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
