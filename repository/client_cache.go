package repository

import (
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/datastructure"
)

type ClientCache struct {
	core.CacheMap
}

type ClientHashWraper struct {
	id string
	e  interface{}
}

func (w *ClientHashWraper) Hash() model.Hash {
	return model.Hash(w.id)
}

func NewClientCache(conf *config.Config) core.ClientCache {
	return &ClientCache{datastructure.NewCacheMap(conf.Cache.ClientLimits)}
}

func apiId(p model.Peer) string {
	return p.GetPeerId() + "#a"
}

func conId(p model.Peer) string {
	return p.GetPeerId() + "#c"
}

func (c *ClientCache) SetConsensus(peer model.Peer, client core.ConsensusGateClient) error {
	return c.CacheMap.Set(&ClientHashWraper{conId(peer), client})
}

func (c *ClientCache) GetConsensus(p model.Peer) (core.ConsensusGateClient, bool) {
	ret, ok := c.Get(model.Hash(conId(p)))
	if !ok {
		return nil, false
	}
	cw, ok := ret.(*ClientHashWraper)
	if !ok {
		return nil, false
	}
	client, ok := cw.e.(core.ConsensusGateClient)
	if !ok {
		return nil, false
	}
	return client, true
}

func (c *ClientCache) SetAPI(peer model.Peer, client core.APIGateClient) error {
	return c.CacheMap.Set(&ClientHashWraper{apiId(peer), client})
}

func (c *ClientCache) GetAPI(p model.Peer) (core.APIGateClient, bool) {
	ret, ok := c.Get(model.Hash(apiId(p)))
	if !ok {
		return nil, false
	}
	cw, ok := ret.(*ClientHashWraper)
	if !ok {
		return nil, false
	}
	client, ok := cw.e.(core.APIGateClient)
	if !ok {
		return nil, false
	}
	return client, true
}
