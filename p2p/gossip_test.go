package p2p

import (
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGossip_GossipBlock(t *testing.T) {
	fc, _, _, c, rp, _, conf := NewTestFactories()
	cf := NewMockClientFactory()
	gossip := NewBroadCastGossip(rp, fc, cf, c, conf)

	// previous setting, commit genesis
	txList := RandomGenesisTxList(t)
	require.NoError(t, rp.GenesisCommit(txList))

	ps := []model.Peer{
		RandomPeer(),
		RandomPeer(),
		RandomPeer(),
		RandomPeer(),
		RandomPeer(),
	}

	b := fc.NewTxBuilder()
	for _, p := range ps {
		b = b.AddPeer("root@root", p.GetPeerId(), p.GetAddress(), RandomPublicKey())
	}
	CommitTxWrapBlock(t, rp, fc, b.Build())

	block, txList := RandomValidSignedBlockAndTxList(t)
	err := gossip.GossipBlock(block, txList)
	require.NoError(t, err)
	for _, p := range ps {
		client, err := cf.ConsensusClient(p)
		require.NoError(t, err)
		assert.Equal(t, client.(*MockConsensusClient).Id, p.GetPeerId())
		assert.Equal(t, client.(*MockConsensusClient).PropagateBlockIn1.Hash(), block.Hash())
		assert.Equal(t, client.(*MockConsensusClient).PropagateBlockIn2.Hash(), txList.Hash())
	}

	nq, err := cf.ConsensusClient(fc.NewPeer(conf.Peer.Id, ":", RandomPublicKey()))
	assert.Nil(t, nq.(*MockConsensusClient).PropagateBlockIn1)
	assert.Nil(t, nq.(*MockConsensusClient).PropagateBlockIn2)
}
