package client

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
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

func (c *SyncClient) Sync(blockHash model.Hash, blockChan chan model.Block, txListChan chan core.TxList) error {
	stream, err := c.SyncClient.Sync(context.TODO())
	if err != nil {
		return err
	}
	// TODO
	return stream.CloseSend()
}
