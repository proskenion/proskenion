package core

import (
	"github.com/pkg/errors"
	. "github.com/proskenion/proskenion/core/model"
)

var (
	ErrMerklePatriciaTreeNotSearchKey = errors.Errorf("Failed MerklePatriciaTree can not search key")
	ErrMerklePatriciaTreeNotFoundKey  = errors.Errorf("Failed MerklePatriciaTree Not Found key")
	ErrInvalidKVNodes                 = errors.Errorf("Failed Key Value Nodes Invalid")
	ErrWSVNotFound                    = errors.Errorf("Failed WSV Query Not Found")
	ErrWSVQueryUnmarshal              = errors.Errorf("Failed WSV Query Unmarshal")
	ErrTxHistoryNotFound              = errors.Errorf("Failed TxHistory Query Not Found")
	ErrTxHistoryQueryUnmarshal        = errors.Errorf("Failed TxHistory Query Unmarshal")
	ErrBlockchainNotFound             = errors.Errorf("Failed Blockchain Get Not Found")
	ErrBlockchainQueryUnmarshal       = errors.Errorf("Failed Blocchain Get Unmarshal")
)

var (
	ErrRepositoryCommitLoadPreBlock  = errors.Errorf("Failed Repository Commit Load PreBlockchain")
	ErrRepositoryCommitLoadWSV       = errors.Errorf("Failed Repository Commit Load WSV")
	ErrRepositoryCommitLoadTxHistory = errors.Errorf("Failed Repository Commit Load TxHistory")
)

// Transaction 列の管理
type MerkleTree interface {
	Push(hash Hasher) error
	Top() Hash
}

// TxList Wrap MerkleTree
type TxList interface {
	Push(tx Transaction) error
	Top() Hash
	List() []Transaction
	Size() int
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

// WSV (MerklePatriciaTree で管理)
type WSV interface {
	Hasher
	// Query gets value from targetId
	Query(targetId Address, value Unmarshaler) error
	// Query All gets value from fromId
	QueryAll(fromId Address, value UnmarshalerFactory) ([]Unmarshaler, error)
	// Get PeerService
	PeerService(peerRootId Address) (PeerService, error)
	// Append [targetId] = value
	Append(targetId Address, value Marshaler) error
	// Commit appenging nodes
	Commit() error
	// RollBack
	Rollback() error
}

// 全Tx履歴 (MerklePatriciaTree で管理)
type TxHistory interface {
	Hasher
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
	// blockHash を指定して Block を取得
	Get(blockHash Hash) (Block, error)
	// Commit block
	Append(block Block) error
}

// 提案された Transaction を保持する Queue
type ProposalTxQueue interface {
	Push(tx Transaction) error
	Erase(hash Hash) error
	Pop() (Transaction, bool)
}

type Repository interface {
	Begin() (RepositoryTx, error)
	Top() (Block, bool)
	Commit(Block, TxList) error
	GenesisCommit(TxList) error
}

type RepositoryTx interface {
	WSV(Hash) (WSV, error)
	TxHistory(Hash) (TxHistory, error)
	Blockchain(Hash) (Blockchain, error)
	Commit() error
	Rollback() error
}

// Peer 取得機構
type PeerService interface {
	List() []Peer
}
