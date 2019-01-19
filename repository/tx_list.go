package repository

import (
	"fmt"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"io/ioutil"
)

type TxList struct {
	tree core.MerkleTree
	txs  []model.Transaction
}

func NewTxListFromConf(cryptor core.Cryptor, pr core.Prosl, conf *config.Config) (core.TxList, error) {
	buf, err := ioutil.ReadFile(conf.Prosl.Genesis.Path)
	if err != nil {
		return nil, err
	}
	err = pr.ConvertFromYaml(buf)
	if err != nil {
		return nil, err
	}
	ret, vars, err := pr.Execute()
	if err != nil {
		return nil, fmt.Errorf("Error Genesis prosl: %s\nvariables: %+v\n", err.Error(), vars)
	}
	if ret.GetTransaction() == nil {
		return nil, fmt.Errorf("Error Genesis prosl return nil.")
	}
	txList := NewTxList(cryptor)
	if err := txList.Push(ret.GetTransaction()); err != nil {
		return nil, err
	}
	return txList, nil
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

func (t *TxList) Size() int {
	return len(t.txs)
}
