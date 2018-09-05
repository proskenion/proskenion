package core

import . "github.com/proskenion/proskenion/core/model"

type APIGateRequest interface {
	Tx(tx Transaction) error
	Query(query Query) (QueryResponse, error)
}

type ConsensusGateRquest interface {
	PropagateTx(tx Transaction) error
	PropagateBlock(block Block) error
	// chan が Stream の返り値
	CollectTx(blockHash Hash, txChan chan Transaction) error
}

type SyncGateRequest interface {
	Sync(blockHash Hash, blockChan chan Block) error
}
