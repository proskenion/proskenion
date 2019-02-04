package synchronize_test

import (
	"github.com/proskenion/proskenion/client"
	"github.com/proskenion/proskenion/config"
	. "github.com/proskenion/proskenion/synchronize"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"io"
	"sync"
	"testing"
	"time"
)

func TestSynchronizer_Sync(t *testing.T) {
	conf := RandomConfig()
	conf.Peer.Port = "60001"
	s := RandomServer()
	rp := RandomRepository()
	fc := RandomFactory()
	require.NoError(t, rp.GenesisCommit(RandomGenesisTxList(t)))
	for i := 0; i < conf.Sync.Limits*2+10; i++ {
		RandomCommitableBlockAndTxList(t, rp)
	}
	// server setup
	go func(conf *config.Config, server *grpc.Server) {
		RandomSetUpSyncServer(t, conf, rp, s)
	}(conf, s)
	time.Sleep(time.Second)

	t.Run("case 1 : success", func(t *testing.T) {
		newRp := RandomRepository()
		cf := client.NewClientFactory(fc, RandomCryptor(), conf)
		require.NoError(t, newRp.GenesisCommit(RandomGenesisTxList(t)))

		active := false
		peer := fc.NewPeer(conf.Peer.Id, conf.Peer.Host+":"+conf.Peer.Port, conf.Peer.PublicKeyBytes())
		syn := NewSynchronizer(newRp, cf, &active)
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			err := syn.Sync(peer)
			assert.Equal(t, io.EOF, err)
			wg.Done()
		}()

		for {
			if MusTop(rp).GetPayload().GetHeight() == MusTop(newRp).GetPayload().GetHeight() {
				assert.Equal(t, MusTop(rp).Hash(), MusTop(newRp).Hash())
				active = true
				break
			}
			time.Sleep(time.Second)
		}
		wg.Wait()
	})
	s.GracefulStop()
}
