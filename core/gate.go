package core

import (
	"fmt"
	. "github.com/proskenion/proskenion/core/model"
)

// API
var (
	ErrAPIWriteVerifyError    = fmt.Errorf("Failed API Write Stateless Verify Error")
	ErrAPIWriteTxAlreadyExist = fmt.Errorf("Failed API Write Transaction is already exists")
	ErrAPIWriteGossipTxError = fmt.Errorf("Failed API Write Gossip Tx occures error")

	ErrAPIQueryVerifyError   = fmt.Errorf("Failed API Read query Verify Error")
	ErrAPIQueryValidateError = fmt.Errorf("Failed API Read query Validate Error")
	ErrAPIQueryNotFound      = fmt.Errorf("Failed API Read query not found")
)

type API interface {
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

	PropagateBlockAck(block Block) (Signature, error)
	PropagateBlockStreamTx(block Block, txChan chan Transaction, errChan chan error) error
}


// Sync
var (
	ErrSyncGateSyncNotFoundBlockHash = fmt.Errorf("Failde SyncGate Sync Not found blockHash.")
)

type SyncGate interface {
	Sync(blockHash Hash, blockChan chan Block, txListChan chan TxList) error
}
