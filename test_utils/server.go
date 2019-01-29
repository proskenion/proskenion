package test_utils

import (
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/inconshreveable/log15"
	"github.com/proskenion/proskenion/command"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/controller"
	"github.com/proskenion/proskenion/convertor"
	"github.com/proskenion/proskenion/crypto"
	"github.com/proskenion/proskenion/gate"
	"github.com/proskenion/proskenion/prosl"
	"github.com/proskenion/proskenion/proto"
	"github.com/proskenion/proskenion/query"
	"github.com/proskenion/proskenion/repository"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"net"
	"testing"
)

func RandomSetUpConsensusServer(t *testing.T, conf *config.Config, s *grpc.Server) {
	logger := log15.New()
	logger.Info(fmt.Sprintf("=================== boot proskenion %s ==========================", conf.Peer.Port))

	cryptor := crypto.NewEd25519Sha256Cryptor()

	cmdExecutor := command.NewCommandExecutor(conf)
	cmdValidator := command.NewCommandValidator(conf)
	qVerifier := query.NewQueryVerifier()
	fc := convertor.NewModelFactory(cryptor, cmdExecutor, cmdValidator, qVerifier)

	txQueue := repository.NewProposalTxQueueOnMemory(conf)
	blockQueue := repository.NewProposalBlockQueueOnMemory(conf)
	txListCache := repository.NewTxListCache(conf)

	pr := prosl.NewProsl(fc, cryptor, conf)

	cmdExecutor.SetField(fc, pr)
	cmdValidator.SetField(fc, pr)

	// ==================== gate =======================
	logger.Info("================= Consensus Gate Boot =================")
	l, err := net.Listen("tcp", ":"+conf.Peer.Port)
	require.NoError(t, err)

	cg := gate.NewConsensusGate(fc, cryptor, txQueue, txListCache, blockQueue, logger, conf)
	proskenion.RegisterConsensusGateServer(s, controller.NewConsensusGateServer(fc, cg, cryptor, logger, conf))

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
