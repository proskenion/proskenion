package core

import (
	. "github.com/proskenion/proskenion/core/model"
)

type APIClient interface {
	Write(in Transaction) error
	Read(in Query) (QueryResponse, error)
}

type ConsensusClient interface {
	PropagateTx(tx Transaction) error
	PropagateBlockStreamTx(block Block, txLit TxList) error
}

type ClientFactory interface {
	APIClient(peer Peer) (APIClient, error)
	ConsensusClient(peer Peer) (ConsensusClient, error)
}
