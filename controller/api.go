package controller

import (
	"github.com/inconshreveable/log15"
	"github.com/proskenion/proskenion/convertor"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/gate"
	"github.com/proskenion/proskenion/proto"
	"golang.org/x/net/context"
)

// APIGateServer is the server API for APIGate service.
type APIGateServer struct {
	fc     model.ModelFactory
	logger log15.Logger
	api    *gate.APIGate
}

func (s *APIGateServer) Write(ctx context.Context, tx *proskenion.Transaction) (*proskenion.TxResponse, error) {
	modelTx := s.fc.NewEmptyTx()
	modelTx.(*convertor.Transaction).Transaction = tx

	err := s.api.Write(modelTx)
	if err != nil {
		// grpc error code
	}
	return &proskenion.TxResponse{}, nil
}

func (s *APIGateServer) Read(ctx context.Context, query *proskenion.Query) (*proskenion.QueryResponse, error) {
	modelQuery := s.fc.NewEmptyQuery()
	modelQuery.(*convertor.Query).Query = query

	res, err := s.api.Read(modelQuery)
	if err != nil {
		// grpc error code
	}
	return res.(*convertor.QueryResponse).QueryResponse, nil
}
