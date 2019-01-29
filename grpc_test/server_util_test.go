package grpc_test

import (
	"fmt"
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
	cf := client.NewClientFactory(fc, cryptor, conf)

	rp := repository.NewRepository(db.DBA("kvstore"), cryptor, fc, conf)
	txQueue := repository.NewProposalTxQueueOnMemory(conf)
	blockQueue := repository.NewProposalBlockQueueOnMemory(conf)
	bq := repository.NewProposalBlockQueueOnMemory(conf)
	txListCache := repository.NewTxListCache(conf)

	pr := prosl.NewProsl(fc, cryptor, conf)

	cmdExecutor.SetField(fc, pr)
	cmdValidator.SetField(fc, pr)

	qp := query.NewQueryProcessor(fc, conf)
	qv := query.NewQueryValidator(fc, conf)

	commitChan := make(chan struct{})
	cs := commit.NewCommitSystem(fc, cryptor, txQueue, rp, conf)

	gossip := p2p.NewGossip(rp, fc, cf, cryptor, conf)
	css := consensus.NewConsensus(rp, cs, bq, txListCache, gossip, pr, logger, conf, commitChan)

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
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Peer.Port))
	require.NoError(t, err)

	api := gate.NewAPIGate(rp, txQueue, qp, qv, logger)
	proskenion.RegisterAPIGateServer(s, controller.NewAPIGateServer(fc, api, logger))
	cg := gate.NewConsensusGate(fc, cryptor, txQueue, txListCache, blockQueue, logger, conf)
	proskenion.RegisterConsensusGateServer(s, controller.NewConsensusGateServer(fc, cg, cryptor, logger, conf))

	logger.Info("================= Consensus Boot =================")
	go func() {
		css.Boot()
	}()
	go func() {
		css.Receiver()
	}()

	if err := s.Serve(l); err != nil {
		require.NoError(t, err)
	}
}
