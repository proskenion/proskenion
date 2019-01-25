package client

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/convertor"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
)

type ConsensusGateClient struct {
	proskenion.ConsensusGateClient
	fc model.ModelFactory
	c  core.Cryptor
}

func NewConsensusGateClient(peer model.Peer, fc model.ModelFactory, c core.Cryptor) (core.ConsensusGateClient, error) {
	gc, err := grpc.Dial(peer.GetAddress(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &ConsensusGateClient{
		proskenion.NewConsensusGateClient(gc),
		fc,
		c,
	}, nil
}

func (c *ConsensusGateClient) PropagateTx(in model.Transaction) error {
	tx := in.(*convertor.Transaction).Transaction
	_, err := c.ConsensusGateClient.PropagateTx(context.TODO(), tx)
	return err
}

func (c *ConsensusGateClient) PropagateBlockStreamTx(block model.Block, txList core.TxList) error {
	stream, err := c.ConsensusGateClient.PropagateBlock(context.TODO())
	defer stream.CloseSend()
	if err != nil {
		return err
	}
	req := &proskenion.PropagagateBlockRequest{
		Req: &proskenion.PropagagateBlockRequest_Block{Block: block.(*convertor.Block).Block},
	}

	if err := stream.Send(req); err != nil {
		return err
	}
	res, err := stream.Recv()
	if err != nil {
		return err
	}
	if err := c.c.Verify(res.GetSignature().GetPublicKey(), block, res.GetSignature().GetSignature()); err != nil {
		return err
	}

	for _, tx := range txList.List() {
		err := stream.Send(&proskenion.PropagagateBlockRequest{
			Req: &proskenion.PropagagateBlockRequest_Transaction{Transaction: tx.(*convertor.Transaction).Transaction}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ConsensusGateClient) CollectTx(blockHash model.Hash, txChan chan model.Transaction) error {
	stream, err := c.ConsensusGateClient.CollectTx(context.TODO(), &proskenion.CollectTxRequest{BlockHash: blockHash})
	if err != nil {
		return err
	}
	for {
		rtx, err := stream.Recv()
		if err == io.EOF { // サーバ側でストリーミングが正常に終了(return nil)された
			break
		}
		if err != nil {
			return errors.Errorf("%+v.CollectTx(_), error: %s\n", c.ConsensusGateClient, err.Error())
		}
		tx := c.fc.NewEmptyTx()
		tx.(*convertor.Transaction).Transaction = rtx
		txChan <- tx
	}
	return err
}
