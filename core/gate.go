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

type ConsensusGate interface {
	PropagateTx(tx Transaction) error
	PropagateBlock(block Block) error
	// chan が Stream の返り値
	CollectTx(blockHash Hash, txChan chan Transaction) error
}

type SyncGate interface {
	Sync(blockHash Hash, blockChan chan Block) error
}
