package dba

import (
	"encoding/hex"
	"github.com/pkg/errors"
	. "github.com/proskenion/proskenion/core"
	. "github.com/proskenion/proskenion/core/model"
	"sync"
)

type DBOnMemory struct {
	dba map[string]DBA
}

func NewDBOnMemory() DB {
	return &DBOnMemory{make(map[string]DBA)}
}

func (db *DBOnMemory) DBA(table string) DBA {
	if _, ok := db.dba[table]; !ok {
		db.dba[table] = NewDBAOnMemory()
	}
	return db.dba[table]
}

type syncMapApplyBytes struct {
	*sync.Map
}

func newSyncMapApplyBytes() *syncMapApplyBytes {
	return &syncMapApplyBytes{new(sync.Map)}
}

func (s *syncMapApplyBytes) Store(key, value interface{}) {
	if b, ok := key.(Hash); ok {
		s.Map.Store(hex.EncodeToString(b), value)
		return
	}
	s.Map.Store(key, value)
}

func (s *syncMapApplyBytes) Load(key interface{}) (interface{}, bool) {
	if b, ok := key.(Hash); ok {
		return s.Map.Load(hex.EncodeToString(b))
	}
	return s.Map.Load(key)
}

type DBAOnMemory struct {
	db *syncMapApplyBytes
}

func (d *DBAOnMemory) Begin() (DBATx, error) {
	return &DBAOnMemoryTx{d.db, newSyncMapApplyBytes()}, nil
}

func (d *DBAOnMemory) Load(key Hash, value Unmarshaler) error {
	tx, _ := d.Begin()
	return tx.Load(key, value)
}

func (d *DBAOnMemory) Store(key Hash, value Marshaler) error {
	tx, _ := d.Begin()
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

func (t *DBAOnMemoryTx) Load(key Hash, value Unmarshaler) error {
	if v, ok := t.origin.Load(key); ok {
		return t.castAndUnmarshal(v, value)
	}
	if v, ok := t.tmp.Load(key); ok {
		return t.castAndUnmarshal(v, value)
	}
	return errors.Wrapf(ErrDBANotFoundLoad, hex.EncodeToString(key))
}

func (t *DBAOnMemoryTx) checkDuplicate(key Hash) error {
	if _, ok := t.origin.Load(key); ok {
		return errors.Wrap(ErrDBADuplicateStore, hex.EncodeToString(key))
	}
	if _, ok := t.tmp.Load(key); ok {
		return errors.Wrap(ErrDBADuplicateStore, hex.EncodeToString(key))
	}
	return nil
}

func (t *DBAOnMemoryTx) Store(key Hash, value Marshaler) error {
	if err := t.checkDuplicate(key); err != nil {
		return err
	}
	v, err := value.Marshal()
	if err != nil {
		return errors.Wrap(ErrMarshal, err.Error())
	}
	t.tmp.Store(key, v)
	return nil
}

func (t *DBAOnMemoryTx) Commit() error {
	t.tmp.Range(func(key, value interface{}) bool {
		t.origin.Store(key, value)
		return true
	})
	return nil
}
