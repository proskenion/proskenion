package main

import (
	"fmt"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/convertor"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/crypto"
	"github.com/proskenion/proskenion/query"
	. "github.com/proskenion/proskenion/test_utils"
	"time"
)

func WaitSecond(second int) {
	time.Sleep(time.Duration(second) * time.Second)
}

func main() {
	cryptor := crypto.NewEd25519Sha256Cryptor()
	fc := convertor.NewModelFactory(cryptor, nil, nil, query.NewQueryVerifier())

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
	fmt.Println(rootPeer.GetAddress())

	// 1. authorizer を登録
	authorizer := NewSenderManager(NewAccountWithPri("authorizer@pr"), rootPeer, fc)
	authorizer.SetAuthorizer()

	WaitSecond(2)
	authorizer.QueryAccountPassed(authorizer.Authorizer)
	fmt.Println("Registered Authorizer PublicKey.")

	creators := make([]*SenderManager, 0)
	for i := 0; i < 10; i++ {
		creators = append(creators,
			NewSenderManager(NewAccountWithPri(fmt.Sprintf("creator%d@creator.pr", i)), rootPeer, fc))
	}

	// 1. 1台サーバを構築 (root)
	// `docker run proskenion -c example/configRoot.yaml`
	// 2. Creator を追加。
	for _, creator := range creators {
		authorizer.CreateAccount(creator.Authorizer)
	}

	WaitSecond(2)
	for _, creator := range creators {
		creator.QueryAccountPassed(creator.Authorizer)
	}
	fmt.Println("================== Scenario 1 :: Create Creators. ==================")

	// 3. +3台AddPeerす
	// 4. Creator がそれぞれ 信頼する Peer を選択する。
	// 5. Creator 同士で信頼(有効辺）を貼る。
	// 6. Creator の一人が新しい Consensus アルゴリズムを提案する。
	// 7. 他の Creator がそれに合意する。
	// 8. アルゴリズムの更新を行う。
	// 9. 合意形成を行うPeerが切り替わる
}
