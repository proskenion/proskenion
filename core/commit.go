package core

import . "github.com/proskenion/proskenion/core/model"

type Commit interface {
	Commit(wfa WFA, blockchain Blockchain, block Block, tree MerkleTree) error
}
