package controller_test

import (
	"fmt"
	. "github.com/proskenion/proskenion/controller"
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

func newRandomConsensusServer() proskenion.ConsensusServer {
	fc, _, _, c, _, _, conf := NewTestFactories()
	cg := gate.NewConsensusGate(fc, c,
		repository.NewProposalTxQueueOnMemory(conf), repository.NewTxListCache(conf),
		repository.NewProposalBlockQueueOnMemory(conf), conf)
	return NewConsensusServer(fc, cg, c, RandomLogger(), RandomConfig())
}

type MockConsensus_PropagateBlockServer struct {
	Req    chan *proskenion.PropagateBlockRequest
	Res    chan *proskenion.PropagateBlockResponse
	Err    chan error
	failed bool
	grpc.ServerStream
}

func newMockPropagateBlockServerStream() *MockConsensus_PropagateBlockServer {
	return &MockConsensus_PropagateBlockServer{make(chan *proskenion.PropagateBlockRequest),
		make(chan *proskenion.PropagateBlockResponse),
		make(chan error),
		false,
		RandomMockServerStream()}
}

func (s *MockConsensus_PropagateBlockServer) Send(res *proskenion.PropagateBlockResponse) error {
	if s.failed {
		return fmt.Errorf("connection failed.")
	}
	s.Res <- res
	return nil
}

func (s *MockConsensus_PropagateBlockServer) Recv() (*proskenion.PropagateBlockRequest, error) {
	if s.failed {
		return nil, fmt.Errorf("connection failed.")
	}
	select {
	case req := <-s.Req:
		return req, nil
	case err := <-s.Err:
		return nil, err
	}
}

func (s *MockConsensus_PropagateBlockServer) _destructor() {
	s.failed = true
	close(s.Res)
	close(s.Req)
	close(s.Err)
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

func TestConsensusServer_PropagateBlock(t *testing.T) {
	ctrl := newRandomConsensusServer()

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
		statusCheck(t, err, codes.InvalidArgument)
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
		statusCheck(t, err, codes.Internal)
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
			stream.Err <- io.EOF
		}(t)
		err := ctrl.PropagateBlock(stream)
		statusCheck(t, err, codes.InvalidArgument)
	})
}
