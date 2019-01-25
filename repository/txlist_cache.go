package repository

import "github.com/proskenion/proskenion/core"

type TxListCache struct {
	core.CacheMap
}

func (c *TxListCache) Set(txList core.TxList) error {
	return c.CacheMap.Set(txList)
}