package datastructure_test

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/datastructure"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCacheMap_GetSet(t *testing.T) {
	hashers := []model.Hasher{
		RandomAccount(),
		RandomAccount(),
		RandomAccount(),
		RandomAccount(),
		RandomAccount(),
		RandomAccount(),
		RandomAccount(),
		RandomAccount(),
		RandomAccount(),
		RandomAccount(),
	}

	t.Run("case nil set error", func(t *testing.T) {
		cm := NewCacheMap(10)
		err := cm.Set(nil)
		assert.EqualError(t, errors.Cause(err), core.ErrCacheMapPushNil.Error())
	})

	t.Run("case 1 set and get", func(t *testing.T) {
		cm := NewCacheMap(10)
		for _, hasher := range hashers {
			require.NoError(t, cm.Set(hasher))
		}
		for i, hasher := range hashers {
			if i%2 == 0 {
				ret, ok := cm.Get(hasher.Hash())
				assert.True(t, ok)
				assert.Equal(t, hasher.Hash(), ret.Hash())
			}
		}
		for i := 0; i < 5; i++ {
			require.NoError(t, cm.Set(RandomAccount()))
		}
		for i, hasher := range hashers {
			if i%2 == 0 {
				ret, ok := cm.Get(hasher.Hash())
				assert.True(t, ok)
				assert.Equal(t, hasher.Hash(), ret.Hash())
			} else { // recently get hasher is deleted
				_, ok := cm.Get(hasher.Hash())
				assert.False(t, ok)
			}
		}
	})
}
