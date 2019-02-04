package p2p

import (
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"go.uber.org/multierr"
	"sync"
)

type BroadCastGossip struct {
	rp core.Repository
	fc model.ModelFactory
	cf core.ClientFactory
	c  core.Cryptor

	conf *config.Config
}

func NewBroadCastGossip(rp core.Repository, fc model.ModelFactory, cf core.ClientFactory, c core.Cryptor, conf *config.Config) core.Gossip {
	return &BroadCastGossip{rp, fc, cf, c, conf}
}

func (g *BroadCastGossip) GossipBlock(block model.Block, txList core.TxList) error {
	wsv, err := g.rp.TopWSV()
	if err != nil {
		return err
	}
	unmarshalers, err := wsv.QueryAll(model.MustAddress("/"+model.PeerStorageName), model.NewPeerUnmarshalerFactory(g.fc))
	if err != nil {
		return err
	}
	if err := core.CommitTx(wsv); err != nil {
		return err
	}

	var errs error
	mutex := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	for _, unmarshaler := range unmarshalers {
		peer := unmarshaler.(model.Peer)
		if peer.GetPeerId() == g.conf.Peer.Id {
			continue
		}
		client, err := g.cf.ConsensusClient(peer)
		if err != nil {
			errs = multierr.Append(errs, err)
			continue
		}
		wg.Add(1)
		go func(block model.Block, txList core.TxList) {
			err := client.PropagateBlockStreamTx(block, txList)
			if err != nil {
				mutex.Lock()
				errs = multierr.Append(errs, err)
				mutex.Unlock()
			}
			wg.Done()
		}(block, txList)
	}
	wg.Wait()
	if errs != nil {
		return errs
	}
	return nil
}
