package incentive

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type TempIncentiveTxBuilder struct {
	fc model.ModelFactory
}

func NewTemporaryIncentiveTxBuilder(fc model.ModelFactory) core.IncentiveTxBuilder {
	return &IncentiveTxBuilder{fc}
}

//
// Temporary Incentive Tx Builder の概要
//
// ブロックの生成者に "root@domain/token" を Transfer する
//
func (i *TempIncentiveTxBuilder) Build(rpTx core.RepositoryTx, block model.Block, txList core.TxList) (model.Transaction, error) {
	builder := i.fc.NewTxBuilder()
	for _, tx := range txList.List() {
	}
	return builder.Build(), nil
}
