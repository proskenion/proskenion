package main

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/inconshreveable/log15"
	"github.com/jessevdk/go-flags"
	"github.com/mattn/go-colorable"
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
	"log"
	"net"
)

var opts struct {
	// save to file name
	ConfigPath string `short:"c" long:"config" description:"A config path." value-name:"config/config.yaml" default-mask:"-"`
}

func main() {
	// ======= Arguents =======
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	// 1. set config
	configFile := "config/config.yaml"
	if opts.ConfigPath != "" {
		configFile = opts.ConfigPath
	}
	conf := config.NewConfig(configFile)

	// ========= loger setting ===========
	logger := log15.New("peerId", conf.Peer.Id)
	logger.SetHandler(log15.LvlFilterHandler(log15.LvlDebug, log15.StreamHandler(colorable.NewColorableStdout(), log15.TerminalFormat())))

	logger.Info("=================== boot proskenion ==========================")

	cryptor := crypto.NewEd25519Sha256Cryptor()

	db := dba.NewDBSQLite(conf)
	cmdExecutor := command.NewCommandExecutor(conf)
	cmdValidator := command.NewCommandValidator(conf)
	qVerifyier := query.NewQueryVerifier()
	fc := convertor.NewModelFactory(cryptor, cmdExecutor, cmdValidator, qVerifyier)

	rp := repository.NewRepository(db.DBA("kvstore"), cryptor, fc, conf)
	txQueue := repository.NewProposalTxQueueOnMemory(conf)
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
	csc := consensus.NewConsensus(rp, fc, cs, sync, bq, txListCache, gossip, pr, logger, conf, commitChan)

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
	cg := gate.NewConsensusGate(fc, cryptor, txQueue, txListCache, bq, conf)
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
