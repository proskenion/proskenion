package client

import (
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestNewSyncClientSync(t *testing.T) {
	conf := RandomConfig()
	s := RandomServer()
	rp := RandomRepository()
	fc := RandomFactory()
	require.NoError(t, rp.GenesisCommit(RandomGenesisTxList(t)))
	topHash := MusTop(rp).Hash()

	for i := 0; i < conf.Sync.Limits*2+10; i++ {
		RandomCommitableBlockAndTxList(t, rp)
	}
	// server setup
	go func(conf *config.Config, server *grpc.Server) {
		RandomSetUpSyncServer(t, conf, rp, s)
	}(conf, s)
	time.Sleep(time.Second * 2)

	blockChan := make(chan model.Block)
	txListChan := make(chan core.TxList)
	errChan := make(chan error)

	retErrChan := make(chan error)
	defer close(retErrChan)

	client, err := NewSyncClient(fc.NewPeer(conf.Peer.Id, "127.0.0.1:"+conf.Peer.Port, conf.Peer.PublicKeyBytes()), fc, RandomCryptor())
	require.NoError(t, err)

	newRp := RandomRepository()
	require.NoError(t, newRp.GenesisCommit(RandomGenesisTxList(t)))
	go func() {
		defer close(blockChan)
		defer close(txListChan)
		defer close(errChan)
		err := client.Sync(topHash, blockChan, txListChan, errChan)
		require.NoError(t, err)
		retErrChan <- err
	}()
	var newBlock model.Block
	var newTxList core.TxList
	for {
		select {
		case newBlock = <-blockChan:
		case newTxList = <-txListChan:
			require.NoError(t, newRp.Commit(newBlock, newTxList))
		case err := <-retErrChan:
			require.NoError(t, err)
			goto afterFor
		}
		if MusTop(newRp).GetPayload().GetHeight() == MusTop(rp).GetPayload().GetHeight() {
			errChan <- nil
		}
	}
afterFor:
	assert.Equal(t, MusTop(rp).Hash(), MusTop(newRp).Hash())
}
