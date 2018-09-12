package dba_test

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/dba"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type ErrUnmarshaler struct{}

func (a *ErrUnmarshaler) Unmarshal(b []byte) error {
	return errors.Errorf("Unmarshal error")
}

func testDBA_Store_Load(t *testing.T, dba core.DBA) {
	for _, c := range []struct {
		name     string
		key      model.Hash
		expValue core.Marshaler
		actValue core.Unmarshaler
		expErr   error
	}{
		{
			"case 1",
			RandomByte(),
			RandomMarshaler(),
			RandomMarshaler(),
			nil,
		},
		{
			"case 2",
			RandomByte(),
			RandomMarshaler(),
			RandomMarshaler(),
			nil,
		},
		{
			"case 3",
			RandomByte(),
			RandomMarshaler(),
			RandomMarshaler(),
			nil,
		},
		{
			"failed unmarshal",
			RandomByte(),
			RandomMarshaler(),
			&ErrUnmarshaler{},
			core.ErrUnmarshal,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			err := dba.Store(c.key, c.expValue)
			require.NoError(t, err)
			err = dba.Load(c.key, c.actValue)
			if c.expErr != nil {
				assert.EqualError(t, errors.Cause(err), c.expErr.Error())
				return
			} else {
				require.NoError(t, err)
			}
			assert.EqualValues(t, c.expValue, c.actValue)

			err = dba.Store(c.key, c.expValue)
			assert.EqualError(t, errors.Cause(err), core.ErrDBADuplicateStore.Error())

			err = dba.Load(RandomByte(), c.actValue)
			assert.EqualError(t, errors.Cause(err), core.ErrDBANotFoundLoad.Error())
		})
	}
}

func testDBATx_Store_Load(t *testing.T, dba core.DBA) {
	for _, c := range []struct {
		name     string
		key      model.Hash
		expValue core.Marshaler
		actValue core.Unmarshaler
		expErr   error
	}{
		{
			"case 1",
			RandomByte(),
			RandomMarshaler(),
			RandomMarshaler(),
			nil,
		},
		{
			"case 2",
			RandomByte(),
			RandomMarshaler(),
			RandomMarshaler(),
			nil,
		},
		{
			"case 3",
			RandomByte(),
			RandomMarshaler(),
			RandomMarshaler(),
			nil,
		},
		{
			"failed unmarshal",
			RandomByte(),
			RandomMarshaler(),
			&ErrUnmarshaler{},
			core.ErrUnmarshal,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			btx, err := dba.Begin()
			require.NoError(t, err)

			defer btx.Rollback()

			err = btx.Store(c.key, c.expValue)
			require.NoError(t, err)

			err = btx.Load(c.key, c.actValue)
			if c.expErr != nil {
				assert.EqualError(t, errors.Cause(err), c.expErr.Error())
				return
			} else {
				require.NoError(t, err)
			}
			assert.EqualValues(t, c.expValue, c.actValue)

			if _, ok := dba.(*DBASQLite); !ok {
				// SQLite only supports one writer at a time per database file.
				// https://qiita.com/Yuki_312/items/7a7dff204e67af0c613a#sqlite3
				err = dba.Load(c.key, c.actValue)
				assert.EqualError(t, errors.Cause(err), core.ErrDBANotFoundLoad.Error())
			}

			err = btx.Store(c.key, c.expValue)
			assert.EqualError(t, errors.Cause(err), core.ErrDBADuplicateStore.Error())

			err = btx.Load(RandomByte(), c.actValue)
			assert.EqualError(t, errors.Cause(err), core.ErrDBANotFoundLoad.Error())

			err = btx.Commit()
			assert.NoError(t, err)

			dba.Load(c.key, c.actValue)
			if c.expErr != nil {
				assert.EqualError(t, errors.Cause(err), c.expErr.Error())
				return
			} else {
				assert.NoError(t, err)
			}
			assert.EqualValues(t, c.expValue, c.actValue)
		})
	}
}
