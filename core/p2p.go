package core

import "github.com/proskenion/proskenion/core/model"

// 伝搬アルゴリズム
type Gossip interface {
	GossipBlock(model.Block, TxList) error
}

// Peer 取得機構
type PeerService interface {
	List() []model.Peer
}