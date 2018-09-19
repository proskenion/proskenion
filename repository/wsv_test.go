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

func test_WSV_Upserts(t *testing.T, wsv core.WSV, id string, ac model.Account) {
	err := wsv.Query(id, ac)
	require.EqualError(t, errors.Cause(err), core.ErrWSVNotFound.Error())
	err = wsv.Append(id, ac)
	require.NoError(t, err)

	unmarshaler := RandomAccount()
	err = wsv.Query(id, unmarshaler)
	require.NoError(t, err)
	assert.Equal(t, MustHash(ac), MustHash(unmarshaler.(model.Account)))
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
		fmt.Println("==upserts", i)
		test_WSV_Upserts(t, wsv, ids[i], ac)
	}
	require.NoError(t, wsv.Commit())
}

func TestWSV(t *testing.T) {
	wsv, err := NewWSV(RandomDBATx(), RandomCryptor(), model.Hash(nil))
	require.NoError(t, err)
	test_WSV(t, wsv)
}
