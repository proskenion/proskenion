package core

import (
	. "github.com/proskenion/proskenion/core/model"
)

type CommitSystem interface {
	VerifyCommit(block Block, txList TxList) error
	Commit(block Block, txList TxList) error
	CreateBlock() (Block, TxList, error)
}
