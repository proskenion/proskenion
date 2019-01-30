package test_utils

import (
	. "github.com/proskenion/proskenion/core"
	. "github.com/proskenion/proskenion/core/model"
)

type MockClientFactory struct {
	cacheAPI map[string]APIClient
	cacheCon map[string]ConsensusClient
}

func NewMockClientFactory() ClientFactory {
	return &MockClientFactory{
		make(map[string]APIClient),
		make(map[string]ConsensusClient)}
}

func (f *MockClientFactory) APIClient(peer Peer) (APIClient, error) {
	if _, ok := f.cacheAPI[peer.GetPeerId()]; !ok {
		f.cacheAPI[peer.GetPeerId()] = &MockAPIClient{Id: peer.GetPeerId()}
	}
	return f.cacheAPI[peer.GetPeerId()], nil
}

func (f *MockClientFactory) ConsensusClient(peer Peer) (ConsensusClient, error) {
	if _, ok := f.cacheCon[peer.GetPeerId()]; !ok {
		f.cacheCon[peer.GetPeerId()] = &MockConsensusClient{Id: peer.GetPeerId()}
	}
	return f.cacheCon[peer.GetPeerId()], nil
}

type MockAPIClient struct {
	Id      string
	WriteIn Transaction
	ReadIn  Query
}

func (c *MockAPIClient) Write(in Transaction) error {
	c.WriteIn = in
	return nil
}
func (c *MockAPIClient) Read(in Query) (QueryResponse, error) {
	c.ReadIn = in
	return nil, nil
}

type MockConsensusClient struct {
	Id                string
	PropagateTxIn     Transaction
	PropagateBlockIn1 Block
	PropagateBlockIn2 TxList
}

func (c *MockConsensusClient) PropagateTx(tx Transaction) error {
	c.PropagateTxIn = tx
	return nil
}
func (c *MockConsensusClient) PropagateBlockStreamTx(block Block, txLit TxList) error {
	c.PropagateBlockIn1 = block
	c.PropagateBlockIn2 = txLit
	return nil
}
