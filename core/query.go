package core

import "github.com/proskenion/proskenion/core/model"

type QueryValidator interface {
	Validate(query model.Query) error
}
