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

func syncId(p model.Peer) string {
	return p.GetPeerId() + "#s"
}

func (c *ClientCache) SetConsensus(peer model.Peer, client core.ConsensusClient) error {
	return c.CacheMap.Set(&ClientHashWraper{conId(peer), client})
}

func (c *ClientCache) GetConsensus(p model.Peer) (core.ConsensusClient, bool) {
	ret, ok := c.Get(model.Hash(conId(p)))
	if !ok {
		return nil, false
	}
	cw, ok := ret.(*ClientHashWraper)
	if !ok {
		return nil, false
	}
	client, ok := cw.e.(core.ConsensusClient)
	if !ok {
		return nil, false
	}
	return client, true
}

func (c *ClientCache) SetAPI(peer model.Peer, client core.APIClient) error {
	return c.CacheMap.Set(&ClientHashWraper{apiId(peer), client})
}

func (c *ClientCache) GetAPI(p model.Peer) (core.APIClient, bool) {
	ret, ok := c.Get(model.Hash(apiId(p)))
	if !ok {
		return nil, false
	}
	cw, ok := ret.(*ClientHashWraper)
	if !ok {
		return nil, false
	}
	client, ok := cw.e.(core.APIClient)
	if !ok {
		return nil, false
	}
	return client, true
}

func (c *ClientCache) SetSync(peer model.Peer, client core.SyncClient) error {
	return c.CacheMap.Set(&ClientHashWraper{syncId(peer), client})
}

func (c *ClientCache) GetSync(p model.Peer) (core.SyncClient, bool) {
	ret, ok := c.Get(model.Hash(syncId(p)))
	if !ok {
		return nil, false
	}
	cw, ok := ret.(*ClientHashWraper)
	if !ok {
		return nil, false
	}
	client, ok := cw.e.(core.SyncClient)
	if !ok {
		return nil, false
	}
	return client, true
}
