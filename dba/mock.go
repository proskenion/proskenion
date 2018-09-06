package dba

import (
	"encoding/hex"
	"github.com/pkg/errors"
	. "github.com/proskenion/proskenion/core"
	"sync"
)

type syncMapApplyBytes struct {
	*sync.Map
}

func newSyncMapApplyBytes() *syncMapApplyBytes {
	return &syncMapApplyBytes{new(sync.Map)}
}

func (s *syncMapApplyBytes) Store(key, value interface{}) {
	if b, ok := key.([]byte); ok {
		s.Map.Store(hex.EncodeToString(b), value)
		return
	}
	s.Map.Store(key, value)
}

func (s *syncMapApplyBytes) Load(key interface{}) (interface{}, bool) {
	if b, ok := key.([]byte); ok {
		return s.Map.Load(hex.EncodeToString(b))
	}
	return s.Map.Load(key)
}

type DBAOnMemory struct {
	db *syncMapApplyBytes
}

func (d *DBAOnMemory) Begin() DBATx {
	return &DBAOnMemoryTx{d.db, newSyncMapApplyBytes()}
}

func (d *DBAOnMemory) Load(key Marshaler, value Unmarshaler) error {
	tx := d.Begin()
	return tx.Load(key, value)
}

func (d *DBAOnMemory) Store(key Marshaler, value Marshaler) error {
	tx := d.Begin()
	if err := tx.Store(key, value); err != nil {
		return err
	}
	return tx.Commit()
}

func NewDBAOnMemory() DBA {
	return &DBAOnMemory{newSyncMapApplyBytes()}
}

type DBAOnMemoryTx struct {
	origin *syncMapApplyBytes
	tmp    *syncMapApplyBytes
}

func (t *DBAOnMemoryTx) Rollback() error {
	return nil
}

func (t *DBAOnMemoryTx) castAndUnmarshal(v interface{}, value Unmarshaler) error {
	b, ok := v.([]byte)
	if !ok {
		return errors.Errorf("Unexpected Error can not cast byte %v", v)
	}
	err := value.Unmarshal(b)
	if err != nil {
		return errors.Wrapf(ErrUnmarshal, err.Error())
	}
	return nil
}

func (t *DBAOnMemoryTx) Load(key Marshaler, value Unmarshaler) error {
	k, err := key.Marshal()
	if err != nil {
		return errors.Wrap(ErrMarshal, err.Error())
	}
	if v, ok := t.origin.Load(k); ok {
		return t.castAndUnmarshal(v, value)
	}
	if v, ok := t.tmp.Load(k); ok {
		return t.castAndUnmarshal(v, value)
	}
	return errors.Wrapf(ErrNotFoundLoad, hex.EncodeToString(k))
}

func (t *DBAOnMemoryTx) checkDuplicate(key []byte) error {
	if _, ok := t.origin.Load(key); ok {
		return errors.Wrap(ErrDuplicateStore, hex.EncodeToString(key))
	}
	if _, ok := t.tmp.Load(key); ok {
		return errors.Wrap(ErrDuplicateStore, hex.EncodeToString(key))
	}
	return nil
}

func (t *DBAOnMemoryTx) Store(key Marshaler, value Marshaler) error {
	k, err := key.Marshal()
	if err != nil {
		return errors.Wrap(ErrMarshal, err.Error())
	}
	if err = t.checkDuplicate(k); err != nil {
		return err
	}
	v, err := value.Marshal()
	if err != nil {
		return errors.Wrap(ErrMarshal, err.Error())
	}
	t.tmp.Store(k, v)
	return nil
}

func (t *DBAOnMemoryTx) Commit() error {
	t.tmp.Range(func(key, value interface{}) bool {
		t.origin.Store(key, value)
		return true
	})
	return nil
}
