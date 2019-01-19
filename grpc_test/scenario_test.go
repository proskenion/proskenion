package grpc_test

import (
	"github.com/proskenion/proskenion/config"
	. "github.com/proskenion/proskenion/test_utils"
	"sync"
	"testing"
	"time"
)

func TestScenario(t *testing.T) {
	return 
	// Boot Server
	conf := config.NewConfig("config.yaml")
	server := RandomServer()
	go func() {
		SetUpTestServer(t, conf, server)
	}()
	time.Sleep(time.Second * 2)

	fc := RandomFactory()
	serverPeer := fc.NewPeer(RandomStr(), ":50023", conf.Peer.PublicKeyBytes())
	am := NewAccountManager(t, serverPeer)

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
	time.Sleep(time.Second * 5)
	ams := []*AccountManager{
		NewAccountManager(t, serverPeer),
		NewAccountManager(t, serverPeer),
		NewAccountManager(t, serverPeer),
		NewAccountManager(t, serverPeer),
		NewAccountManager(t, serverPeer),
	}
	w := &sync.WaitGroup{}
	for i, ac := range acs {
		w.Add(1)
		go func(ac *AccountWithPri) {
			ams[i].QueryAccountPassed(t, ac)
			w.Done()
		}(ac)
	}
	w.Wait()

	// server stop
	server.GracefulStop()
}
