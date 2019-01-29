package dba_test

import (
	. "github.com/proskenion/proskenion/dba"
	. "github.com/proskenion/proskenion/test_utils"
	"testing"
)

func TestDBASQLite_StoreAndLoad(t *testing.T) {
	conf := RandomConfig()
	db := NewDBSQLite(conf)
	testDBA_Store_Load(t, db.DBA("test"))
}

func TestDBASQLiteTx_StoreAndLoad(t *testing.T) {
	conf := RandomConfig()
	db := NewDBSQLite(conf)
	testDBATx_Store_Load(t, db.DBA("test"))
}

func TestDBASQLite_Parallel(t *testing.T) {
	conf := RandomConfig()
	db := NewDBSQLite(conf)
	testDBA_Parallel(t, db.DBA("test"))
}
