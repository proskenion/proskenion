package grpc_test

import (
	"encoding/hex"
	"fmt"
	"github.com/fatih/color"
	"github.com/inconshreveable/log15"
	"github.com/mattn/go-colorable"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/prosl"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"io/ioutil"
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
	confs[1].Peer.PublicKey = "7ae0937f747fff11760db2f1a08d2e1892a25fdc7adb39714fb596081478d0a7"
	confs[1].Peer.PrivateKey = "03bc00e4e618de88a9f3148c53855284e46275367d57cb7357107c00624431977ae0937f747fff11760db2f1a08d2e1892a25fdc7adb39714fb596081478d0a7"
	confs[2].Peer.PublicKey = "f7dc24e3ac16779f071cc0bcc4971f0bc9d2ca3bf78047282796a0dcb9da7278"
	confs[2].Peer.PrivateKey = "67806d47c7b782d2691fa87cf1b45ceb38f32d187062120ff3d6f599068ace6df7dc24e3ac16779f071cc0bcc4971f0bc9d2ca3bf78047282796a0dcb9da7278"
	confs[3].Peer.PublicKey = "b3918c70db7e308d6b686c01ab0e08f3f677066eb8aba72c33f22b2798799635"
	confs[3].Peer.PrivateKey = "6fd47967cc1389e7e0c7838f35a8d5d42277a931c072251f7478a2544592c21db3918c70db7e308d6b686c01ab0e08f3f677066eb8aba72c33f22b2798799635"

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
	am.QueryPeersStatePassed(t, serversPeer)
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

	// Scenario 5 ===== AcceptEdge Creator <-> Creator ========
	for i, cm := range cms {
		go cm.CreateEdgeStorage(t, creators[i])
	}
	time.Sleep(time.Second * 2)
	edges := make([][]model.Object, 5)
	for i, cm := range cms {
		edges[i] = make([]model.Object, 0)
		go cm.AddEdge(t, creators[i], creators[(i+1)%5])
		edges[i] = append(edges[i], fc.NewObjectBuilder().Address(creators[(i+1)%5].AccountId))
		if i < 4 {
			go cm.AddEdge(t, creators[i], creators[(i+2)%5])
			edges[i] = append(edges[i], fc.NewObjectBuilder().Address(creators[(i+2)%5].AccountId))
		}
		if i < 3 {
			go cm.AddEdge(t, creators[i], creators[(i+3)%5])
			edges[i] = append(edges[i], fc.NewObjectBuilder().Address(creators[(i+3)%5].AccountId))
		}
		if i < 2 {
			go cm.AddEdge(t, creators[i], creators[(i+4)%5])
			edges[i] = append(edges[i], fc.NewObjectBuilder().Address(creators[(i+4)%5].AccountId))
		}
	}
	time.Sleep(time.Second * 3)
	w = &sync.WaitGroup{}
	for i, cm := range cms {
		w.Add(1)
		go func(cm *AccountManager, ac *AccountWithPri, es []model.Object) {
			cm.QueryStorageEdgesPassed(t, fmt.Sprintf("%s/%s", ac.AccountId, FollowStorage), es)
			w.Done()
		}(cm, creators[i], edges[i])
	}
	w.Wait()
	logger.Info(color.GreenString("Passed Scenario 5 : Follow Edge Creator <-> Creator"))

	// Scenario 6 ===== Propose NewConsensusAlgorithm =====
	pr := prosl.NewProsl(fc, RandomCryptor(), confs[0])
	newConY, err := ioutil.ReadFile("rep_consensus.yaml")
	require.NoError(t, err)
	err = pr.ConvertFromYaml(newConY)
	newCon, err := pr.Marshal()
	require.NoError(t, err)

	newIncY, err := ioutil.ReadFile("rep_incentive.yaml")
	require.NoError(t, err)
	err = pr.ConvertFromYaml(newIncY)
	require.NoError(t, err)
	newInc, err := pr.Marshal()
	require.NoError(t, err)

	// proposer == cms[0]
	cms[0].ProposeNewConsensus(t, newCon, newInc)
	cms[0].CreateProslSignStorage(t)
	time.Sleep(time.Second * 2)

	cms[0].QueryProslPassed(t, core.ConsensusKey, newCon)
	cms[0].QueryProslPassed(t, core.IncentiveKey, newInc)
	logger.Info(color.GreenString("Passed Scenario 6 : Propose NewConsensusAlgorithm."))

	// Senario 7 ===== Verify Creators new Conensus Algorithm =====
	incStj := fc.NewObjectBuilder().Storage(cms[0].QueryStorage(t, MakeIncentiveWalletId(cms[0].authorizer).Id()))
	conStj := fc.NewObjectBuilder().Storage(cms[0].QueryStorage(t, MakeConsensusWalletId(cms[0].authorizer).Id()))

	for _, cm := range cms[1:] {
		cm.VoteNewConsensus(t, cms[0].authorizer, core.IncentiveKey, incStj)
		cm.VoteNewConsensus(t, cms[0].authorizer, core.ConsensusKey, conStj)
	}
	time.Sleep(time.Second * 2)

	cms[0].QueryCollectSigsPassed(t, core.IncentiveKey, incStj, 4)
	cms[0].QueryCollectSigsPassed(t, core.ConsensusKey, conStj, 4)
	logger.Info(color.GreenString("Passed Scenario 7 : Verify Creators new Conensus Algorithm."))

	// Senario 8 ===== CheckAndCommit new Consensus Algorithm =====
	cms[0].CheckAndCommit(t)
	time.Sleep(time.Second * 2)
	cms[0].QueryRootProslPassed(t, incStj.GetStorage())
	cms[0].QueryRootProslPassed(t, conStj.GetStorage())
	logger.Info(color.GreenString("Passed Scenario 8 :  CheckAndCommit new Consensus Algorithm."))

	// server stop
	for _, server := range servers {
		server.GracefulStop()
	}
}
