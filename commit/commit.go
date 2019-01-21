package commit

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"time"
)

type CommitSystem struct {
	factory model.ModelFactory
	cryptor core.Cryptor
	queue   core.ProposalTxQueue
	rp      core.Repository
	conf    *config.Config
}

func NewCommitSystem(factory model.ModelFactory, cryptor core.Cryptor, queue core.ProposalTxQueue, rp core.Repository, conf *config.Config) core.CommitSystem {
	return &CommitSystem{factory, cryptor, queue, rp, conf}
}

func UnixTime(t time.Time) int64 {
	return t.UnixNano()
}

func Now() int64 {
	return UnixTime(time.Now())
}

// Stateless Validate
func (c *CommitSystem) VerifyCommit(block model.Block, txList core.TxList) error {
	if err := block.Verify(); err != nil {
		return errors.Wrapf(core.ErrCommitSystemVerifyCommitBlockVerify, err.Error())
	}
	if !bytes.Equal(block.GetPayload().GetTxsHash(), txList.Top()) {
		return errors.Wrapf(core.ErrCommitSystemVerifyCommitNotMatchedTxsHash,
			"block.txsHash:%x, actual txsHash:%x", block.GetPayload().GetTxsHash(), txList.Top())
	}
	for _, tx := range txList.List() {
		if err := tx.Verify(); err != nil {
			return errors.Wrapf(core.ErrCommitSystemVerifyCommitTxVerify, err.Error())
		}
	}
	return nil
}

func getTopWSV(rp core.Repository) (core.RepositoryTx, core.WSV, model.Block, error) {
	rtx, err := rp.Begin()
	if err != nil {
		return nil, nil, nil, err
	}
	top, _ := rp.Top()
	wsv, err := rtx.WSV(top.GetPayload().GetWSVHash())
	if err != nil {
		return nil, nil, nil, err
	}
	return rtx, wsv, top, nil
}

func (c *CommitSystem) ValidateCommit(block model.Block, txList core.TxList) error {
	acs, err := c.rp.GetDelegatedAccounts()
	if err != nil {
		return err
	}
	if len(acs) <= int(block.GetPayload().GetRound()) {
		return errors.Wrapf(core.ErrCommitSystemValidateCommitRoundOutOfRange,
			"round: %d", block.GetPayload().GetRound())
	}
	peerId := model.MustAddress(model.MustAddress(acs[block.GetPayload().GetRound()].GetDelegatePeerId()).PeerId())
	rtx, wsv, top, err := getTopWSV(c.rp)
	peer := c.factory.NewEmptyPeer()
	if err := wsv.Query(peerId, peer); err != nil {
		return errors.Wrapf(core.ErrCommitSystemValidateCommitInternal, err.Error())
	}
	if !bytes.Equal(peer.GetPublicKey(), block.GetSignature().GetPublicKey()) {
		return errors.Wrapf(core.ErrCommitSystemValidateCommitInvalidPeer,
			"expected peer: %s, expected pubkey: %x, actual: %x",
			peer.GetPeerId(), peer.GetPublicKey(), block.GetSignature().GetPublicKey())
	}
	if !bytes.Equal(top.Hash(), block.GetPayload().GetPreBlockHash()) {
		return errors.Wrapf(core.ErrCommitSystemValidateCommitInvalidPreBlock,
			"expected: %x, but receive block's preBlockHash: %x", top.Hash(), block.GetPayload().GetPreBlockHash())

	}
	expAfter := top.GetPayload().GetCreatedTime() + int64(c.conf.Commit.WaitInterval)*int64(block.GetPayload().GetRound())
	now := Now()
	if expAfter > now {
		return errors.Wrapf(core.ErrCommitSystemValidateCommitSoFastTime,
			"round: %d, expected after time: %d, but now: %d", block.GetPayload().GetRound(), expAfter, now)
	}

	return rtx.Commit()
}

func (c *CommitSystem) Commit(block model.Block, txList core.TxList) error {
	return c.rp.Commit(block, txList)
}

// CreateBlock
func (c *CommitSystem) CreateBlock(round int32) (model.Block, core.TxList, error) {
	return c.rp.CreateBlock(c.queue, round, Now())
}
