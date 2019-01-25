package core

import (
	"fmt"
	. "github.com/proskenion/proskenion/core/model"
)

// API
var (
	ErrAPIGateWriteVerifyError    = fmt.Errorf("Failed APIGate Write Stateless Verify Error")
	ErrAPIGateWriteTxAlreadyExist = fmt.Errorf("Failed APIGate Write Transaction is already exists")

	ErrAPIGateQueryVerifyError   = fmt.Errorf("Failed APIGate Read query Verify Error")
	ErrAPIGateQueryValidateError = fmt.Errorf("Failed APIGate Read query Validate Error")
	ErrAPIGateQueryNotFound      = fmt.Errorf("Failed APIGate Read query not found")
)

type APIGate interface {
	Write(tx Transaction) error
	Read(query Query) (QueryResponse, error)
}

// Consensus
var (
	ErrConsensusGatePropagateTxVerifyError = fmt.Errorf("Failed ConsensusGate PropagateTx Verify error")

	ErrConsensusGatePropagateBlockVerifyError   = fmt.Errorf("Failed ConsensusGate PropagateBlock Verify error")
	ErrConsensusGatePropagateBlockAlreadyExist  = fmt.Errorf("Failed ConsensusGate PropagateBlock block is already exists")
	ErrConsensusGatePropagateBlockDifferentHash = fmt.Errorf("Failed ConsensusGate PropagateBlock txList hash and block's txsHash is different")
)

type ConsensusGate interface {
	PropagateTx(tx Transaction) error
	PropagateBlock(block Block) error

	PropagateBlockAck(block Block) (Signature, error)
	PropagateBlockStreamTx(block Block, txChan chan Transaction, errChan chan error) error

	// chan が Stream の返り値
	CollectTx(blockHash Hash, txChan chan Transaction, errChan chan error) error
}

type SyncGate interface {
	Sync(blockHash Hash, blockChan chan Block) error
}
