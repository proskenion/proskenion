package consensus

import (
	"fmt"
	"github.com/inconshreveable/log15"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"time"
)

type Consensus struct {
	rp   core.Repository
	fc   model.ModelFactory
	cs   core.CommitSystem
	sync core.Synchronizer
	bq   core.ProposalBlockQueue
	tc   core.TxListCache

	gossip core.Gossip
	logger log15.Logger
	pr     core.Prosl
	conf   *config.Config

	commitChan chan struct{}
	// 前回の Commit から次の Commit までの間隔の最大値
	WaitngInterval time.Duration
}

func NewConsensus(rp core.Repository, fc model.ModelFactory, cs core.CommitSystem, sync core.Synchronizer, bq core.ProposalBlockQueue, tc core.TxListCache,
	gossip core.Gossip, pr core.Prosl, logger log15.Logger, conf *config.Config, commitChan chan struct{}) core.Consensus {
	return &Consensus{rp, fc, cs, sync, bq, tc, gossip, logger, pr, conf,
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
		c.logger.Error(fmt.Sprintf("GetDelegatePeers error: %s", err.Error()))
		return false
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
	c.logger.Info("================= Consensus Boot =================")
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
				c.logger.Info(fmt.Sprintf("============= Create Block : height %d, round %d =============", top.GetPayload().GetHeight()+1, round))
				// 2. block を生成
				block, txList, err := c.cs.CreateBlock(int32(round))
				if err != nil {
					c.logger.Error(err.Error())
					continue
				}
				c.logger.Info(fmt.Sprintf("txLen :: %d", len(txList.List())))
				// 3. block を Gosship
				c.logger.Info("============= Gossip Block And TxList =============")
				err = c.gossip.GossipBlock(block, txList)
				if err != nil {
					c.logger.Error(err.Error())
				}
				c.logger.Info("============= Finish Gossiped  =============")
				break
			}
		}
	}
}

func (c *Consensus) syncFrom(fromPeer model.Peer) {
	c.logger.Info("================= Start Synchronize =================", "From:", fromPeer.GetPeerId())
	err := c.sync.Sync(fromPeer)
	if err != nil {
		c.logger.Error(err.Error())
	} else {
		c.rp.Me().Activate()
		c.logger.Info("============= Sucess SyncBlockChain !!! =============")
	}
}

func (c *Consensus) Patrol() {
	c.logger.Info("================= Consensus Patrol =================")
	// start sync
	// WIP : Initialize Sync only.
	if !c.rp.Me().GetActive() {
		fromPeer := config.NewPeerFromConf(c.fc, c.conf.Sync.From)
		c.syncFrom(fromPeer)
	}
	for {
		if c.rp.Me().GetActive() {
			time.Sleep(time.Second)
			continue
		}

		// start sync
		// WIP : Initialize Sync only.
		fromPeer := config.NewPeerFromConf(c.fc, c.conf.Sync.From)
		c.syncFrom(fromPeer)
	}
}

func (c *Consensus) Receiver() {
	c.logger.Info("================= Consensus Receiver =================")
	for {
		c.logger.Info("============= Wait Receive Block =============")
		c.bq.WaitPush()
		block, ok := c.bq.Pop()
		if !ok {
			continue
		}
		txList, ok := c.tc.Get(block.GetPayload().GetTxsHash())
		if !ok {
			continue
		}

		// Commit Phase
		c.logger.Info("============= Receive Block and TxList =============")
		if err := c.cs.VerifyCommit(block, txList); err != nil {
			c.logger.Error(err.Error())
			continue
		}
		if err := c.cs.ValidateCommit(block, txList); err != nil {
			c.logger.Error(err.Error())
			continue
		}
		if err := c.cs.Commit(block, txList); err != nil {
			c.logger.Error(err.Error())
			continue
		}
		c.logger.Info("============= Commit Received Block and TxList =============")
		c.commitChan <- struct{}{}
		if !c.rp.Me().GetActive() {
			c.rp.Me().Activate()
		}
	}
}
