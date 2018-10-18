package core

type Consensus interface {
	Boot()
}

type ConsensusCustomize interface {
	WaitUntilComeNextBlock()
	IsBlockCreator() bool
}
