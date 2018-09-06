package core

import "github.com/pkg/errors"

var (
	ErrDBADuplicateStore = errors.Errorf("Failed DBA Duplicate Store")
	ErrDBANotFoundLoad   = errors.Errorf("Failed DBA Load not found")
	ErrDBABeginErr       = errors.Errorf("Failed DBA BeignTx Error")
)

type DB interface {
	DBA(table string) DBA
}

type DBATx interface {
	Rollback() error
	Load(key Marshaler, value Unmarshaler) error // value = Load(key)
	Store(key Marshaler, value Marshaler) error  // Duplicate Insert error
	Commit() error
}

type DBA interface {
	Begin() (DBATx, error)
	Load(key Marshaler, value Unmarshaler) error // value = Load(key)
	Store(key Marshaler, value Marshaler) error  // Duplicate Insert error
}
