package core

import . "github.com/proskenion/proskenion/core/model"

// Transaction 列の管理
type MerkleTree interface {
	Push(hash Hasher) error
	Top() Hash
}

// World State の管理
type MerkleParticleTree interface {
	// TODO あとで考える
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
