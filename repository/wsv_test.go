package repository_test

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/repository"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func test_WSV_Upserts(t *testing.T, wsv core.WSV, id model.Address, ac model.Account) {
	err := wsv.Query(id, ac)
	require.EqualError(t, errors.Cause(err), core.ErrWSVNotFound.Error())
	err = wsv.Append(id, ac)
	require.NoError(t, err)

	unmarshaler := RandomAccount()
	err = wsv.Query(id, unmarshaler)
	require.NoError(t, err)
	assert.Equal(t, MustHash(ac), MustHash(unmarshaler))
}

func test_WSV_Upserts_Peer(t *testing.T, wsv core.WSV, id model.Address, peer model.Peer) {
	err := wsv.Query(id, peer)
	require.EqualError(t, errors.Cause(err), core.ErrWSVNotFound.Error())
	err = wsv.Append(id, peer)
	require.NoError(t, err)

	unmarshaler := RandomPeer()
	err = wsv.Query(id, unmarshaler)
	require.NoError(t, err)
	assert.Equal(t, MustHash(peer), MustHash(unmarshaler))
}

func test_WSV(t *testing.T, wsv core.WSV) {
	acs := []model.Account{
		RandomAccount(),
		RandomAccount(),
		RandomAccount(),
		RandomAccount(),
		RandomAccount(),
	}
	ids := []model.Address{
		model.MustAddress("targeta@a"),
		model.MustAddress("tartb@a"),
		model.MustAddress("tartbc@a"),
		model.MustAddress("xyz@a"),
		model.MustAddress("target@a"),
	}

	for i, ac := range acs {
		test_WSV_Upserts(t, wsv, ids[i], ac)
	}
	peerRootAddress := model.MustAddress("/peer")
	_, err := wsv.PeerService(peerRootAddress)
	assert.Error(t, err)

	ps := []model.Peer{
		RandomPeer(),
		RandomPeer(),
		RandomPeer(),
		RandomPeer(),
	}
	pids := []model.Address{
		model.MustAddress("p1@peer/peer"),
		model.MustAddress("p2@peer/peer"),
		model.MustAddress("p3@peer/peer"),
		model.MustAddress("p4@peer/peer"),
	}

	for i, p := range ps {
		test_WSV_Upserts_Peer(t, wsv, pids[i], p)
	}
	require.NoError(t, wsv.Commit())

	peerService, err := wsv.PeerService(peerRootAddress)
	assert.NoError(t, err)
	assert.Equal(t, len(peerService.List()), 4)
	exPs := ps
	assertPs := func(p model.Peer) {
		for i, exP := range exPs {
			if bytes.Equal(MustHash(exP), MustHash(p)) {
				exPs = append(exPs[:i], exPs[i+1:]...)
				return
			}
		}
		assert.Failf(t, "assert peers.", "%x is not found.", MustHash(p))
	}
	for _, p := range peerService.List() {
		assertPs(p)
	}
	assert.Equal(t, len(exPs), 0)
}

func TestWSV(t *testing.T) {
	wsv, err := NewWSV(RandomDBATx(), RandomCryptor(), model.Hash(nil))
	require.NoError(t, err)
	test_WSV(t, wsv)
}
