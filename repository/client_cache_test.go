package repository_test

import (
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/repository"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClientCache_GetSetAPI(t *testing.T) {
	conf := RandomConfig()
	cc := NewClientCache(conf)
	ps := []model.Peer{
		RandomPeer(),
		RandomPeer(),
		RandomPeer(),
		RandomPeer(),
		RandomPeer(),
		RandomPeer(),
	}
	t.Run("case 1 : no error", func(t *testing.T) {
		for _, p := range ps {
			err := cc.SetAPI(p, &MockAPIClient{Id: p.GetPeerId()})
			require.NoError(t, err)
		}
		for _, p := range ps {
			ret, ok := cc.GetAPI(p)
			require.True(t, ok)
			assert.Equal(t, p.GetPeerId(), ret.(*MockAPIClient).Id)
		}

	})
	t.Run("case 2 : failed get", func(t *testing.T) {
		_, ok := cc.GetAPI(RandomPeer())
		assert.False(t, ok)
	})
}

func TestClientCache_GetSetConsensus(t *testing.T) {
	conf := RandomConfig()
	cc := NewClientCache(conf)
	ps := []model.Peer{
		RandomPeer(),
		RandomPeer(),
		RandomPeer(),
		RandomPeer(),
		RandomPeer(),
		RandomPeer(),
	}
	t.Run("case 1 : no error", func(t *testing.T) {
		for _, p := range ps {
			err := cc.SetConsensus(p, &MockConsensusClient{Id: p.GetPeerId()})
			require.NoError(t, err)
		}
		for _, p := range ps {
			ret, ok := cc.GetConsensus(p)
			assert.True(t, ok)
			assert.Equal(t, p.GetPeerId(), ret.(*MockConsensusClient).Id)
		}

	})
	t.Run("case 2 : failed get", func(t *testing.T) {
		_, ok := cc.GetConsensus(RandomPeer())
		assert.False(t, ok)
	})

}
