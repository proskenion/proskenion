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

func test_WSV_Upserts(t *testing.T, wsv core.WSV, id string, ac model.Account) {
	err := wsv.Query(id, ac)
	require.EqualError(t, errors.Cause(err), core.ErrWSVNotFound.Error())
	err = wsv.Append(id, ac)
	require.NoError(t, err)

	unmarshaler := RandomAccount()
	err = wsv.Query(id, unmarshaler)
	require.NoError(t, err)
	assert.Equal(t, MustHash(ac), MustHash(unmarshaler))
}

func test_WSV_Upserts_Peer(t *testing.T, wsv core.WSV, id string, peer model.Peer) {
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
	ids := []string{
		"targeta",
		"tartb",
		"tartbc",
		"xyz",
		"target",
	}

	for i, ac := range acs {
		test_WSV_Upserts(t, wsv, ids[i], ac)
	}
	_, err := wsv.PeerService()
	assert.Error(t, err)

	ps := []model.Peer{
		RandomPeer(),
		RandomPeer(),
		RandomPeer(),
		RandomPeer(),
	}
	pids := []string{
		":127.32.2.1",
		":127.32.2.2",
		":127.32.2.3",
		":127.32.2.4",
	}

	for i, p := range ps {
		test_WSV_Upserts_Peer(t, wsv, pids[i], p)
	}
	require.NoError(t, wsv.Commit())

	peerService, err := wsv.PeerService()
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
