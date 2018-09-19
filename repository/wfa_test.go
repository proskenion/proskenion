package repository_test

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/repository"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func test_WFA_Upserts(t *testing.T, wfa core.WFA, id string, ac model.Account) {
	err := wfa.Query(id, ac)
	require.EqualError(t, errors.Cause(err), core.ErrWFANotFound.Error())
	err = wfa.Append(id, ac)
	require.NoError(t, err)

	unmarshaler := RandomAccount()
	err = wfa.Query(id, unmarshaler)
	require.NoError(t, err)
	assert.Equal(t, MustHash(ac), MustHash(unmarshaler.(model.Account)))
}

func test_WFA(t *testing.T, wfa core.WFA) {
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
		fmt.Println("==upserts", i)
		test_WFA_Upserts(t, wfa, ids[i], ac)
	}
	require.NoError(t, wfa.Commit())
}

func TestWFA(t *testing.T) {
	wfa, err := NewWFA(RandomDBATx(), RandomCryptor(), model.Hash(nil))
	require.NoError(t, err)
	test_WFA(t, wfa)
}
