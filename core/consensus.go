package core

import (
	. "github.com/proskenion/proskenion/core/model"
)

type Consensus interface {
	// TODO
	CommitValidate(block Block) error
}

type ConsensusCustomize interface {
	WaitUntilComeNextBlock()
	IsBlockCreator() bool
}
