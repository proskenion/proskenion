package controller

import (
	"fmt"
	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/config"
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
	c      core.Cryptor
	logger log15.Logger

	conf *config.Config
}

func NewConsensusGateServer(fc model.ModelFactory, cg core.ConsensusGate, c core.Cryptor, logger log15.Logger, conf *config.Config) proskenion.ConsensusGateServer {
	return &ConsensusGateServer{
		fc,
		cg,
		c,
		logger,
		conf,
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

func (s *ConsensusGateServer) internalError(err error) error {
	s.logger.Error(err.Error())
	return status.Error(codes.Internal, err.Error())
}

func (s *ConsensusGateServer) PropagateBlock(stream proskenion.ConsensusGate_PropagateBlockServer) error {
	req, err := stream.Recv()
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}
	block := req.GetBlock()
	if block == nil {
		s.logger.Error("First receive block is nil.")
		return status.Error(codes.InvalidArgument, "First receive block is nil.")
	}

	modelBlock := s.fc.NewEmptyBlock()
	modelBlock.(*convertor.Block).Block = block

	// PropagateBlockAck Reply :: signature...
	signature, err := s.cg.PropagateBlockAck(modelBlock)
	if err != nil {
		s.logger.Error(err.Error())
		return s.internalError(err)
	}
	// send signature
	if err := stream.Send(&proskenion.PropagateBlockResponse{
		Signature: &proskenion.Signature{
			PublicKey: signature.GetPublicKey(),
			Signature: signature.GetSignature(),
		}}); err != nil {
		if errors.Cause(err) == core.ErrConsensusGatePropagateBlockVerifyError {
			return status.Errorf(codes.InvalidArgument, err.Error())
		}
		return s.internalError(err)
	}

	txChan := make(chan model.Transaction)
	errChan := make(chan error)

	go func() {
		defer close(txChan)
		defer close(errChan)

		for i := 0; ; i++ {
			req, err := stream.Recv()
			if err != nil {
				errChan <- err
				return
			}
			tx := req.GetTransaction()
			if tx == nil {
				errChan <- fmt.Errorf("%d-th transaction is nil", i)
				return
			}
			modelTx := s.fc.NewEmptyTx()
			modelTx.(*convertor.Transaction).Transaction = tx
			txChan <- modelTx
		}
	}()

	if err := s.cg.PropagateBlockStreamTx(modelBlock, txChan, errChan); err != nil {
		s.logger.Error(err.Error())
		cause := errors.Cause(err)
		if cause == core.ErrConsensusGatePropagateBlockAlreadyExist {
			return status.Error(codes.AlreadyExists, err.Error())
		} else if cause == core.ErrConsensusGatePropagateBlockDifferentHash {
			return status.Error(codes.InvalidArgument, err.Error())
		}
		return s.internalError(err)
	}
	return nil
}
