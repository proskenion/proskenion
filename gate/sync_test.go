package gate_test

import (
	. "github.com/proskenion/proskenion/gate"
	. "github.com/proskenion/proskenion/test_utils"
	"testing"
)

func TestNewSyncGate(t *testing.T) {
	//model.ModelFactory,
	//	core.CommandExecutor, core.CommandValidator,
	//	core.Cryptor, core.Repository, core.Prosl, *config.Config
	fc, _, _, c, rp, _, conf := NewTestFactories()
	NewSyncGate(rp, fc, c, RandomLogger(), conf)
}
