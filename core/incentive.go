package core

import (
	. "github.com/proskenion/proskenion/core/model"
)

type IncentiveTxBuilder interface {
	Build(wsv WSV, txList TxList) Transaction
}
