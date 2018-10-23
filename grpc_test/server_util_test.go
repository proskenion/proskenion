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
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/crypto"
	"github.com/proskenion/proskenion/dba"
	"github.com/proskenion/proskenion/gate"
	"github.com/proskenion/proskenion/p2p"
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
	cmdExecutor := command.NewCommandExecutor()
	cmdValidator := command.NewCommandValidator()
	qValidator := query.NewQueryValidator()
	fc := convertor.NewModelFactory(cryptor, cmdExecutor, cmdValidator, qValidator)

	rp := repository.NewRepository(db.DBA("kvstore"), cryptor, fc)
	queue := repository.NewProposalTxQueueOnMemory(conf)

	qp := query.NewQueryProcessor(rp, fc, conf)

	commitChan := make(chan interface{})
	cs := commit.NewCommitSystem(fc, cryptor, queue, commit.DefaultCommitProperty(conf), rp)
	cc := consensus.NewMockCustomize(rp, commitChan)

	// WIP : mock
	gossip := &p2p.MockGossip{}
	consensus := consensus.NewConsensus(cc, cs, gossip, logger)

	// Genesis Commit
	logger.Info("================= Genesis Commit =================")
	genesisTxList := func() core.TxList {
		txList := repository.NewTxList(cryptor)
		txList.Push(fc.NewTxBuilder().
			CreateAccount("root", "root@com").
			AddPublicKey("root", "root@com", conf.Peer.PublicKeyBytes()).
			CreateAccount("root", "authorizer@com").
			Build())
		return txList
	}
	require.NoError(t, rp.GenesisCommit(genesisTxList()))

	// ==================== gate =======================
	logger.Info("================= Gate Boot =================")
	l, err := net.Listen("tcp", ":"+conf.Peer.Port)
	require.NoError(t, err)

	api := gate.NewAPIGate(queue, qp, logger)
	proskenion.RegisterAPIGateServer(s, controller.NewAPIGateServer(fc, api, logger))

	logger.Info("================= Consensus Boot =================")
	go func() {
		consensus.Boot()
	}()

	if err := s.Serve(l); err != nil {
		require.NoError(t, err)
	}
}

func NewTestServer() *grpc.Server {
	return grpc.NewServer([]grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_validator.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	}...)
}
