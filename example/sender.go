package main

import (
	"fmt"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/test_utils"
)

func main() {
	fc := RandomFactory()

	// servers Peer を登録
	confs := []*config.Config{
		config.NewConfig("configRoot.yaml"),
		config.NewConfig("config1.yaml"),
		config.NewConfig("config2.yaml"),
		config.NewConfig("config3.yaml"),
	}
	serversPeer := make([]model.Peer, 0)
	for _, conf := range confs {
		serversPeer = append(serversPeer,
			fc.NewPeer(conf.Peer.Id, fmt.Sprintf("%s:%s", conf.Peer.Host, conf.Peer.Port), conf.Peer.PublicKeyBytes()))
	}
	rootPeer := serversPeer[0]

	auditors := make([]*SenderManager, 0)
	for i := 0; i < 10; i++ {
		auditors = append(auditors,
			NewSenderManager(NewAccountWithPri(fmt.Sprintf("auditor%d@pr", i)), rootPeer))
	}

	// 1. 1台サーバを構築 (root)
	// 2. Creator を追加。
	// 3. +3台AddPeerす
	// 4. Creator がそれぞれ 信頼する Peer を選択する。
	// 5. Creator 同士で信頼(有効辺）を貼る。
	// 6. Creator の一人が新しい Consensus アルゴリズムを提案する。
	// 7. 他の Creator がそれに合意する。
	// 8. アルゴリズムの更新を行う。
	// 9. 合意形成を行うPeerが切り替わる
}
