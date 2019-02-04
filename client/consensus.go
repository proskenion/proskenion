package client

import (
	"github.com/proskenion/proskenion/convertor"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type ConsensusClient struct {
	proskenion.ConsensusClient
	fc model.ModelFactory
	c  core.Cryptor
}

func NewConsensusClient(peer model.Peer, fc model.ModelFactory, c core.Cryptor) (core.ConsensusClient, error) {
	gc, err := grpc.Dial(peer.GetAddress(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &ConsensusClient{
		proskenion.NewConsensusClient(gc),
		fc,
		c,
	}, nil
}

func (c *ConsensusClient) PropagateTx(in model.Transaction) error {
	tx := in.(*convertor.Transaction).Transaction
	_, err := c.ConsensusClient.PropagateTx(context.TODO(), tx)
	return err
}

func (c *ConsensusClient) PropagateBlockStreamTx(block model.Block, txList core.TxList) error {
	stream, err := c.ConsensusClient.PropagateBlock(context.TODO())
	if err != nil {
		return err
	}
	defer stream.CloseSend()
	req := &proskenion.PropagateBlockRequest{
		Req: &proskenion.PropagateBlockRequest_Block{Block: block.(*convertor.Block).Block},
	}

	if err := stream.Send(req); err != nil {
		return err
	}

	// ack reply. (verify block)
	res, err := stream.Recv()
	if err != nil {
		return err
	}
	if err := c.c.Verify(res.GetSignature().GetPublicKey(), block, res.GetSignature().GetSignature()); err != nil {
		return err
	}

	for _, tx := range txList.List() {
		err := stream.Send(&proskenion.PropagateBlockRequest{
			Req: &proskenion.PropagateBlockRequest_Transaction{Transaction: tx.(*convertor.Transaction).Transaction}})
		if err != nil {
			return err
		}
	}
	return nil
}
