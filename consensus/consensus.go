package consensus

import (
	"github.com/inconshreveable/log15"
	"github.com/proskenion/proskenion/commit"
	"github.com/proskenion/proskenion/core"
	"time"
)

type Consensus struct {
	commitChan chan interface{}
	cs         core.CommitSystem
	rp         core.Repository
	gossip     core.Gossip
	logger     log15.Logger
}

func (c *Consensus) Boot() {
	top, ok := c.rp.Top()
	if !ok {
		panic("Must be genesis commit after call this func.")
	}
	for {
		timer := time.NewTimer(time.Duration(top.GetPayload().GetCreatedTime() - commit.Now()))
		// commit を待つ
		select {
		case <-c.commitChan:
			break
		case <-timer.C:
			break
		}
		top, ok = c.rp.Top()

		// 1. 自分が Block の生成者か判定
		if true {
			// 2. block を生成
			block, txList, err := c.cs.CreateBlock()
			if err != nil {
				c.logger.Error(err.Error())
				continue
			}
			// 3. block を Gosship
			err = c.gossip.GossipBlock(block, txList)
			if err != nil {
				c.logger.Error(err.Error())
				continue
			}
		}
	}

}
