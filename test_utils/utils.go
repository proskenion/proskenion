package test_utils

import (
	"bytes"
	"github.com/proskenion/proskenion/core/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func CastedHasher(o model.Hasher) model.Hasher {
	return o
}

func CastedHasherPeerList(o []model.Peer) []model.Hasher {
	ar := make([]model.Hasher, 0)
	for _, c := range o {
		ar = append(ar, CastedHasher(c))
	}
	return ar
}

func CastedHasherAccountList(o []model.Account) []model.Hasher {
	ar := make([]model.Hasher, 0)
	for _, c := range o {
		ar = append(ar, CastedHasher(c))
	}
	return ar
}

func CastedHasherStorageList(o []model.Storage) []model.Hasher {
	ar := make([]model.Hasher, 0)
	for _, c := range o {
		ar = append(ar, CastedHasher(c))
	}
	return ar
}

func AssertSetEqual(t *testing.T, actI interface{}, expI interface{}) {
	act := make([]model.Hasher, 0)
	exp := make([]model.Hasher, 0)

	switch sl := actI.(type) {
	case []model.Peer:
		act = CastedHasherPeerList(sl)
	case []model.Account:
		act = CastedHasherAccountList(sl)
	case []model.Storage:
		act = CastedHasherStorageList(sl)
	}

	switch sl := expI.(type) {
	case []model.Peer:
		exp = CastedHasherPeerList(sl)
	case []model.Account:
		exp = CastedHasherAccountList(sl)
	case []model.Storage:
		exp = CastedHasherStorageList(sl)
	}

	assert.Equal(t, len(act), len(exp))
	for i, ac := range act {
		for i, ex := range exp {
			if bytes.Equal(ac.Hash(), ex.Hash()) {
				exp = append(exp[:i], exp[i+1:]...)
				break
			}
		}
		assert.Equalf(t, len(act)-i-1, len(exp), "assert hasher.", "%x is not found.", ac.Hash())
	}
}
