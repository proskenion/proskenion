package consensus

import (
	"github.com/inconshreveable/log15"
	"github.com/proskenion/proskenion/core"
)

type Consensus struct {
	cc     core.ConsensusCustomize
	cs     core.CommitSystem
	gossip core.Gossip
	logger log15.Logger
}

func NewConsensus(cc core.ConsensusCustomize, cs core.CommitSystem, gossip core.Gossip, logger log15.Logger) core.Consensus {
	return &Consensus{cc, cs, gossip, logger}
}

func (c *Consensus) Boot() {
	for {
		c.logger.Info("============= Wait Until Come Next Block =============")
		c.cc.WaitUntilComeNextBlock()

		// 1. 自分が Block の生成者か判定
		if c.cc.IsBlockCreator() {
			c.logger.Info("============= Create Block =============")
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
