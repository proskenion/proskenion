package core

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core/model"
)

var (
	ErrDBADuplicateStore = errors.Errorf("Failed DBA Duplicate Store")
	ErrDBANotFoundLoad   = errors.Errorf("Failed DBA Load not found")
	ErrDBABeginErr       = errors.Errorf("Failed DBA BeignTx Error")
)

type KeyValueStore interface {
	Load(key model.Hash, value Unmarshaler) error // value = Load(key)
	Store(key model.Hash, value Marshaler) error  // Duplicate Insert error
}

type DB interface {
	DBA(table string) DBA
}

type DBATx interface {
	Rollback() error
	KeyValueStore
	Commit() error
}

type DBA interface {
	Begin() (DBATx, error)
	KeyValueStore
}
