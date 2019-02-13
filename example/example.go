package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/inconshreveable/log15"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/convertor"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/crypto"
	"github.com/proskenion/proskenion/prosl"
	"github.com/proskenion/proskenion/query"
	. "github.com/proskenion/proskenion/test_utils"
	"io/ioutil"
	"sync"
	"time"
)

func WaitSecond(second int) {
	time.Sleep(time.Duration(second) * time.Second)
}

const NUM_CREATORS = 10

func main() {
	cryptor := crypto.NewEd25519Sha256Cryptor()
	fc := convertor.NewModelFactory(cryptor, nil, nil, query.NewQueryVerifier())

	logger := log15.New()

	// servers Peer を登録
	confs := []*config.Config{
		config.NewConfig("example/configRoot.yaml"),
		config.NewConfig("example/config1.yaml"),
		config.NewConfig("example/config2.yaml"),
		config.NewConfig("example/config3.yaml"),
	}
	serversPeer := make([]model.Peer, 0)
	for _, conf := range confs {
		serversPeer = append(serversPeer,
			fc.NewPeer(conf.Peer.Id, fmt.Sprintf("%s:%s", conf.Peer.Host, conf.Peer.Port), conf.Peer.PublicKeyBytes()))
	}
	rootPeer := serversPeer[0]
	logger.Debug(rootPeer.GetAddress())

	// 1. authorizer を登録
	authorizer := NewSenderManager(NewAccountWithPri("authorizer@pr"), rootPeer, fc, confs[0])
	authorizer.SetAuthorizer()

	WaitSecond(2)
	authorizer.QueryAccountPassed(authorizer.Authorizer)
	logger.Info(color.BlueString("Registered Authorizer PublicKey."))

	creators := make([]*SenderManager, 0)
	for i := 0; i < NUM_CREATORS; i++ {
		creators = append(creators,
			NewSenderManager(NewAccountWithPri(fmt.Sprintf("creator%d@creator.pr", i)), rootPeer, fc, confs[0]))
	}

	// 1. 1台サーバを構築 (root)
	// `docker run proskenion -c example/configRoot.yaml`
	// 2. +3台AddPeer
	// ```
	// export LOCAL_HOST_IP=`ifconfig en0 | grep inet | grep -v inet6 | sed -E "s/inet ([0-9]{1,3}.[0-9]{1,3}.[0-9].{1,3}.[0-9]{1,3}) .*$/\1/" | tr -d "\t"`
	// docker run -p $LOCAL_HOST_IP:50052:50052 proskenion -c example/configRoot.yaml
	// docker run -p $LOCAL_HOST_IP:50053:50053 proskenion:latest -c example/config1.yaml
	// docker run -p $LOCAL_HOST_IP:50054:50054 proskenion:latest -c example/config2.yaml
	// docker run-p $LOCAL_HOST_IP:50055:50055 proskenion:latest -c example/config3.yaml
	// ```

	// 3. Creator を追加。
	logger.Info(color.BlueString("================== Scenario 1 :: Create Creators. =================="))
	for _, creator := range creators {
		authorizer.CreateAccount(creator.Authorizer)
	}

	WaitSecond(2)
	for _, creator := range creators {
		creator.QueryAccountPassed(creator.Authorizer)
	}
	logger.Info(color.GreenString("===================== :: Passed Scenario 1 :: ====================="))

	// 4. Creator がそれぞれ 信頼する Peer を選択する。
	logger.Info(color.BlueString("================ Scenario 2 :: Degrade 5 Creators  ================"))
	for i, creator := range creators {
		go creator.Consign(creator.Authorizer, serversPeer[i%4])
	}

	WaitSecond(3)
	w := &sync.WaitGroup{}
	for i, cm := range creators {
		w.Add(1)
		go func(cm *SenderManager, ac *AccountWithPri, peer model.Peer) {
			cm.QueryAccountDegradedPassed(ac, peer)
			w.Done()
		}(cm, cm.Authorizer, serversPeer[i%4])
	}
	w.Wait()
	logger.Info(color.GreenString("===================== :: Passed Scenario 2 :: ====================="))

	// 5. Creator 同士で信頼(有効辺）を貼る。
	logger.Info(color.BlueString("=================== Scenario 3 :: Follow Edges  ==================="))
	for _, cm := range creators {
		go cm.CreateEdgeStorage(cm.Authorizer)
	}
	time.Sleep(time.Second * 2)
	edges := make([][]model.Object, NUM_CREATORS)
	for i, cm := range creators {
		edges[i] = make([]model.Object, 0)
		go cm.AddEdge(cm.Authorizer, creators[(i+1)%NUM_CREATORS].Authorizer)
		edges[i] = append(edges[i], fc.NewObjectBuilder().Address(creators[(i+1)%NUM_CREATORS].Authorizer.AccountId))
		if i < 8 {
			go cm.AddEdge(cm.Authorizer, creators[(i+2)%NUM_CREATORS].Authorizer)
			edges[i] = append(edges[i], fc.NewObjectBuilder().Address(creators[(i+2)%NUM_CREATORS].Authorizer.AccountId))
		}
		if i < 6 {
			go cm.AddEdge(cm.Authorizer, creators[(i+3)%NUM_CREATORS].Authorizer)
			edges[i] = append(edges[i], fc.NewObjectBuilder().Address(creators[(i+3)%NUM_CREATORS].Authorizer.AccountId))
		}
		if i < 4 {
			go cm.AddEdge(cm.Authorizer, creators[(i+4)%NUM_CREATORS].Authorizer)
			edges[i] = append(edges[i], fc.NewObjectBuilder().Address(creators[(i+4)%NUM_CREATORS].Authorizer.AccountId))
		}
	}

	WaitSecond(3)
	w = &sync.WaitGroup{}
	for i, cm := range creators {
		w.Add(1)
		go func(cm *SenderManager, ac *AccountWithPri, es []model.Object) {
			cm.QueryStorageEdgesPassed(fmt.Sprintf("%s/%s", ac.AccountId, FollowStorage), es)
			w.Done()
		}(cm, cm.Authorizer, edges[i])
	}
	w.Wait()
	logger.Info(color.GreenString("===================== :: Passed Scenario 3 :: ====================="))

	// 6. Creator の一人が新しい Consensus アルゴリズムを提案する。
	logger.Info(color.BlueString("=========== Scenario 4 :: Propose NewConsensusAlgorithm  ==========="))
	pr := prosl.NewProsl(fc, RandomCryptor(), confs[0])
	newConY, err := ioutil.ReadFile("example/rep_consensus.yaml")
	RequireNoError(err)
	err = pr.ConvertFromYaml(newConY)
	newCon, err := pr.Marshal()
	RequireNoError(err)

	newIncY, err := ioutil.ReadFile("example/rep_incentive.yaml")
	RequireNoError(err)
	err = pr.ConvertFromYaml(newIncY)
	RequireNoError(err)
	newInc, err := pr.Marshal()
	RequireNoError(err)

	creators[0].ProposeNewConsensus(newCon, newInc) // proposer is creators[0]
	creators[0].CreateProslSignStorage()

	WaitSecond(2)
	creators[0].QueryProslPassed(core.ConsensusKey, newCon)
	creators[0].QueryProslPassed(core.IncentiveKey, newInc)
	logger.Info(color.GreenString("===================== :: Passed Scenario 4 :: ====================="))

	// 7. 他の Creator がそれに合意する。
	logger.Info(color.BlueString("============ Scenario 5 :: Verify NewConensusAlgorithm  ============"))
	incStj := fc.NewObjectBuilder().Storage(creators[0].QueryStorage(MakeIncentiveWalletId(creators[0].Authorizer).Id()))
	conStj := fc.NewObjectBuilder().Storage(creators[0].QueryStorage(MakeConsensusWalletId(creators[0].Authorizer).Id()))

	for _, cm := range creators[1:] {
		cm.VoteNewConsensus(creators[0].Authorizer, core.IncentiveKey, incStj)
		cm.VoteNewConsensus(creators[0].Authorizer, core.ConsensusKey, conStj)
	}

	WaitSecond(2)
	creators[0].QueryCollectSigsPassed(core.IncentiveKey, incStj, 9)
	creators[0].QueryCollectSigsPassed(core.ConsensusKey, conStj, 9)
	logger.Info(color.GreenString("===================== :: Passed Scenario 5 :: ====================="))

	// 8. アルゴリズムの更新を行う。
	logger.Info(color.BlueString("======= Scenario 6 :: CheckAndCommit NewConsensusAlgorithm  ======="))
	creators[0].CheckAndCommit()

	WaitSecond(2)
	creators[0].QueryRootProslPassed(incStj.GetStorage())
	creators[0].QueryRootProslPassed(conStj.GetStorage())
	logger.Info(color.GreenString("===================== :: Passed Scenario 6 :: ====================="))

	// 9. 合意形成を行うPeerが切り替わる. fin
	for {
		WaitSecond(5)
		authorizer.QueryAccountsBalances()
		logger.Info(color.GreenString("===================== :: Waiting 5 seconds :: ====================="))
	}
}
