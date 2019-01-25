package repository

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type TxListCache struct {
	core.CacheMap
}

func (c *TxListCache) Set(txList core.TxList) error {
	return c.CacheMap.Set(txList)
}

func (c *TxListCache) Get(hash model.Hash) (core.TxList, bool) {
	ret, ok := c.CacheMap.Get(hash)
	if !ok {
		return nil, false
	}
	txList, ok := ret.(core.TxList)
	if !ok {
		return nil, false
	}
	return txList, true
}
