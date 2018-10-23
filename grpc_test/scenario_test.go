package grpc_test

import (
	"github.com/proskenion/proskenion/config"
	. "github.com/proskenion/proskenion/test_utils"
	"testing"
	"time"
)

func TestScenario(t *testing.T) {
	conf := config.NewConfig("config.yaml")
	fc := NewTestFactory()
	server := fc.NewPeer("localhost:50023", conf.Peer.PublicKeyBytes())
	am := NewAccountMnager(t, server)

	// set authorizer
	am.SetAuthorizer(t)

	acs := []*AccountWithPri{
		NewAccountWithPri("target1@com"),
		NewAccountWithPri("target2@com"),
		NewAccountWithPri("target3@com"),
		NewAccountWithPri("target4@com"),
		NewAccountWithPri("target5@com"),
	}

	// Scenario 1 ====== Create 5 Accounts ===================
	for _, ac := range acs {
		go func() {
			am.CreateAccount(t, ac)
		}()
	}
	time.Sleep(time.Second * 2)
	for _, ac := range acs {
		go func() {
			am.QueryAccountPassed(t, ac)
		}()
	}
}
