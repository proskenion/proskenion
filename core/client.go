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
}

type ClientFactory interface {
	APIClient(peer Peer) (APIGateClient, error)
	ConsensusClient(peer Peer) (ConsensusGateClient, error)
}
