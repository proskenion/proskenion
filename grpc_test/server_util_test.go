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
	"github.com/proskenion/proskenion/synchronize"
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

	logger := log15.New("peerId", conf.Peer.Id)
	logger.SetHandler(log15.LvlFilterHandler(log15.LvlDebug, log15.StdoutHandler))
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
	txListCache := repository.NewTxListCache(conf)

	pr := prosl.NewProsl(fc, cryptor, conf)

	cmdExecutor.SetField(fc, pr)
	cmdValidator.SetField(fc, pr)

	qp := query.NewQueryProcessor(fc, conf)
	qv := query.NewQueryValidator(fc, conf)

	commitChan := make(chan struct{})
	cs := commit.NewCommitSystem(fc, cryptor, txQueue, rp, conf)

	gossip := p2p.NewBroadCastGossip(rp, fc, cf, cryptor, conf)
	sync := synchronize.NewSynchronizer(rp, cf, fc)
	css := consensus.NewConsensus(rp, fc, cs, sync, blockQueue, txListCache, gossip, pr, logger, conf, commitChan)

	// Genesis Commit
	logger.Info("================= Genesis Commit =================")
	genTxList, err := repository.GenesisTxListFromConf(cryptor, fc, rp, pr, conf)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}
	if err := rp.GenesisCommit(genTxList); err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	// ==================== gate =======================
	logger.Info("================= Gate Boot =================")
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Peer.Port))
	require.NoError(t, err)

	api := gate.NewAPI(rp, txQueue, qp, qv, gossip, logger)
	proskenion.RegisterAPIServer(s, controller.NewAPIServer(fc, api, logger))
	cg := gate.NewConsensusGate(fc, cryptor, txQueue, txListCache, blockQueue, conf)
	proskenion.RegisterConsensusServer(s, controller.NewConsensusServer(fc, cg, cryptor, logger, conf))
	sg := gate.NewSyncGate(rp, fc, cryptor, conf)
	proskenion.RegisterSyncServer(s, controller.NewSyncServer(fc, sg, cryptor, logger, conf))

	// ==================== SetUp Consensus =======================
	go css.Boot()
	go css.Receiver()
	go css.Patrol()

	if err := s.Serve(l); err != nil {
		require.NoError(t, err)
	}
}
