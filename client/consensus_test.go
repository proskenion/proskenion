package client

import (
	"github.com/proskenion/proskenion/config"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestConsensusClient_PropagateBlockStreamTx(t *testing.T) {
	conf := RandomConfig()
	s := RandomServer()

	// server setup
	go func(conf *config.Config, server *grpc.Server) {
		RandomSetUpConsensusServer(t, conf, s)
	}(conf, s)
	time.Sleep(time.Second * 2)

	client, err := NewConsensusClient(RandomFactory().NewPeer(conf.Peer.Id, "127.0.0.1:"+conf.Peer.Port, conf.Peer.PublicKeyBytes()), RandomFactory(), RandomCryptor())
	require.NoError(t, err)
	block, txList := RandomValidSignedBlockAndTxList(t)
	require.NoError(t, client.PropagateBlockStreamTx(block, txList))

	s.GracefulStop()
}
