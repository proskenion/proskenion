package convertor

import (
	"github.com/gogo/protobuf/proto"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
)

type Account struct {
	cryptor core.Cryptor
	*proskenion.Account
}

func (a *Account) GetPublicKeys() []model.PublicKey {
	if a.Account == nil {
		return nil
	}
	return model.PublicKeysFromBytesSlice(a.Account.GetPublicKeys())
}

func (a *Account) Marshal() ([]byte, error) {
	return proto.Marshal(a.Account)
}

func (a *Account) Unmarshal(pb []byte) error {
	return proto.Unmarshal(pb, a.Account)
}

func (a *Account) Hash() model.Hash {
	if a.Account == nil {
		return nil
	}
	return a.cryptor.Hash(a)
}

type Peer struct {
	cryptor core.Cryptor
	*proskenion.Peer
}

func (p *Peer) GetPublicKey() model.PublicKey {
	if p.Peer == nil {
		return nil
	}
	return p.Peer.GetPublicKey()
}

func (a *Peer) Marshal() ([]byte, error) {
	return proto.Marshal(a.Peer)
}

func (a *Peer) Unmarshal(pb []byte) error {
	return proto.Unmarshal(pb, a.Peer)
}

func (a *Peer) Hash() model.Hash {
	if a.Peer == nil {
		return nil
	}
	return a.cryptor.Hash(a)
}

type ObjectList struct {
	cryptor core.Cryptor
	*proskenion.ObjectList
}

func (a *ObjectList) Marshal() ([]byte, error) {
	return proto.Marshal(a.ObjectList)
}

func (a *ObjectList) Unmarshal(pb []byte) error {
	return proto.Unmarshal(pb, a.ObjectList)
}

type Object struct {
	cryptor core.Cryptor
	*proskenion.Object
}
