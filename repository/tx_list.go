package repository

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type TxList struct {
	tree core.MerkleTree
	txs  []model.Transaction
}

func NewTxList(cryptor core.Cryptor) core.TxList {
	return &TxList{NewAccumulateHash(cryptor), make([]model.Transaction, 0)}
}

func (t *TxList) Push(tx model.Transaction) error {
	t.txs = append(t.txs, tx)
	return t.tree.Push(tx)
}

func (t *TxList) Top() model.Hash {
	return t.tree.Top()
}

func (t *TxList) List() []model.Transaction {
	return t.txs
}
