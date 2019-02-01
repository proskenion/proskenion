package controller_test

import (
	. "github.com/proskenion/proskenion/controller"
	"github.com/proskenion/proskenion/convertor"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/gate"
	"github.com/proskenion/proskenion/proto"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"io"
	"testing"
)

type MockSync_SyncServer struct {
	Req chan *proskenion.SyncRequest
	Res chan *proskenion.SyncResponse
	Err chan error
	grpc.ServerStream
}

func newRandomSyncServer(t *testing.T) (proskenion.SyncServer, core.Repository) {
	fc, _, _, c, rp, _, conf := NewTestFactories()
	require.NoError(t, rp.GenesisCommit(RandomGenesisTxList(t)))
	for i := 0; i < conf.Sync.Limits*2; i++ {
		RandomCommitableBlockAndTxList(t, rp)
	}
	sg := gate.NewSyncGate(rp, fc, c, conf)
	return NewSyncServer(fc, sg, c, RandomLogger(), conf), rp
}

func newMockSyncServerStream() *MockSync_SyncServer {
	return &MockSync_SyncServer{make(chan *proskenion.SyncRequest),
		make(chan *proskenion.SyncResponse),
		make(chan error),
		RandomMockServerStream()}
}

func (s *MockSync_SyncServer) Send(res *proskenion.SyncResponse) error {
	select {
	case err := <-s.Err:
		return err
	default:
		break
	}
	s.Res <- res
	return nil
}

func (s *MockSync_SyncServer) Recv() (*proskenion.SyncRequest, error) {
	select {
	case req := <-s.Req:
		return req, nil
	case err := <-s.Err:
		return nil, err
	}
}

func TestNewSyncServer(t *testing.T) {
	ctrl, mrp := newRandomSyncServer(t)
	t.Run("case 1 : correct", func(t *testing.T) {
		rp := RandomRepository()
		require.NoError(t, rp.GenesisCommit(RandomGenesisTxList(t)))
		fc := RandomFactory()
		limits := RandomConfig().Sync.Limits
		stream := newMockSyncServerStream()
		go func(t *testing.T, limits int) {
			defer close(stream.Req)
			defer close(stream.Res)
			defer close(stream.Err)

			for i := 0; i < limits*2; i += limits {
				stream.Req <- &proskenion.SyncRequest{BlockHash: MusTop(rp).Hash()}
				for j := 0; j < limits; j++ {
					res := <-stream.Res
					block := res.GetBlock()
					require.NotNil(t, block)
					modelBlock := fc.NewEmptyBlock()
					modelBlock.(*convertor.Block).Block = block

					txList := EmptyTxList()
					for {
						res := <-stream.Res
						if res.Res == nil {
							break
						}
						tx := res.GetTransaction()
						require.NotNil(t, tx)
						modelTx := fc.NewEmptyTx()
						modelTx.(*convertor.Transaction).Transaction = tx
						require.NoError(t, txList.Push(modelTx))
					}
					require.NoError(t, rp.Commit(modelBlock, txList))
				}
				res := <-stream.Res
				require.Nil(t, res.Res)
			}
			stream.Err <- io.EOF
		}(t, limits)
		err := ctrl.Sync(stream)
		require.NoError(t, err)
		assert.Equal(t, MusTop(mrp).Hash(), MusTop(rp).Hash())
	})

	t.Run("case 2 : send first nil hash", func(t *testing.T) {
		rp := RandomRepository()
		fc := RandomFactory()
		require.NoError(t, rp.GenesisCommit(RandomGenesisTxList(t)))
		limits := RandomConfig().Sync.Limits
		stream := newMockSyncServerStream()
		go func(t *testing.T, limits int) {
			defer close(stream.Req)
			defer close(stream.Res)
			defer close(stream.Err)

			stream.Req <- &proskenion.SyncRequest{BlockHash: nil}
			res := <-stream.Res
			block := res.GetBlock()
			require.NotNil(t, block)
			modelBlock := fc.NewEmptyBlock()
			modelBlock.(*convertor.Block).Block = block

			txList := EmptyTxList()
			for {
				res := <-stream.Res
				if res.Res == nil {
					break
				}
				tx := res.GetTransaction()
				require.NotNil(t, tx)
				modelTx := fc.NewEmptyTx()
				modelTx.(*convertor.Transaction).Transaction = tx
				require.NoError(t, txList.Push(modelTx))
			}
			err := rp.Commit(modelBlock, txList)
			assert.Error(t, err)
			stream.Err <- err
		}(t, limits)
		err := ctrl.Sync(stream)
		statusCheck(t, err, codes.Internal)
	})
}
