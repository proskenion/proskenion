package grpc_test
/*
import (
	"fmt"
	"github.com/proskenion/proskenion/config"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/proskenion/proskenion/core/model"
	"google.golang.org/grpc"
	"sync"
	"testing"
	"time"
)

func TestScenario(t *testing.T) {
	// Boot Server
	confs := []*config.Config{
		config.NewConfig("config.yaml"),
		config.NewConfig("config.yaml"),
		config.NewConfig("config.yaml"),
		config.NewConfig("config.yaml"),
	}
	for i, _ := range confs {
		confs[i].Peer.Port = fmt.Sprintf("5005%d", 3+i)
		confs[i].DB.Name = fmt.Sprintf("testdb%d", i)
		if i > 0 {
			confs[i].Peer.Id = fmt.Sprintf("p%d@peer", i)
		}
	}

	fc := RandomFactory()
	servers := make([]*grpc.Server, 0)
	serversPeer := make([]model.Peer, 0)
	for i, conf := range confs {
		servers = append(servers, RandomServer())
		serversPeer = append(serversPeer,
			fc.NewPeer(conf.Peer.Id, fmt.Sprintf("%s:%s", conf.Peer.Host, conf.Peer.Port), conf.Peer.PublicKeyBytes()))
		go func(conf *config.Config, server *grpc.Server) {
			SetUpTestServer(t, conf, server)
		}(conf, servers[i])
	}
	time.Sleep(time.Second * 2)

	rootPeer := serversPeer[0]
	am := NewAccountManager(t, rootPeer)

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
	time.Sleep(time.Second * 5)
	ams := []*AccountManager{
		NewAccountManager(t, rootPeer),
		NewAccountManager(t, rootPeer),
		NewAccountManager(t, rootPeer),
		NewAccountManager(t, rootPeer),
		NewAccountManager(t, rootPeer),
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
	for _, server := range servers {
		server.GracefulStop()
	}
}
*/