package client

import (
	"context"
	"github.com/proskenion/proskenion/convertor"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
	"google.golang.org/grpc"
)

type APIGateClient struct {
	proskenion.APIGateClient
	publicKey  model.PublicKey
	privateKey model.PrivateKey
	fc         model.ModelFactory
}

func (c *APIGateClient) Write(in model.Transaction, opts ...grpc.CallOption) error {
	tx := in.(*convertor.Transaction).Transaction
	_, err := c.APIGateClient.Write(context.TODO(), tx)
	return err
}

func (c *APIGateClient) Read(in model.Query, opts ...grpc.CallOption) (model.QueryResponse, error) {
	query := in.(*convertor.Query).Query
	res, err := c.APIGateClient.Read(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	qres := c.fc.NewEmptyQueryResponse()
	qres.(*convertor.QueryResponse).QueryResponse = res
	return qres, nil
}
