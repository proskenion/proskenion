package grpc_test

import (
	"github.com/proskenion/proskenion/config"
	. "github.com/proskenion/proskenion/test_utils"
	"testing"
	"time"
)

func TestScenario(t *testing.T) {

	// Boot Server
	conf := config.NewConfig("config.yaml")
	server := NewTestServer()
	go func() {
		SetUpTestServer(t, conf, server)
	}()
	time.Sleep(time.Second * 2)

	fc := NewTestFactory()
	serverPeer := fc.NewPeer(":50023", conf.Peer.PublicKeyBytes())
	am := NewAccountMnager(t, serverPeer)

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
		go func(ac *AccountWithPri) {
			am.CreateAccount(t, ac)
		}(ac)
	}
	time.Sleep(time.Second * 2)
	ams := []*AccountManager{
		NewAccountMnager(t, serverPeer),
		NewAccountMnager(t, serverPeer),
		NewAccountMnager(t, serverPeer),
		NewAccountMnager(t, serverPeer),
		NewAccountMnager(t, serverPeer),
	}
	for i, ac := range acs {
		go func(ac *AccountWithPri) {
			ams[i].QueryAccountPassed(t, ac)
		}(ac)
	}

	// server stop
	server.GracefulStop()
}
