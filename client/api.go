package client

import (
	"context"
	"github.com/proskenion/proskenion/convertor"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
	"google.golang.org/grpc"
)

type APIClient struct {
	proskenion.APIClient
	fc model.ModelFactory
}

func NewAPIClient(peer model.Peer, fc model.ModelFactory) (core.APIClient, error) {
	gc, err := grpc.Dial(peer.GetAddress(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &APIClient{
		proskenion.NewAPIClient(gc),
		fc,
	}, nil
}

func (c *APIClient) Write(in model.Transaction) error {
	tx := in.(*convertor.Transaction).Transaction
	_, err := c.APIClient.Write(context.TODO(), tx)
	return err
}

func (c *APIClient) Read(in model.Query) (model.QueryResponse, error) {
	query := in.(*convertor.Query).Query
	res, err := c.APIClient.Read(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	qres := c.fc.NewEmptyQueryResponse()
	qres.(*convertor.QueryResponse).QueryResponse = res
	return qres, nil
}
