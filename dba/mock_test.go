package dba_test

import (
	. "github.com/proskenion/proskenion/dba"
	"testing"
)

func TestDBAOnMemory_StoreAndLoad(t *testing.T) {
	db := NewDBAnMemory()
	testDBA_Store_Load(t, db.DBA("test"))
}

func TestDBAOnMemoryTx_StoreAndLoad(t *testing.T) {
	db := NewDBAnMemory()
	testDBATx_Store_Load(t, db.DBA("test"))
}
