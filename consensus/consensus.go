package consensus

import (
	"fmt"
	"github.com/inconshreveable/log15"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/core"
	"time"
)

type Consensus struct {
	rp     core.Repository
	cs     core.CommitSystem
	bq     core.ProposalBlockQueue

	gossip core.Gossip
	logger log15.Logger
	pr     core.Prosl
	conf   *config.Config

	commitChan chan struct{}
	// 前回の Commit から次の Commit までの間隔の最大値
	WaitngInterval time.Duration
}

func NewConsensus(rp core.Repository, cs core.CommitSystem, bq core.ProposalBlockQueue, gossip core.Gossip, pr core.Prosl,
	logger log15.Logger, conf *config.Config, commitChan chan struct{}) core.Consensus {
	return &Consensus{rp, cs, bq, gossip, logger, pr, conf,
		commitChan, time.Duration(conf.Commit.WaitInterval) * time.Millisecond}
}

type ConsensusWaitFlag int32

const (
	UpdateFlag ConsensusWaitFlag = iota
	TimeOutFlag
)

func (c *Consensus) waitUntilComeNextBlock() ConsensusWaitFlag {
	timer := time.NewTimer(c.WaitngInterval)
	// commit を待つ
	select {
	case <-c.commitChan:
		return UpdateFlag
	case <-timer.C:
		return TimeOutFlag
	}
}

func (c *Consensus) isBlockCreator(round int) bool {
	acs, err := c.rp.GetDelegatedAccounts()
	if err != nil {
		c.logger.Error("GetDelegatePeers error: %s", err.Error())
	}
	if len(acs) <= round {
		return false
	}
	if c.conf.Peer.Id == acs[round].GetDelegatePeerId() {
		return true
	}
	return false
}

func (c *Consensus) Boot() {
	// Height - loop
	for {
		top, _ := c.rp.Top()
		c.logger.Info(fmt.Sprintf(fmt.Sprintf("============= Wait Until Come Next Block %d =============", top.GetPayload().GetHeight())))
		for round := 0; ; round++ { // Round - loop
			c.logger.Info(fmt.Sprintf("        ============= Round : %d =============", round))
			// 前回の Commit から次の Commit までの間隔の最大値
			if c.waitUntilComeNextBlock() == UpdateFlag {
				// top が Update されたら New Height からスタート
				break
			}
			// top が更新されないまま一定時間経過で来なかったら -> next round

			// 1. 自分が i(i<=round) 番目の Block 生成者か判定
			if c.isBlockCreator(round) {
				c.logger.Info("============= Create Block =============")
				// 2. block を生成
				block, txList, err := c.cs.CreateBlock(int32(round))
				if err != nil {
					c.logger.Error(err.Error())
					continue
				}
				c.logger.Info(fmt.Sprintf("txLen :: %d", len(txList.List())))
				// 3. block を Gosship
				err = c.gossip.GossipBlock(block, txList)
				if err != nil {
					c.logger.Error(err.Error())
					continue
				}
				break
			}
		}
	}
}

func (c *Consensus) Receiver() {
	for {
		c.bq.WaitPush()
		_, ok := c.bq.Pop()
		if !ok {
			continue
		}

	}
}
