package core

import (
	"github.com/pkg/errors"
)

var (
	ErrDBADuplicateStore = errors.Errorf("Failed DBA Duplicate Store")
	ErrDBANotFoundLoad   = errors.Errorf("Failed DBA Load not found")
	ErrDBABeginErr       = errors.Errorf("Failed DBA BeignTx Error")
)

type Txx interface {
	Commit() error
	Rollback() error
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

func RollBackTx(tx Txx, mtErr error) error {
	if err := tx.Rollback(); err != nil {
		return errors.Wrap(err, mtErr.Error())
	}
	return mtErr
}

func CommitTx(tx Txx) error {
	if err := tx.Commit(); err != nil {
		return RollBackTx(tx, err)
	}
	return nil
}