package consensus

import (
	"github.com/inconshreveable/log15"
	"github.com/proskenion/proskenion/core"
)

type Consensus struct {
	commitChan chan interface{}
	cc         core.ConsensusCustomize
	cs         core.CommitSystem
	rp         core.Repository
	gossip     core.Gossip
	logger     log15.Logger
}

func (c *Consensus) Boot() {
	for {
		c.cc.WaitUntilComeNextBlock()

		// TODO 1. 自分が Block の生成者か判定
		if c.cc.IsBlockCreator() {
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
