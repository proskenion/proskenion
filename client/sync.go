package client

import (
	"github.com/proskenion/proskenion/convertor"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
	"github.com/proskenion/proskenion/repository"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type SyncClient struct {
	proskenion.SyncClient
	fc model.ModelFactory
	c  core.Cryptor
}

func NewSyncClient(peer model.Peer, fc model.ModelFactory, c core.Cryptor) (core.SyncClient, error) {
	gc, err := grpc.Dial(peer.GetAddress(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &SyncClient{
		proskenion.NewSyncClient(gc),
		fc,
		c,
	}, nil
}

func (c *SyncClient) Sync(blockHash model.Hash, blockChan chan model.Block, txListChan chan core.TxList, errChan chan error) error {
	stream, err := c.SyncClient.Sync(context.TODO())
	if err != nil {
		return err
	}
	defer stream.CloseSend()

	for {
		if err := stream.Send(&proskenion.SyncRequest{BlockHash: blockHash}); err != nil {
			return err
		}

		for {
			res, err := stream.Recv()
			if err != nil {
				return err
			}
			block := res.GetBlock()
			if block == nil {
				break
			}
			modelBlock := c.fc.NewEmptyBlock()
			modelBlock.(*convertor.Block).Block = block
			blockHash = modelBlock.Hash()
			blockChan <- modelBlock

			txList := repository.NewTxList(c.c, c.fc)
			for {
				res, err := stream.Recv()
				if err != nil {
					return err
				}
				tx := res.GetTransaction()
				if tx == nil {
					break
				}
				modelTx := c.fc.NewEmptyTx()
				modelTx.(*convertor.Transaction).Transaction = tx
				if err := txList.Push(modelTx); err != nil {
					return err
				}
			}
			txListChan <- txList
		}
		select {
		case <-errChan:
			return err
		default:
		}
	}
}
