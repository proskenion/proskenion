package core

import "github.com/proskenion/proskenion/core/model"

// 伝搬アルゴリズム
type Gossip interface {
	GossipBlock(model.Block, TxList) error
}

type Fetcher interface {
	NewClient(peer model.Peer) error
	FetchTxList(block model.Block) (TxList, error)
}