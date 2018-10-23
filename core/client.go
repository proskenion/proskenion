package core

import (
	"github.com/proskenion/proskenion/core/model"
)

type APIGateClient interface {
	Write(in model.Transaction) error
	Read(in model.Query) (model.QueryResponse, error)
}
