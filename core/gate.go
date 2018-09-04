package core

import . "github.com/proskenion/proskenion/core/model"

type APIGate interface {
	Tx(tx Transaction) error
	Query(query Query) (QueryResponse, error)
}

type ConsensusGate interface {
	PropagateTx(tx Transaction) error
	PropagateBlock(block Block) error
	// chan が Stream の返り値
	CollectTx(blockHash Hash, txChan chan Transaction) error
}

type SyncGate interface {
	Sync(blockHash Hash, blockChan chan Block) error
}
