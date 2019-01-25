package core

import (
	. "github.com/proskenion/proskenion/core/model"
)

type APIGateClient interface {
	Write(in Transaction) error
	Read(in Query) (QueryResponse, error)
}

type ConsensusGateClient interface {
	PropagateTx(tx Transaction) error
	PropagateBlockStreamTx(block Block, txLit TxList) error
	// chan が Stream の返り値
	CollectTx(blockHash Hash, txChan chan Transaction) error
}
