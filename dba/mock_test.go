package dba_test

import (
	. "github.com/proskenion/proskenion/dba"
	"testing"
)

func TestDBAOnMemory_StoreAndLoad(t *testing.T) {
	dba := NewDBAOnMemory()
	testDBA_Store_Load(t, dba)
}

func TestDBAOnMemoryTx_StoreAndLoad(t *testing.T) {
	dba := NewDBAOnMemory()
	testDBATx_Store_Load(t, dba)
}
