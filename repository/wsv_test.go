package repository_test

import (
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

type AccountUnmarshalerFactory struct {
	fc model.ModelFactory
}

func (f *AccountUnmarshalerFactory) CreateUnmarshaler() model.Unmarshaler {
	return f.fc.NewEmptyAccount()
}

func NewAccountUnmarshalerFactory() *AccountUnmarshalerFactory {
	return &AccountUnmarshalerFactory{NewTestFactory()}
}

func testWSV_QueryAll(t *testing.T, wsv core.WSV, acs []model.Account, id model.Address) {
	res, err := wsv.QueryAll(id, NewAccountUnmarshalerFactory())
	assert.NoError(t, err)
	resAc := make([]model.Account, 0)
	for _, xxx := range res {
		resAc = append(resAc, xxx.(model.Account))
	}
	AssertSetEqual(t, resAc, acs)
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
		model.MustAddress("targeta@a/account"),
		model.MustAddress("tartb@a/account"),
		model.MustAddress("tartbc@a/account"),
		model.MustAddress("xyz@b/account"),
		model.MustAddress("target@b/account"),
	}

	for i, ac := range acs {
		test_WSV_Upserts(t, wsv, ids[i], ac)
	}
	// Queryll test
	testWSV_QueryAll(t, wsv, acs, model.MustAddress("/account"))
	testWSV_QueryAll(t, wsv, acs[3:], model.MustAddress("b/account"))

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
	AssertSetEqual(t, peerService.List(), ps)
}

func TestWSV(t *testing.T) {
	wsv, err := NewWSV(RandomDBATx(), RandomCryptor(), model.Hash(nil))
	require.NoError(t, err)
	test_WSV(t, wsv)
}
