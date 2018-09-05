package convertor

import (
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
)

type Account struct {
	*proskenion.Account
}

func (a *Account) GetPublicKeys() []model.PublicKey {
	if a.Account == nil {
		return nil
	}
	return model.PublicKeysFromBytesSlice(a.Account.GetPublicKeys())
}

type Peer struct {
	*proskenion.Peer
}

func (p *Peer) GetPublicKey() model.PublicKey {
	if p.Peer == nil {
		return nil
	}
	return p.Peer.GetPublicKey()
}