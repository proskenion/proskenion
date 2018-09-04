package core

import (
	"github.com/pkg/errors"
	. "github.com/proskenion/proskenion/core/model"
)

var (
	ErrStatefulValidate       = errors.Errorf("Failed StatefulValidator Validate")
	ErrStatelessBlockValidate = errors.Errorf("Failed StatelessBlockValidator Validate")
	ErrStatelessTxValidate    = errors.Errorf("Failed StatelessTxValidator Validate")
)

type StatefulValidator interface {
	Validate(block Block) error
}

type StatelessValidator interface {
	BlockValidate(block Block) error
	TxValidate(tx Transaction) error
}
