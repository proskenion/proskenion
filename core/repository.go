package core

import (
	"github.com/pkg/errors"
	. "github.com/proskenion/proskenion/core/model"
)

var (
	ErrMerklePatriciaTreeNotFoundKey = errors.Errorf("Failed MerklePatriciaTree Not Found key")
	ErrInvalidKVNodes                = errors.Errorf("Failed Key Value Nodes Invalid")
	ErrWSVNotFound                   = errors.Errorf("Failed WSV Query Not Found")
	ErrWSVQueryUnmarshal             = errors.Errorf("Failed WSV Query Unmarshal")
	ErrTxHistoryNotFound             = errors.Errorf("Failed WSV Query Not Found")
	ErrTxHistoryQueryUnmarshal       = errors.Errorf("Failed WSV Query Unmarshal")
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

// Merkle Patricia Tree に対する操作
type MerklePatriciaController interface {
	// key で参照した先の iterator を取得
	Find(key []byte) (MerklePatriciaNodeIterator, error)
	// Upsert したあとの Iterator を生成して取得
	Upsert(KVNode) (MerklePatriciaNodeIterator, error)
	// hash を Root にする
	Set(hash Hash) error
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
}

// WSV (MerklePatriciaTree で管理)
type WSV interface {
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

// 全Tx履歴 (MerklePatriciaTree で管理)
type TxHistory interface {
	Hash() (Hash, error)
	// Query gets tx from txHash
	Query(txHash Hash) (Transaction, error)
	// Append tx
	Append(tx Transaction) error
	// Commit appenging nodes
	Commit() error
	// RollBack
	Rollback() error
}

// BlockChain
type Blockchain interface {
	Top() (Block, bool)
	// Commit block
	Commit(block Block) error
	VerifyCommit(block Block) error
}
