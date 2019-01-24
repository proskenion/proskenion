package core

import (
	"github.com/pkg/errors"
)

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
	KeyValueStore
	Commit() error
}

type DBA interface {
	Begin() (DBATx, error)
	KeyValueStore
}
