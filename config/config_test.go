package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfig(t *testing.T) {
	conf := NewConfig("config.yaml")

	assert.Equal(t, conf.DB.Path, "database")
	assert.Equal(t, conf.DB.Kind, "sqlite3")
	assert.Equal(t, conf.DB.Name, "db")

	assert.Equal(t, conf.Commit.WaitInterval, 1000)
	assert.Equal(t, conf.Commit.NumTxInBlock, 1000)

	assert.Equal(t, conf.Queue.TxsLimits, 1000)
	assert.Equal(t, conf.Queue.BlockLimits, 30)

	assert.Equal(t, conf.Cache.TxListLimits, 100)

	assert.Equal(t, conf.Peer.Port, "50023")

	assert.Equal(t, conf.Prosl.Id, "/prosl")

	assert.Equal(t, conf.Prosl.Incentive.Id, "incentive/prosl")
	assert.Equal(t, conf.Prosl.Consensus.Id, "consensus/prosl")
	assert.Equal(t, conf.Prosl.Update.Id, "update/prosl")

	assert.Equal(t, conf.Prosl.Genesis.Path, "./grpc_test/genesis.yaml")
	assert.Equal(t, conf.Prosl.Incentive.Path, "./grpc_test/incentive.yaml")
	assert.Equal(t, conf.Prosl.Consensus.Path, "./grpc_test/consensus.yaml")
	assert.Equal(t, conf.Prosl.Update.Path, "./grpc_test/update.yaml")

}
