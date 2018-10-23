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

	assert.Equal(t, conf.Commit.NumTxInBlock, 1000)

	assert.Equal(t, conf.ProposalTxsLimits, 1000)

	assert.Equal(t, conf.Peer.Port, "50023")
}
