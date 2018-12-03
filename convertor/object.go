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

type Object struct {
	cryptor core.Cryptor
	*proskenion.Object
}

func (o *Object) GetType() model.ObjectCode {
	return model.ObjectCode(o.Type)
}

func (o *Object) GetSig() model.Signature {
	if o.Object == nil {
		return nil
	}
	return &Signature{o.Object.GetSig()}
}

func (o *Object) GetAccount() model.Account {
	if o.Object == nil {
		return nil
	}
	return &Account{
		o.cryptor,
		o.Object.GetAccount(),
	}
}

func (o *Object) GetPeer() model.Peer {
	if o.Object == nil {
		return nil
	}
	return &Peer{
		o.cryptor,
		o.Object.GetPeer(),
	}
}

func (o *Object) modelObjectListFromProtoObjectList(objects *proskenion.ObjectList) []model.Object {
	ret := make([]model.Object, len(objects.GetList()))
	for i, object := range objects.GetList() {
		ret[i] = &Object{
			o.cryptor,
			object,
		}
	}
	return ret
}

func (o *Object) GetList() []model.Object {
	if o.Object == nil {
		return nil
	}
	return o.modelObjectListFromProtoObjectList(o.Object.GetList())
}

func (o *Object) modelObjectDictFromProtoObjectDict(objects *proskenion.ObjectDict) map[string]model.Object {
	ret := make(map[string]model.Object)
	for key, object := range objects.GetDict() {
		ret[key] = &Object{
			o.cryptor,
			object,
		}
	}
	return ret
}

func (o *Object) GetDict() map[string]model.Object {
	if o.Object == nil {
		return nil
	}
	return o.modelObjectDictFromProtoObjectDict(o.Object.GetDict())
}

func (o *Object) Marshal() ([]byte, error) {
	return proto.Marshal(o.Object)
}

func (o *Object) Unmarshal(pb []byte) error {
	return proto.Unmarshal(pb, o.Object)
}

func (o *Object) Hash() model.Hash {
	return o.cryptor.Hash(o)
}
