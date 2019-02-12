package main

import (
	"encoding/hex"
	"flag"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/inconshreveable/log15"
	"github.com/proskenion/proskenion/client"
	"github.com/proskenion/proskenion/command"
	"github.com/proskenion/proskenion/commit"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/consensus"
	"github.com/proskenion/proskenion/controller"
	"github.com/proskenion/proskenion/convertor"
	"github.com/proskenion/proskenion/crypto"
	"github.com/proskenion/proskenion/dba"
	"github.com/proskenion/proskenion/gate"
	"github.com/proskenion/proskenion/p2p"
	"github.com/proskenion/proskenion/prosl"
	"github.com/proskenion/proskenion/proto"
	"github.com/proskenion/proskenion/query"
	"github.com/proskenion/proskenion/repository"
	"github.com/proskenion/proskenion/synchronize"
	"google.golang.org/grpc"
	"net"
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
	cmdExecutor := command.NewCommandExecutor(conf)
	cmdValidator := command.NewCommandValidator(conf)
	qVerifyier := query.NewQueryVerifier()
	fc := convertor.NewModelFactory(cryptor, cmdExecutor, cmdValidator, qVerifyier)

	rp := repository.NewRepository(db.DBA("kvstore"), cryptor, fc, conf)
	txQueue := repository.NewProposalTxQueueOnMemory(conf)
	blockQueue := repository.NewProposalBlockQueueOnMemory(conf)
	bq := repository.NewProposalBlockQueueOnMemory(conf)
	txListCache := repository.NewTxListCache(conf)

	pr := prosl.NewProsl(fc, cryptor, conf)

	// cmd executor and validator set field.
	cmdExecutor.SetField(fc, pr)
	cmdValidator.SetField(fc, pr)

	qp := query.NewQueryProcessor(fc, conf)
	qv := query.NewQueryValidator(fc, conf)

	commitChan := make(chan struct{})
	cs := commit.NewCommitSystem(fc, cryptor, txQueue, rp, conf)

	// Gossip
	cf := client.NewClientFactory(fc, cryptor, conf)
	gossip := p2p.NewBroadCastGossip(rp, fc, cf, cryptor, conf)

	// sync
	sync := synchronize.NewSynchronizer(rp, cf, fc)

	// consensus
	csc := consensus.NewConsensus(rp, fc,cs,sync,bq, txListCache, gossip, pr, logger, conf, commitChan)

	// Genesis Commit
	logger.Info("================= Genesis Commit =================")
	genTxList, err := repository.GenesisTxListFromConf(cryptor, fc, rp, pr, conf)
	if err != nil {
		panic(err)
	}
	if err := rp.GenesisCommit(genTxList); err != nil {
		panic(err)
	}

	// ==================== gate =======================
	logger.Info("================= Gate Boot =================")
	l, err := net.Listen("tcp", ":"+conf.Peer.Port)
	if err != nil {
		panic(err.Error())
	}

	s := grpc.NewServer([]grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_validator.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	}...)
	api := gate.NewAPI(rp, txQueue, qp, qv, gossip, logger)
	proskenion.RegisterAPIServer(s, controller.NewAPIServer(fc, api, logger))
	cg := gate.NewConsensusGate(fc, cryptor, txQueue, txListCache, blockQueue, conf)
	proskenion.RegisterConsensusServer(s, controller.NewConsensusServer(fc, cg, cryptor, logger, conf))
	sg := gate.NewSyncGate(rp, fc, cryptor, conf)
	proskenion.RegisterSyncServer(s, controller.NewSyncServer(fc, sg, cryptor, logger, conf))

	// SetUp Consensus Loop
	go csc.Boot()
	go csc.Receiver()
	go csc.Patrol()

	if err := s.Serve(l); err != nil {
		logger.Error("Failed to server grpc: %s", err.Error())
	}
}
