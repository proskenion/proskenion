package core

import (
	"fmt"
	. "github.com/proskenion/proskenion/core/model"
)

var (
	ErrCommitSystemVerifyCommitBlockVerify       = fmt.Errorf("Failed Verify Commit block verify.")
	ErrCommitSystemVerifyCommitNotMatchedTxsHash = fmt.Errorf("Failed Verify Commit not matched txsHash.")
	ErrCommitSystemVerifyCommitTxVerify          = fmt.Errorf("Failed Verify Commit tx verify.")
)

var (
	ErrCommitSystemValidateCommitInternal        = fmt.Errorf("Failed Validate Commit internal error.")
	ErrCommitSystemValidateCommitRoundOutOfRange = fmt.Errorf("Failed Validate Commit round out of range.")
	ErrCommitSystemValidateCommitInvalidPeer     = fmt.Errorf("Failed Validate Commit invalid consensus peer create this block.")
	ErrCommitSystemValidateCommitSoFastTime     = fmt.Errorf("Failed Validate Commit so fast time for this round.")
	ErrCommitSystemValidateCommitInvalidPreBlock = fmt.Errorf("Failed Validate Commit invalid preblock hash is different.")
)

type CommitSystem interface {
	VerifyCommit(block Block, txList TxList) error
	ValidateCommit(block Block, txList TxList) error
	Commit(block Block, txList TxList) error
	CreateBlock(round int32) (Block, TxList, error)
}
