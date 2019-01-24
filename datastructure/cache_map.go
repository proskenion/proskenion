package datastructure

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"sync"
)

type HasherAndInd struct {
	model.Hasher
	int
}

const CacheStageNum = 3

// LRU 簡易実装 3-段キャッシュ
type CacheMap struct {
	mutex    *sync.Mutex
	limit    int
	oldestId int
	field    map[string]*HasherAndInd
	index    []map[string]struct{}
}

func NewCacheMap(limit int) core.CacheMap {
	return &CacheMap{new(sync.Mutex), limit, 0, make(map[string]*HasherAndInd), make([]map[string]struct{}, 3)}
}

func nextCacheStage(oldestId int) int {
	return (oldestId + 1) % CacheStageNum
}
func recentCacheStage(oldestId int) int {
	return (oldestId + CacheStageNum - 1) % CacheStageNum
}

func (c *CacheMap) Set(hasher model.Hasher) error {
	if hasher == nil {
		return core.ErrCacheMapPushNil
	}
	keyHash := string(hasher.Hash())
	if len(c.field) >= c.limit {
		for len(c.index[c.oldestId]) == 0 {
			c.oldestId = nextCacheStage(c.oldestId)
		}
		for k, _ := range c.index[c.oldestId] {
			delete(c.field, k)
			delete(c.index[c.oldestId], k)
			break
		}
	}
	c.field[keyHash] = &HasherAndInd{hasher, recentCacheStage(c.oldestId)}
	c.index[recentCacheStage(c.oldestId)][keyHash] = struct{}{}
	return nil
}

func (c *CacheMap) Get(hash model.Hash) (model.Hasher, bool) {
	keyHash := string(hash)
	ret, ok := c.field[keyHash]
	if !ok {
		return nil, false
	}
	// recently use
	delete(c.index[ret.int], keyHash)
	if len(c.index[c.oldestId]) == 0 {
		c.oldestId = nextCacheStage(c.oldestId)
	}
	c.index[recentCacheStage(c.oldestId)][keyHash] = struct{}{}
	return ret, ok
}
