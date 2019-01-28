package core

import (
	"github.com/pkg/errors"
	. "github.com/proskenion/proskenion/core/model"
)

const (
	AccountStorageName = "account"
	PeerStorageName    = "peer"
)

var (
	ErrWSVNotFound       = errors.Errorf("Failed WSV Query Not Found")
	ErrWSVQueryUnmarshal = errors.Errorf("Failed WSV Query Unmarshal")

	ErrTxHistoryNotFound       = errors.Errorf("Failed TxHistory Query Not Found")
	ErrTxHistoryQueryUnmarshal = errors.Errorf("Failed TxHistory Query Unmarshal")

	ErrBlockchainNotFound       = errors.Errorf("Failed Blockchain Get Not Found")
	ErrBlockchainQueryUnmarshal = errors.Errorf("Failed Blocchain Get Unmarshal")

	ErrProposalBlockQueuePush = errors.Errorf("Failed ProposalBlockQueue Push")
	ErrProposalTxListCacheSet = errors.Errorf("Failed ProposalTXListCache Set")

	ErrRepositoryCommitLoadPreBlock  = errors.Errorf("Failed Repository Commit Load PreBlockchain")
	ErrRepositoryCommitLoadWSV       = errors.Errorf("Failed Repository Commit Load WSV")
	ErrRepositoryCommitLoadTxHistory = errors.Errorf("Failed Repository Commit Load TxHistory")
)

// TxList Wrap MerkleTree
type TxList interface {
	Push(tx Transaction) error
	List() []Transaction
	Size() int
	Modelor
}

type TxListCache interface {
	Set(txList TxList) error
	Get(hash Hash) (TxList, bool)
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
	// GetTxList gets txList from txHash
	GetTxList(txListHash Hash) (TxList, error)
	// GetTxList gets
	GetTx(txHash Hash) (Transaction, error)
	// Append tx
	Append(txList TxList) error
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

type ProposalBlockQueue interface {
	Push(block Block) error
	Erase(hash Hash) error
	Pop() (Block, bool)
	WaitPush() struct{}
}

type Repository interface {
	Begin() (RepositoryTx, error)
	Top() (Block, bool)
	GetDelegatedAccounts() ([]Account, error)
	Commit(Block, TxList) error
	GenesisCommit(TxList) error
	CreateBlock(queue ProposalTxQueue, round int32, now int64) (Block, TxList, error)
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
