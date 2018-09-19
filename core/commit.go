package core

import . "github.com/proskenion/proskenion/core/model"

type Commit interface {
	Commit(wsv WSV, blockchain Blockchain, block Block, tree MerkleTree) error
}
