package repository

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type Blockchain struct {
	tx      core.DBATx
	factory model.ModelFactory
	tree    core.MerklePatriciaTree
}

var BLOCKCHAIN_ROOT_KEY byte = 3

func NewBlockchain(tx core.DBATx, factory model.ModelFactory, cryptor core.Cryptor, rootHash model.Hash) (core.Blockchain, error) {
	tree, err := NewMerklePatriciaTree(tx, cryptor, rootHash, BLOCKCHAIN_ROOT_KEY)
	if err != nil {
		return nil, err
	}
	return &Blockchain{tx, factory, tree}, nil
}

func BlockHashToKey(blockHash model.Hash) []byte {
	return append([]byte{BLOCKCHAIN_ROOT_KEY}, blockHash...)
}

func (b *Blockchain) Get(blockHash model.Hash) (model.Block, error) {
	blockHash = BlockHashToKey(blockHash)
	it, err := b.tree.Find(blockHash)
	if err != nil {
		if errors.Cause(err) == core.ErrMerklePatriciaTreeNotFoundKey {
			return nil, errors.Wrap(core.ErrBlockchainNotFound, err.Error())
		}
		return nil, err
	}
	retBlock := b.factory.NewEmptyBlock()
	if err = it.Data(retBlock); err != nil {
		return nil, errors.Wrap(core.ErrBlockchainQueryUnmarshal, err.Error())
	}
	return retBlock, nil
}

// Commit is allowed only Commitable Block, ohterwise panic
func (b *Blockchain) Append(block model.Block) (err error) {
	blockHash, err := block.Hash()
	blockHash = BlockHashToKey(blockHash)
	if err != nil {
		return err
	}
	_, err = b.tree.Upsert(&KVNode{blockHash, block})
	return err
}
