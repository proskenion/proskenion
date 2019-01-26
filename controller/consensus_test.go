package controller

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/convertor"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/gate"
	"github.com/proskenion/proskenion/proto"
	"github.com/proskenion/proskenion/repository"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"io"
	"testing"
)

func newRandomConsensusGateServer() proskenion.ConsensusGateServer {
	fc, _, _, c, _, _, conf := NewTestFactories()
	cg := gate.NewConsensusGate(fc, c,
		repository.NewProposalTxQueueOnMemory(conf), repository.NewTxListCache(conf),
		repository.NewProposalBlockQueueOnMemory(conf), RandomLogger(), conf)
	return NewConsensusGateServer(fc, cg, c, RandomLogger(), RandomConfig())
}

type MockConsensusGate_PropagateBlockServer struct {
	Req chan *proskenion.PropagateBlockRequest
	Res chan *proskenion.PropagateBlockResponse
	Err chan error
	grpc.ServerStream
}

func newMockPropagateBlockServerStream() *MockConsensusGate_PropagateBlockServer {
	return &MockConsensusGate_PropagateBlockServer{make(chan *proskenion.PropagateBlockRequest),
		make(chan *proskenion.PropagateBlockResponse),
		make(chan error),
		RandomMockServerStream()}
}

func (s *MockConsensusGate_PropagateBlockServer) Send(res *proskenion.PropagateBlockResponse) error {
	s.Res <- res
	return nil
}

func (s *MockConsensusGate_PropagateBlockServer) Recv() (*proskenion.PropagateBlockRequest, error) {
	select {
	case req := <-s.Req:
		return req, nil
	case err := <-s.Err:
		return nil, err
	}
}

func createRequestBlock(block model.Block) *proskenion.PropagateBlockRequest {
	return &proskenion.PropagateBlockRequest{
		Req: &proskenion.PropagateBlockRequest_Block{block.(*convertor.Block).Block},
	}
}

func createRequestTx(tx model.Transaction) *proskenion.PropagateBlockRequest {
	return &proskenion.PropagateBlockRequest{
		Req: &proskenion.PropagateBlockRequest_Transaction{tx.(*convertor.Transaction).Transaction},
	}
}

func TestConsensusGateServer_PropagateBlock(t *testing.T) {
	ctrl := newRandomConsensusGateServer()

	t.Run("case 1 : correct", func(t *testing.T) {
		stream := newMockPropagateBlockServerStream()
		block, txList := RandomValidSignedBlockAndTxList(t)
		go func(t *testing.T) {
			defer close(stream.Req)
			defer close(stream.Res)
			defer close(stream.Err)
			stream.Req <- createRequestBlock(block)

			res := <-stream.Res
			assert.NoError(t, RandomCryptor().Verify(res.GetSignature().GetPublicKey(), block, res.GetSignature().GetSignature()))

			for _, tx := range txList.List() {
				stream.Req <- createRequestTx(tx)
			}
			stream.Err <- io.EOF
		}(t)
		err := ctrl.PropagateBlock(stream)
		require.NoError(t, err)
	})

	t.Run("case 2 : block is nil", func(t *testing.T) {
		stream := newMockPropagateBlockServerStream()
		go func(t *testing.T) {
			defer close(stream.Req)
			defer close(stream.Res)
			defer close(stream.Err)
			stream.Req <- nil
		}(t)
		err := ctrl.PropagateBlock(stream)
		require.Error(t, errors.Cause(err), codes.InvalidArgument)
	})

	t.Run("case 3 : tx is nil", func(t *testing.T) {
		stream := newMockPropagateBlockServerStream()
		block, txList := RandomValidSignedBlockAndTxList(t)
		go func(t *testing.T) {
			defer close(stream.Req)
			defer close(stream.Res)
			defer close(stream.Err)
			stream.Req <- createRequestBlock(block)

			res := <-stream.Res
			assert.NoError(t, RandomCryptor().Verify(res.GetSignature().GetPublicKey(), block, res.GetSignature().GetSignature()))

			for i, tx := range txList.List() {
				if i == 10 {
					stream.Req <- nil
					break
				}
				stream.Req <- createRequestTx(tx)
			}
		}(t)
		err := ctrl.PropagateBlock(stream)
		require.Error(t, errors.Cause(err), codes.Internal)
	})

	t.Run("case 4 : txs Hash not txList Hash", func(t *testing.T) {
		stream := newMockPropagateBlockServerStream()
		block, _ := RandomValidSignedBlockAndTxList(t)
		txList := RandomTxList()
		go func(t *testing.T) {
			defer close(stream.Req)
			defer close(stream.Res)
			defer close(stream.Err)
			stream.Req <- createRequestBlock(block)

			res := <-stream.Res
			assert.NoError(t, RandomCryptor().Verify(res.GetSignature().GetPublicKey(), block, res.GetSignature().GetSignature()))

			for _, tx := range txList.List() {
				stream.Req <- createRequestTx(tx)
			}
		}(t)
		err := ctrl.PropagateBlock(stream)
		require.Error(t, errors.Cause(err), codes.InvalidArgument)
	})
}
