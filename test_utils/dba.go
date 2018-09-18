package test_utils

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/dba"
)

func RandomDBA() core.DBA {
	return dba.NewDBAOnMemory()
}
