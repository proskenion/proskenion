package controller

import (
	"github.com/inconshreveable/log15"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/convertor"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
	"io"
)

// SyncServer is the server Sync for SyncGate service.
type SyncServer struct {
	fc     model.ModelFactory
	sg     core.SyncGate
	c      core.Cryptor
	logger log15.Logger

	conf *config.Config
}

func NewSyncServer(fc model.ModelFactory, sg core.SyncGate, c core.Cryptor, logger log15.Logger, conf *config.Config) proskenion.SyncServer {
	return &SyncServer{
		fc,
		sg,
		c,
		logger,
		conf,
	}
}

func (s *SyncServer) newBlockResponse(block model.Block) *proskenion.SyncResponse {
	return &proskenion.SyncResponse{
		Res: &proskenion.SyncResponse_Block{Block: block.(*convertor.Block).Block},
	}
}

func (s *SyncServer) newTxResponse(tx model.Transaction) *proskenion.SyncResponse {
	return &proskenion.SyncResponse{
		Res: &proskenion.SyncResponse_Transaction{Transaction: tx.(*convertor.Transaction).Transaction},
	}
}

func (s *SyncServer) Sync(stream proskenion.Sync_SyncServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		blockHash := req.GetBlockHash()
		blockChan := make(chan model.Block)
		txListChan := make(chan core.TxList)
		errChan := make(chan error)
		go func(blockHash model.Hash) {
			defer close(blockChan)
			defer close(txListChan)
			defer close(errChan)
			err := s.sg.Sync(blockHash, blockChan, txListChan)
			errChan <- err
		}(blockHash)
		for {
			select {
			case newBlock := <-blockChan:
				if err := stream.Send(s.newBlockResponse(newBlock)); err != nil {
					return err
				}
			case newTxList := <-txListChan:
				for _, tx := range newTxList.List() {
					if err := stream.Send(s.newTxResponse(tx)); err != nil {
						return err
					}
				}
				if err := stream.Send(&proskenion.SyncResponse{}); err != nil {
					return err
				}
			case err := <-errChan:
				if err != nil && err != io.EOF {
					return err
				}
				goto afterFor
			}
		}
	afterFor:
	}
}
