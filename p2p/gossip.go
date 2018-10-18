package p2p

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type Gossip struct {
	rp core.Repository
}

func (g *Gossip) GossipBlock(block model.Block, list core.TxList) error {
	return nil
}
