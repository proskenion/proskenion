package client

import (
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/repository"
)

type ClientFactory struct {
	fc    model.ModelFactory
	c     core.Cryptor
	cache core.ClientCache
}

func NewClientFactory(fc model.ModelFactory, c core.Cryptor, conf *config.Config) core.ClientFactory {
	return &ClientFactory{fc, c, repository.NewClientCache(conf)}
}

func (fc *ClientFactory) APIClient(peer model.Peer) (core.APIGateClient, error) {
	ret, ok := fc.cache.GetAPI(peer)
	if ok {
		return ret, nil
	}
	ret, err := NewAPIGateClient(peer, fc.fc)
	if err != nil {
		return nil, err
	}
	if err := fc.cache.SetAPI(peer, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (fc *ClientFactory) ConsensusClient(peer model.Peer) (core.ConsensusGateClient, error) {
	ret, ok := fc.cache.GetConsensus(peer)
	if ok {
		return ret, nil
	}
	ret, err := NewConsensusGateClient(peer, fc.fc, fc.c)
	if err != nil {
		return nil, err
	}
	if err := fc.cache.SetConsensus(peer, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
