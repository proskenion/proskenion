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

// ConsensusGateServer is the server Consensus for ConsensusGate service.
type ConsensusGateServer struct {
	fc     model.ModelFactory
	cg     core.ConsensusGate
	logger log15.Logger
}

func NewConsensusGateServer(fc model.ModelFactory, cg core.ConsensusGate, logger log15.Logger) proskenion.ConsensusGateServer {
	return &ConsensusGateServer{
		fc,
		cg,
		logger,
	}
}

func (s *ConsensusGateServer) PropagateTx(ctx context.Context, tx *proskenion.Transaction) (*proskenion.ConsensusResponse, error) {
	modelTx := s.fc.NewEmptyTx()
	modelTx.(*convertor.Transaction).Transaction = tx

	if err := s.cg.PropagateTx(modelTx); err != nil {
		s.logger.Error(err.Error())
		if errors.Cause(err) == core.ErrConsensusGatePropagateTxVerifyError {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proskenion.ConsensusResponse{}, nil
}

func (s *ConsensusGateServer) PropagateBlock(ctx context.Context, block *proskenion.Block) (*proskenion.ConsensusResponse, error) {
	modelBlock := s.fc.NewEmptyBlock()
	modelBlock.(*convertor.Block).Block = block

	if err := s.cg.PropagateBlock(modelBlock); err != nil {
		s.logger.Error(err.Error())
		if errors.Cause(err) == core.ErrConsensusGatePropagateBlockVerifyError {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Cause(err) == core.ErrConsensusGatePropagateBlockAlreadyExist {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proskenion.ConsensusResponse{}, nil
}

func (s *ConsensusGateServer) CollectTx(req *proskenion.CollectTxRequest, stream proskenion.ConsensusGate_CollectTxServer) error {
	txChan := make(chan model.Transaction)
	errChan := make(chan error)
	defer close(errChan)
	defer close(txChan)

	go func(errChan chan error) {
		for tx := range txChan {
			errChan <- stream.Send(tx.(*convertor.Transaction).Transaction)
		}
	}(errChan)
	if err := s.cg.CollectTx(req.GetBlockHash(), txChan, errChan); err != nil {
		s.logger.Error(err.Error())
		return status.Error(codes.Internal, err.Error())
	}
	return nil
}
