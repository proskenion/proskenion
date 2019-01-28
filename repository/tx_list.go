package repository

import (
	"fmt"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/datastructure"
	"io/ioutil"
)

type TxList struct {
	tree core.MerkleTree
	txs  []model.Transaction
	fc   model.ModelFactory
}

func NewTxListFromConf(cryptor core.Cryptor, fc model.ModelFactory, pr core.Prosl, conf *config.Config) (core.TxList, error) {
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
	txList := NewTxList(cryptor,fc)
	if err := txList.Push(ret.GetTransaction()); err != nil {
		return nil, err
	}
	return txList, nil
}

func NewTxList(cryptor core.Cryptor, fc model.ModelFactory) core.TxList {
	return &TxList{datastructure.NewAccumulateHash(cryptor), make([]model.Transaction, 0), fc}
}

func (t *TxList) Push(tx model.Transaction) error {
	t.txs = append(t.txs, tx)
	return t.tree.Push(tx)
}

func (t *TxList) Hash() model.Hash {
	return t.tree.Hash()
}

func (t *TxList) List() []model.Transaction {
	return t.txs
}

func (t *TxList) Size() int {
	return len(t.txs)
}

func (t *TxList) Marshal() ([]byte, error) {
	mbytes := make([][]byte, 0, len(t.txs))
	for _, tx := range t.txs {
		ret, err := tx.Marshal()
		if err != nil {
			return nil, err
		}
		mbytes = append(mbytes, ret)
	}
	return model.GobMarshal(mbytes)
}

func (t *TxList) Unmarshal(b []byte) error {
	mbytes := make([][]byte, 0)
	if err := model.GobUnmarshal(b, &mbytes); err != nil {
		return err
	}
	for _, bytes := range mbytes {
		etx := t.fc.NewEmptyTx()
		if err := etx.Unmarshal(bytes); err != nil {
			return err
		}
		if err := t.Push(etx); err != nil {
			return err
		}
	}
	return nil
}
