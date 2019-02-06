package grpc_test

import (
	"encoding/hex"
	"fmt"
	"github.com/fatih/color"
	"github.com/inconshreveable/log15"
	"github.com/mattn/go-colorable"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/test_utils"
	"google.golang.org/grpc"
	"sync"
	"testing"
	"time"
)

func TestScenario(t *testing.T) {
	logger := log15.New("test", "ScenarioLog!!")
	logger.SetHandler(log15.StreamHandler(colorable.NewColorableStdout(), log15.TerminalFormat()))

	// Boot Server
	confs := []*config.Config{
		config.NewConfig("config.yaml"),
		config.NewConfig("config.yaml"),
		config.NewConfig("config.yaml"),
		config.NewConfig("config.yaml"),
		config.NewConfig("config.yaml"),
		config.NewConfig("config.yaml"),
		config.NewConfig("config.yaml"),
		config.NewConfig("config.yaml"),
	}
	for i, _ := range confs {
		confs[i].Peer.Port = fmt.Sprintf("5005%d", 2+i)
		confs[i].DB.Name = fmt.Sprintf("testdb%d", i)
		if i > 0 {
			confs[i].Peer.Id = fmt.Sprintf("p%d@peer", i)
		}
		if i > 3 {
			confs[i].Peer.Active = false
			pub, pri := RandomCryptor().NewKeyPairs()
			strPub, strPri := hex.EncodeToString(pub), hex.EncodeToString(pri)
			confs[i].Peer.PublicKey = strPub
			confs[i].Peer.PrivateKey = strPri
		}
	}

	fc := RandomFactory()
	servers := make([]*grpc.Server, 0)
	serversPeer := make([]model.PeerWithPriKey, 0)
	for i, conf := range confs[:4] {
		servers = append(servers, RandomServer())
		serversPeer = append(serversPeer,
			&PeerWithPri{fc.NewPeer(conf.Peer.Id, fmt.Sprintf("%s:%s", conf.Peer.Host, conf.Peer.Port), conf.Peer.PublicKeyBytes()),
				conf.Peer.PrivateKeyBytes()})
		go func(conf *config.Config, server *grpc.Server) {
			SetUpTestServer(t, conf, server)
		}(conf, servers[i])
	}
	time.Sleep(time.Second * 2)

	rootPeer := serversPeer[0]
	authorizer := NewAccountWithPri("authorizer@pr")
	am := NewAccountManager(t, authorizer, rootPeer)

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
		NewAccountManager(t, acs[0], rootPeer),
		NewAccountManager(t, acs[1], rootPeer),
		NewAccountManager(t, acs[2], serversPeer[1]),
		NewAccountManager(t, acs[3], serversPeer[2]),
		NewAccountManager(t, acs[4], serversPeer[3]),
	}
	w := &sync.WaitGroup{}
	for i, ac := range acs {
		w.Add(1)
		go func(am *AccountManager, ac *AccountWithPri) {
			am.QueryAccountPassed(t, ac)
			w.Done()
		}(ams[i], ac)
	}
	w.Wait()
	logger.Info(color.GreenString("Passed Scenario 1 : CreateAccount."))

	// Scenario 2 ===== Create 5 Creators(Accounts) =======
	creators := []*AccountWithPri{
		NewAccountWithPri("alis@creator"),
		NewAccountWithPri("bob@creator"),
		NewAccountWithPri("hiroshi@creator"),
		NewAccountWithPri("migawari@creator"),
		NewAccountWithPri("hihi@creator"),
	}
	for _, ac := range creators {
		go func(ac *AccountWithPri) {
			am.CreateAccount(t, ac)
		}(ac)
	}
	time.Sleep(time.Second * 2)
	am.QueryRangeAccountsPassed(t, "creator/"+model.AccountStorageName, creators)
	logger.Info(color.GreenString("Passed Scenario 2 : Creators"))

	// Scenario 3 ===== Sync another 4 Peers ===============
	for i, conf := range confs[4:] {
		servers = append(servers, RandomServer())
		serversPeer = append(serversPeer,
			&PeerWithPri{
				fc.NewPeer(conf.Peer.Id, fmt.Sprintf("%s:%s", conf.Peer.Host, conf.Peer.Port), conf.Peer.PublicKeyBytes()),
				conf.Peer.PrivateKeyBytes(),
			})
		go func(conf *config.Config, server *grpc.Server) {
			SetUpTestServer(t, conf, server)
		}(conf, servers[i+4])
		am.AddPeer(t, serversPeer[i+4])
	}
	time.Sleep(time.Second * 3)
	am.QueryPeersState(t, serversPeer)
	logger.Info(color.GreenString("Passed Scenario 3 : Sync another 4 Peers"))

	// Scenario 4 ==== Degreade 5 Creators[0...4] -> 5 Peers[1...5] ======
	cms := make([]*AccountManager, 0, len(creators))
	for i, ac := range creators {
		cms = append(cms, NewAccountManager(t, ac, serversPeer[i+1]))
	}
	for i, cm := range cms {
		go cm.Consign(t, creators[i], serversPeer[i+1])
	}
	time.Sleep(time.Second * 3)
	w = &sync.WaitGroup{}
	for i, cm := range cms {
		w.Add(1)
		go func(cm *AccountManager, ac *AccountWithPri, peer model.PeerWithPriKey) {
			cm.QueryAccountDegradedPassed(t, ac, peer)
			w.Done()
		}(cm, creators[i], serversPeer[i+1])
	}
	w.Wait()
	logger.Info(color.GreenString("Passed Scenario 4 : Degrade 5 Creators -> 5 Peers"))

	// server stop
	for _, server := range servers {
		server.GracefulStop()
	}
}
