package main

import (
	"encoding/hex"
	"flag"
	"github.com/inconshreveable/log15"
	"github.com/proskenion/proskenion/command"
	"github.com/proskenion/proskenion/commit"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/consensus"
	"github.com/proskenion/proskenion/convertor"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/crypto"
	"github.com/proskenion/proskenion/dba"
	"github.com/proskenion/proskenion/p2p"
	"github.com/proskenion/proskenion/query"
	"github.com/proskenion/proskenion/repository"
)

func main() {
	logger := log15.New()
	logger.Info("=================== boot proskenion ==========================")

	// Arguents ====
	configFile := "config/config.yaml"
	if len(flag.Args()) != 0 {
		configFile = flag.Arg(0)
	}

	conf := config.NewConfig(configFile)
	cryptor := crypto.NewEd25519Sha256Cryptor()

	// WIP : set public key and private key, this peer
	pub, pri := cryptor.NewKeyPairs()
	conf.Peer.PublicKey = hex.EncodeToString(pub)
	conf.Peer.PrivateKey = hex.EncodeToString(pri)

	db := dba.NewDBSQLite(conf)
	cmdExecutor := command.NewCommandExecutor()
	cmdValidator := command.NewCommandValidator()
	qValidator := query.NewQueryValidator()
	fc := convertor.NewModelFactory(cryptor, cmdExecutor, cmdValidator, qValidator)

	rp := repository.NewRepository(db.DBA("kvstore"), cryptor, fc)
	queue := repository.NewProposalTxQueueOnMemory(conf)

	commitChan := make(chan interface{})
	cs := commit.NewCommitSystem(fc, cryptor, queue, commit.DefaultCommitProperty(conf), rp)
	cc := consensus.NewMockCustomize(rp, commitChan)

	// WIP : mock
	gossip := &p2p.MockGossip{}
	consensus := consensus.NewConsensus(cc, cs, gossip, logger)

	// Genesis Commit
	genesisTxList := func() core.TxList {
		txList := repository.NewTxList(cryptor)
		txList.Push(fc.NewTxBuilder().
			CreateAccount("root", "root@com").
			AddPublicKey("root", "root@com", pub).
			Build())
		return txList
	}
	rp.GenesisCommit(genesisTxList())

	logger.Info("================= Consensus Boot =================")
	consensus.Boot()
}
