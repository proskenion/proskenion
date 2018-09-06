package core

import "github.com/pkg/errors"

var (
	ErrDuplicateStore = errors.Errorf("Failed Duplicate Store")
	ErrNotFoundLoad = errors.Errorf("Failed Load not found")
)

type DBATx interface {
	Rollback() error
	Load(key Marshaler, value Unmarshaler) error // value = Load(key)
	Store(key Marshaler, value Marshaler) error  // Duplicate Insert error
	Commit() error
}

type DBA interface {
	Begin() DBATx
	Load(key Marshaler, value Unmarshaler) error // value = Load(key)
	Store(key Marshaler, value Marshaler) error  // Duplicate Insert error
}
