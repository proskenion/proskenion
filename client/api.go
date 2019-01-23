package client

import (
	"context"
	"github.com/proskenion/proskenion/convertor"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
	"google.golang.org/grpc"
)

type APIGateClient struct {
	proskenion.APIGateClient
	fc model.ModelFactory
}

func NewAPIGateClient(peer model.Peer, fc model.ModelFactory) (core.APIGateClient, error) {
	gc, err := grpc.Dial(peer.GetAddress(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &APIGateClient{
		proskenion.NewAPIGateClient(gc),
		fc,
	}, nil
}

func (c *APIGateClient) Write(in model.Transaction) error {
	tx := in.(*convertor.Transaction).Transaction
	_, err := c.APIGateClient.Write(context.TODO(), tx)
	return err
}

func (c *APIGateClient) Read(in model.Query) (model.QueryResponse, error) {
	query := in.(*convertor.Query).Query
	res, err := c.APIGateClient.Read(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	qres := c.fc.NewEmptyQueryResponse()
	qres.(*convertor.QueryResponse).QueryResponse = res
	return qres, nil
}
