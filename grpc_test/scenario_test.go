package grpc_test

import (
	"github.com/proskenion/proskenion/config"
	. "github.com/proskenion/proskenion/test_utils"
	"sync"
	"testing"
	"time"
)

func TestScenario(t *testing.T) {
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
		NewAccountWithPri("target1@pr"),
		NewAccountWithPri("target2@pr"),
		NewAccountWithPri("target3@pr"),
		NewAccountWithPri("target4@pr"),
		NewAccountWithPri("target5@pr"),
	}

	// Scenario 1 ====== Create 5 Accounts ===================
	for _, ac := range acs {
		go func(ac *AccountWithPri) {
			am.CreateAccount(t, ac)
		}(ac)
	}
	time.Sleep(time.Second * 2)
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
