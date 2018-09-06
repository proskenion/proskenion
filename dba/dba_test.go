package dba_test

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

type ErrUnmarshaler struct{}

func (a *ErrUnmarshaler) Unmarshal(b []byte) error {
	return errors.Errorf("Unmarshal error")
}

func testDBA_Store_Load(t *testing.T, dba core.DBA) {
	for _, c := range []struct {
		name     string
		key      core.Marshaler
		expValue core.Marshaler
		actValue core.Unmarshaler
		expErr   error
	}{
		{
			"case 1",
			RandomMarshaler(),
			RandomMarshaler(),
			RandomMarshaler(),
			nil,
		},
		{
			"case 2",
			RandomMarshaler(),
			RandomMarshaler(),
			RandomMarshaler(),
			nil,
		},
		{
			"case 3",
			RandomMarshaler(),
			RandomMarshaler(),
			RandomMarshaler(),
			nil,
		},
		{
			"failed unmarshal",
			RandomMarshaler(),
			RandomMarshaler(),
			&ErrUnmarshaler{},
			core.ErrUnmarshal,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			dba.Store(c.key, c.expValue)
			err := dba.Load(c.key, c.actValue)
			if c.expErr != nil {
				assert.EqualError(t, errors.Cause(err), c.expErr.Error())
				return
			} else {
				assert.NoError(t, err)
			}
			assert.EqualValues(t, c.expValue, c.actValue)

			err = dba.Store(c.key, c.expValue)
			assert.EqualError(t, errors.Cause(err), core.ErrDuplicateStore.Error())

			err = dba.Load(RandomMarshaler(), c.actValue)
			assert.EqualError(t, errors.Cause(err), core.ErrNotFoundLoad.Error())
		})
	}
}

func testDBATx_Store_Load(t *testing.T, dba core.DBA) {
	for _, c := range []struct {
		name     string
		key      core.Marshaler
		expValue core.Marshaler
		actValue core.Unmarshaler
		expErr   error
	}{
		{
			"case 1",
			RandomMarshaler(),
			RandomMarshaler(),
			RandomMarshaler(),
			nil,
		},
		{
			"case 2",
			RandomMarshaler(),
			RandomMarshaler(),
			RandomMarshaler(),
			nil,
		},
		{
			"case 3",
			RandomMarshaler(),
			RandomMarshaler(),
			RandomMarshaler(),
			nil,
		},
		{
			"failed unmarshal",
			RandomMarshaler(),
			RandomMarshaler(),
			&ErrUnmarshaler{},
			core.ErrUnmarshal,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			btx := dba.Begin()
			btx.Store(c.key, c.expValue)

			err := btx.Load(c.key, c.actValue)
			if c.expErr != nil {
				assert.EqualError(t, errors.Cause(err), c.expErr.Error())
				return
			} else {
				assert.NoError(t, err)
			}
			assert.EqualValues(t, c.expValue, c.actValue)

			err = dba.Load(c.key, c.actValue)
			assert.EqualError(t, errors.Cause(err), core.ErrNotFoundLoad.Error())

			err = btx.Store(c.key, c.expValue)
			assert.EqualError(t, errors.Cause(err), core.ErrDuplicateStore.Error())

			err = btx.Load(RandomMarshaler(), c.actValue)
			assert.EqualError(t, errors.Cause(err), core.ErrNotFoundLoad.Error())

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
