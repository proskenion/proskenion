package p2p

import (
	"fmt"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type MockGossip struct{}

func (g *MockGossip) GossipBlock(block model.Block, list core.TxList) error {
	fmt.Println("====== Mock Gossip Block ========")
	return nil
}

func (g *MockGossip) GossipTx(tx model.Transaction) error {
	fmt.Println("===== Mock Gossip Tx ======")
	return nil
}
