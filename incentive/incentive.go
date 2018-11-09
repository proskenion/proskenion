package incentive

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type IncentiveTxBuilder struct {
	mf model.ModelFactory
}

func NewIncentiveTxBuilder(mf model.ModelFactory) core.IncentiveTxBuilder {
	return &IncentiveTxBuilder{mf}
}

func (i *IncentiveTxBuilder) Build(wsv core.WSV, txList core.TxList) model.Transaction {
	return i.mf.NewEmptyTx()
}
