package consensus

import (
	"fmt"
	"github.com/proskenion/proskenion/commit"
	"github.com/proskenion/proskenion/core"
	"time"
)

type Customize struct {
	rp         core.Repository
	commitChan chan interface{}
	// 前回の Commit から次の Commit までの間隔の最大値
	MaxWaitngCommitInterval time.Duration
}

func (c *Customize) WaitUntilComeNextBlock() {
	top, ok := c.rp.Top()
	if !ok {
		panic("Must be Genesis Commit after boot consensus")
	}

	timer := time.NewTimer(
		c.MaxWaitngCommitInterval +
			time.Duration(top.GetPayload().GetCreatedTime()-commit.Now()))
	// commit を待つ
	select {
	case <-c.commitChan:
		return
	case <-timer.C:
		return
	}
}

func (mc *Customize) IsBlockCreator() bool {
	// TODO
	return true
}

// Mock
type MockCustomize struct {
	rp         core.Repository
	commitChan chan interface{}
}

func NewMockCustomize(rp core.Repository, commitChan chan interface{}) core.ConsensusCustomize {
	return &MockCustomize{rp, commitChan}
}

func (mc *MockCustomize) WaitUntilComeNextBlock() {
	top, ok := mc.rp.Top()
	if !ok {
		panic("Must be Genesis Commit after boot consensus")
	}
	fmt.Println("Height: ", top.GetPayload().GetHeight())

	timer := time.NewTimer(
		time.Second +
			time.Duration(top.GetPayload().GetCreatedTime()-commit.Now()))
	// commit を待つ
	select {
	case <-mc.commitChan:
		return
	case <-timer.C:
		return
	}
}

func (mc *MockCustomize) IsBlockCreator() bool {
	return true
}