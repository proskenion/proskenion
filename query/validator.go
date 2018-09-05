package query

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type QueryValidator struct{}

func NewQueryValidator() core.QueryValidator {
	return &QueryValidator{}
}

func (q *QueryValidator) Validate(query model.Query) error {
	return nil
}
