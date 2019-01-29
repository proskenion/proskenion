package repository

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/datastructure"
)

type Blockchain struct {
	tx      core.DBATx
	factory model.ModelFactory
	tree    core.MerklePatriciaTree
}

var (
	BlockChainRootKey     byte = 9
	BlockChainTopRootKey  byte = 0
	BlockChainNowRootKey  byte = 1
	BlockChainNextRootKey byte = 2
)

type ByteWrapper struct {
	B []byte
}

func (b *ByteWrapper) Marshal() ([]byte, error) {
	return b.B, nil
}

func (b *ByteWrapper) Unmarshal(x []byte) error {
	b.B = x
	return nil
}

func NewBlockchainFromTopBlock(tx core.DBATx, factory model.ModelFactory, cryptor core.Cryptor, topBlockHash model.Hash) (core.Blockchain, error) {
	rootHash := model.Hash(nil)
	if topBlockHash != nil {
		bw := &ByteWrapper{nil}
		if err := tx.Load(BlockHashToMappingKey(topBlockHash), bw); err != nil {
			return nil, err
		}
		rootHash = bw.B
	}
	tree, err := datastructure.NewMerklePatriciaTree(tx, cryptor, rootHash, BlockChainRootKey)
	if err != nil {
		return nil, err
	}
	return &Blockchain{tx, factory, tree}, nil
}

func BlockHashToKey(blockHash model.Hash) []byte {
	return append([]byte{BlockChainRootKey, BlockChainNowRootKey}, blockHash...)
}

func BlockHashToMappingKey(blockHash model.Hash) []byte {
	return append([]byte{BlockChainRootKey, BlockChainTopRootKey}, blockHash...)
}

func BlockHashToNextKey(blockHash model.Hash) []byte {
	return append([]byte{BlockChainRootKey, BlockChainNextRootKey}, blockHash...)
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
	blockHash := block.Hash()
	if err != nil {
		return err
	}
	it, err := b.tree.Upsert(&KVNode{BlockHashToKey(blockHash), block})
	if err != nil {
		return err
	}
	if _, err := b.tree.Upsert(&KVNode{BlockHashToNextKey(block.GetPayload().GetPreBlockHash()),
		&ByteWrapper{blockHash}}); err != nil {
		return err
	}

	rootHash := it.Hash()
	if err != nil {
		return err
	}
	if err := b.tx.Store(BlockHashToMappingKey(blockHash), &ByteWrapper{rootHash}); err != nil {
		return err
	}
	return nil
}
