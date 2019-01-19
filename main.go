package main

import (
	"encoding/hex"
	"flag"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/inconshreveable/log15"
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
	queue := repository.NewProposalTxQueueOnMemory(conf)

	pr := prosl.NewProsl(fc, rp, cryptor, conf)

	// cmd executor and validator set field.
	cmdExecutor.SetField(fc, pr)
	cmdValidator.SetField(fc, pr)

	qp := query.NewQueryProcessor(rp, fc, conf)
	qv := query.NewQueryValidator(rp, fc, conf)

	commitChan := make(chan interface{})
	cs := commit.NewCommitSystem(fc, cryptor, queue, commit.DefaultCommitProperty(conf), rp)
	cc := consensus.NewMockCustomize(rp, commitChan)

	// WIP : mock
	gossip := &p2p.MockGossip{}
	csc := consensus.NewConsensus(cc, cs, gossip, logger)

	// Genesis Commit
	logger.Info("================= Genesis Commit =================")
	genTxList, err := repository.NewTxListFromConf(cryptor, pr, conf)
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

	api := gate.NewAPIGate(queue, qp, qv, logger)
	s := grpc.NewServer([]grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_validator.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	}...)
	proskenion.RegisterAPIGateServer(s, controller.NewAPIGateServer(fc, api, logger))

	logger.Info("================= Consensus Boot =================")
	go func() {
		csc.Boot()
	}()

	if err := s.Serve(l); err != nil {
		logger.Error("Failed to server grpc: %s", err.Error())
	}
}
