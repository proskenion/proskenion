package controller

import (
	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/convertor"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// APIGateServer is the server API for APIGate service.
type APIGateServer struct {
	fc     model.ModelFactory
	api    core.APIGate
	logger log15.Logger
}

func NewAPIGateServer(fc model.ModelFactory, api core.APIGate, logger log15.Logger) proskenion.APIGateServer {
	return &APIGateServer{
		fc,
		api,
		logger,
	}
}

func (s *APIGateServer) Write(ctx context.Context, tx *proskenion.Transaction) (*proskenion.TxResponse, error) {
	modelTx := s.fc.NewEmptyTx()
	modelTx.(*convertor.Transaction).Transaction = tx

	err := s.api.Write(modelTx)
	if err != nil {
		s.logger.Error(err.Error())
		if errors.Cause(err) == core.ErrAPIGateWriteVerifyError {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Cause(err) == core.ErrAPIGateWriteTxAlreadyExist {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proskenion.TxResponse{}, nil
}

func (s *APIGateServer) Read(ctx context.Context, query *proskenion.Query) (*proskenion.QueryResponse, error) {
	modelQuery := s.fc.NewEmptyQuery()
	modelQuery.(*convertor.Query).Query = query

	res, err := s.api.Read(modelQuery)
	if err != nil {
		s.logger.Error(err.Error())
		if errors.Cause(err) == core.ErrAPIGateQueryVerifyError ||
			errors.Cause(err) == core.ErrAPIGateQueryValidateError {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Cause(err) == core.ErrAPIGateQueryNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return res.(*convertor.QueryResponse).QueryResponse, nil
}
