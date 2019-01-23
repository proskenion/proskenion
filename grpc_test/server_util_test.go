package grpc_test

import (
	"fmt"
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
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"net"
	"os"
	"testing"
)

func clearData(t *testing.T, conf *config.Config) {
	os.Remove(fmt.Sprintf("%s/%s.sqlite", conf.DB.Path, conf.DB.Name))
}

func SetUpTestServer(t *testing.T, conf *config.Config, s *grpc.Server) {
	clearData(t, conf)

	logger := log15.New()
	logger.Info(fmt.Sprintf("=================== boot proskenion %s ==========================", conf.Peer.Port))

	cryptor := crypto.NewEd25519Sha256Cryptor()

	db := dba.NewDBSQLite(conf)
	cmdExecutor := command.NewCommandExecutor(conf)
	cmdValidator := command.NewCommandValidator(conf)
	qVerifier := query.NewQueryVerifier()
	fc := convertor.NewModelFactory(cryptor, cmdExecutor, cmdValidator, qVerifier)

	rp := repository.NewRepository(db.DBA("kvstore"), cryptor, fc, conf)
	queue := repository.NewProposalTxQueueOnMemory(conf)

	pr := prosl.NewProsl(fc, rp, cryptor, conf)

	cmdExecutor.SetField(fc, pr)
	cmdValidator.SetField(fc, pr)

	qp := query.NewQueryProcessor(rp, fc, conf)
	qv := query.NewQueryValidator(rp, fc, conf)

	commitChan := make(chan struct{})
	cs := commit.NewCommitSystem(fc, cryptor, queue, rp, conf)

	// WIP : mock
	gossip := &p2p.MockGossip{}
	css := consensus.NewConsensus(rp, cs, gossip, pr, logger, conf, commitChan)

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
	require.NoError(t, err)

	api := gate.NewAPIGate(queue, qp, qv, logger)
	proskenion.RegisterAPIGateServer(s, controller.NewAPIGateServer(fc, api, logger))

	logger.Info("================= Consensus Boot =================")
	go func() {
		css.Boot()
	}()

	if err := s.Serve(l); err != nil {
		require.NoError(t, err)
	}
}

func RandomServer() *grpc.Server {
	return grpc.NewServer([]grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_validator.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	}...)
}
