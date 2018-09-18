package dba

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	sqlite "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/config"
	. "github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type DBSQLite struct {
	dba map[string]DBA
	db  *sqlx.DB
}

func newSQLite(conf *config.Config) DB {
	db, err := sqlx.Open(conf.DB.Kind, fmt.Sprintf("file:%s/%s.sqlite?cache=shared", conf.DB.Path, conf.DB.Name))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return &DBSQLite{make(map[string]DBA), db}
}

func NewDBSQLite(conf *config.Config) DB {
	fmt.Println(conf.DB.Kind)
	switch conf.DB.Kind {
	case "sqlite3":
		return newSQLite(conf)
	}
	return nil
}

func (db *DBSQLite) DBA(table string) DBA {
	if _, ok := db.dba[table]; !ok {
		db.dba[table] = NewDBA(db.db, table)
	}
	return db.dba[table]
}

func sq() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

type DBASQLite struct {
	db    *sqlx.DB
	table string
}

func commitTx(dba DBATx) error {
	if err := dba.Commit(); err != nil {
		return rollBackTx(dba, err)
	}
	return nil
}

func rollBackTx(dba DBATx, mtErr error) error {
	if err := dba.Rollback(); err != nil {
		return errors.Wrap(err, mtErr.Error())
	}
	return mtErr
}

func (d *DBASQLite) Begin() (DBATx, error) {
	tx, err := d.db.Beginx()
	if err != nil {
		return nil, err
	}
	return &DBASQLiteTx{tx, d.table}, nil
}

func (d *DBASQLite) Load(key model.Hash, value Unmarshaler) error {
	tx, err := d.Begin()
	if err != nil {
		return rollBackTx(tx, errors.Wrap(ErrDBABeginErr, err.Error()))
	}
	return rollBackTx(tx, tx.Load(key, value))
}

func (d *DBASQLite) Store(key model.Hash, value Marshaler) error {
	tx, err := d.Begin()
	if err != nil {
		return rollBackTx(tx, errors.Wrap(ErrDBABeginErr, err.Error()))
	}
	if err := tx.Store(key, value); err != nil {
		return rollBackTx(tx, err)
	}
	return commitTx(tx)
}

func NewDBA(db *sqlx.DB, table string) DBA {
	schema := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s ("key" BLOB PRIMARY KEY, "value" BLOB);`, table)
	db.MustExec(schema)
	return &DBASQLite{db, table}
}

type DBASQLiteTx struct {
	*sqlx.Tx
	table string
}

type KVTable struct {
	Key   []byte `db:"key"`
	Value []byte `db:"value"`
}

func (t *DBASQLiteTx) Rollback() error {
	return t.Tx.Rollback()
}

func (t *DBASQLiteTx) loadAndCast(k []byte, v Unmarshaler) error {
	query, args, err := sq().Select("*").
		From(t.table).
		Where(squirrel.Eq{"key": k}).
		ToSql()
	if err != nil {
		return err
	}
	kvTable := KVTable{}
	if err = t.Get(&kvTable, query, args...); err != nil {
		return err
	}
	if err = v.Unmarshal(kvTable.Value); err != nil {
		return errors.Wrap(ErrUnmarshal, err.Error())
	}
	return nil
}

func (t *DBASQLiteTx) store(k []byte, v []byte) error {
	_, err := sq().Insert(t.table).
		Columns("key", "value").
		Values(k, v).
		RunWith(t.Tx.Tx).Exec()
	return err
}

func (t *DBASQLiteTx) Load(key model.Hash, value Unmarshaler) error {
	if err := t.loadAndCast(key, value); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return errors.Wrap(ErrDBANotFoundLoad, err.Error())
		}
		return err
	}
	return nil
}

func (t *DBASQLiteTx) Store(key model.Hash, value Marshaler) error {
	v, err := value.Marshal()
	if err != nil {
		return errors.Wrap(ErrMarshal, err.Error())
	}
	if err = t.store(key, v); err != nil {
		if sqliteErr, ok := err.(sqlite.Error); ok {
			if sqliteErr.Code == sqlite.ErrConstraint {
				return errors.Wrap(ErrDBADuplicateStore, err.Error())
			}
		}
		return err
	}
	return nil
}

func (t *DBASQLiteTx) Commit() error {
	return t.Tx.Commit()
}