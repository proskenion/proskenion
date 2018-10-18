package commit

import (
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/core/model"
)

type CommitProperty struct {
	NumTxInBlock int
	PublicKey    model.PublicKey
	PrivateKey   model.PrivateKey
}

func DefaultCommitProperty(conf *config.Config) *CommitProperty {
	return &CommitProperty{
		conf.Commit.NumTxInBlock,
		conf.Peer.PublicKeyBytes(),
		conf.Peer.PrivateKeyBytes()}
}
