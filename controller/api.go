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

// APIServer is the server API for API service.
type APIServer struct {
	fc     model.ModelFactory
	api    core.API
	logger log15.Logger
}

func NewAPIServer(fc model.ModelFactory, api core.API, logger log15.Logger) proskenion.APIServer {
	return &APIServer{
		fc,
		api,
		logger,
	}
}

func (s *APIServer) Write(ctx context.Context, tx *proskenion.Transaction) (*proskenion.TxResponse, error) {
	modelTx := s.fc.NewEmptyTx()
	modelTx.(*convertor.Transaction).Transaction = tx

	err := s.api.Write(modelTx)
	if err != nil {
		s.logger.Error(err.Error())
		if errors.Cause(err) == core.ErrAPIWriteVerifyError {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Cause(err) == core.ErrAPIWriteTxAlreadyExist {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proskenion.TxResponse{}, nil
}

func (s *APIServer) Read(ctx context.Context, query *proskenion.Query) (*proskenion.QueryResponse, error) {
	modelQuery := s.fc.NewEmptyQuery()
	modelQuery.(*convertor.Query).Query = query

	res, err := s.api.Read(modelQuery)
	if err != nil {
		s.logger.Error(err.Error())
		if errors.Cause(err) == core.ErrAPIQueryVerifyError ||
			errors.Cause(err) == core.ErrAPIQueryValidateError {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Cause(err) == core.ErrAPIQueryNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return res.(*convertor.QueryResponse).QueryResponse, nil
}
