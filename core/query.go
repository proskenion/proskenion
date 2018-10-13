package core

import (
	"fmt"
	"github.com/proskenion/proskenion/core/model"
)

var (
	ErrQueryProcessorQueryEmptyBlockchain          = fmt.Errorf("Failed QueryProcessor Query blockchain is empty")
	ErrQueryProcessorQueryObjectCodeNotImplemented = fmt.Errorf("Failed QueryProcessor Query ObjectCode is not implemented")
	ErrQueryProcessorNotFound                      = fmt.Errorf("Failed QueryProcessor Query Not Found")
)

type QueryProcessor interface {
	Query(query model.Query) (model.QueryResponse, error)
}

type QueryValidator interface {
	Validate(query model.Query) error
}
