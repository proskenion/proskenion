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

func StrObject(str string, cryptor core.Cryptor) model.Object {
	return &Object{
		cryptor, nil, nil,
		&proskenion.Object{
			Type:   proskenion.ObjectCode_StringObjectCode,
			Object: &proskenion.Object_Str{str},
		},
	}
}

func BytesObject(bytes []byte, cryptor core.Cryptor) model.Object {
	return &Object{
		cryptor, nil, nil,
		&proskenion.Object{
			Type:   proskenion.ObjectCode_BytesObjectCode,
			Object: &proskenion.Object_Data{bytes},
		},
	}
}

func Int64Object(a int64, cryptor core.Cryptor) model.Object {
	return &Object{
		cryptor, nil, nil,
		&proskenion.Object{
			Type:   proskenion.ObjectCode_Int64ObjectCode,
			Object: &proskenion.Object_I64{a},
		},
	}
}

func Int32Object(a int32, cryptor core.Cryptor) model.Object {
	return &Object{
		cryptor, nil, nil,
		&proskenion.Object{
			Type:   proskenion.ObjectCode_Int32ObjectCode,
			Object: &proskenion.Object_I32{a},
		},
	}
}

func Uint64Object(a uint64, cryptor core.Cryptor) model.Object {
	return &Object{
		cryptor, nil, nil,
		&proskenion.Object{
			Type:   proskenion.ObjectCode_Uint64ObjectCode,
			Object: &proskenion.Object_U64{a},
		},
	}
}

func Uint32Object(a uint32, cryptor core.Cryptor) model.Object {
	return &Object{
		cryptor, nil, nil,
		&proskenion.Object{
			Type:   proskenion.ObjectCode_Uint32ObjectCode,
			Object: &proskenion.Object_U32{a},
		},
	}
}

func AddressObject(a string, cryptor core.Cryptor) model.Object {
	return &Object{
		cryptor, nil, nil,
		&proskenion.Object{
			Type:   proskenion.ObjectCode_AddressObjectCode,
			Object: &proskenion.Object_Address{a},
		},
	}
}

func ProsObjectListFromModelObjectList(objects []model.Object) []*proskenion.Object {
	ret := make([]*proskenion.Object, 0)
	for _, value := range objects {
		ret = append(ret, value.(*Object).Object)
	}
	return ret
}

func PrimitiveListObject(l []model.Object, cryptor core.Cryptor) model.Object {
	return &Object{
		cryptor, nil, nil,
		&proskenion.Object{
			Type:   proskenion.ObjectCode_ListObjectCode,
			Object: &proskenion.Object_List{&proskenion.ObjectList{List: ProsObjectListFromModelObjectList(l)}},
		},
	}
}

func PublicKeysToListObject(keys []model.PublicKey, cryptor core.Cryptor) model.Object {
	obs := make([]model.Object, 0, len(keys))
	for _, key := range keys {
		obs = append(obs, BytesObject(key, cryptor))
	}
	return PrimitiveListObject(obs, cryptor)
}

func (a *Account) GetFromKey(key string) model.Object {
	switch key {
	case "id", "account_id":
		return AddressObject(a.GetAccountId(), a.cryptor)
	case "name", "account_name":
		return StrObject(a.GetAccountName(), a.cryptor)
	case "keys", "public_keys":
		return PublicKeysToListObject(a.GetPublicKeys(), a.cryptor)
	case "balance":
		return Int64Object(a.GetBalance(), a.cryptor)
	case "quorum":
		return Int32Object(a.GetQuorum(), a.cryptor)
	case "peer_id", "delegate_peer_id", "peer", "delegate_peer":
		return AddressObject(a.GetDelegatePeerId(), a.cryptor)
	}
	return nil
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

func (a *Peer) GetFromKey(key string) model.Object {
	switch key {
	case "id", "peer_id":
		return AddressObject(a.GetPeerId(), a.cryptor)
	case "address", "ip":
		return StrObject(a.GetAddress(), a.cryptor)
	case "key", "public_key":
		return BytesObject(a.GetPublicKey(), a.cryptor)
	}
	return nil
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
	cryptor   core.Cryptor
	executor  core.CommandExecutor
	validator core.CommandValidator
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
			o.executor,
			o.validator,
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
			o.executor,
			o.validator,
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

func (o *Object) GetStorage() model.Storage {
	if o.Object == nil {
		return nil
	}
	return &Storage{o.cryptor, o.executor, o.validator, o.Object.GetStorage()}
}

func (o *Object) GetCommand() model.Command {
	if o.Object == nil {
		return nil
	}
	return &Command{o.Object.GetCommand(),
		o.executor,
		o.validator,
		o.cryptor}
}

func (o *Object) GetTransaction() model.Transaction {
	if o.Object == nil {
		return nil
	}
	return &Transaction{o.Object.GetTransaction(),
		o.cryptor, o.executor, o.validator}
}

func (o *Object) GetBlock() model.Block {
	if o.Object == nil {
		return nil
	}
	return &Block{o.Object.GetBlock(), o.cryptor}
}

func (o *Object) Cast(code model.ObjectCode) (model.Object, bool) {
	if o.GetType() == code {
		return o, true
	}
	switch code {
	case model.Int32ObjectCode:
		switch o.GetType() {
		case model.Int64ObjectCode:
			return Int32Object(int32(o.GetI64()), o.cryptor), true
		case model.Uint32ObjectCode:
			return Int32Object(int32(o.GetU32()), o.cryptor), true
		case model.Uint64ObjectCode:
			return Int32Object(int32(o.GetU64()), o.cryptor), true
		}
	case model.Int64ObjectCode:
		switch o.GetType() {
		case model.Int32ObjectCode:
			return Int64Object(int64(o.GetI32()), o.cryptor), true
		case model.Uint32ObjectCode:
			return Int64Object(int64(o.GetU32()), o.cryptor), true
		case model.Uint64ObjectCode:
			return Int64Object(int64(o.GetU64()), o.cryptor), true
		}
	case model.Uint32ObjectCode:
		switch o.GetType() {
		case model.Int32ObjectCode:
			return Uint32Object(uint32(o.GetI32()), o.cryptor), true
		case model.Int64ObjectCode:
			return Uint32Object(uint32(o.GetI64()), o.cryptor), true
		case model.Uint64ObjectCode:
			return Uint32Object(uint32(o.GetU64()), o.cryptor), true
		}
	case model.Uint64ObjectCode:
		switch o.GetType() {
		case model.Int32ObjectCode:
			return Uint64Object(uint64(o.GetI32()), o.cryptor), true
		case model.Int64ObjectCode:
			return Uint64Object(uint64(o.GetI64()), o.cryptor), true
		case model.Uint32ObjectCode:
			return Uint64Object(uint64(o.GetU32()), o.cryptor), true
		}
	case model.StringObjectCode:
		switch o.GetType() {
		case model.Int32ObjectCode:
			return StrObject(string(o.GetI32()), o.cryptor), true
		case model.Int64ObjectCode:
			return StrObject(string(o.GetI64()), o.cryptor), true
		case model.Uint32ObjectCode:
			return StrObject(string(o.GetU32()), o.cryptor), true
		case model.Uint64ObjectCode:
			return StrObject(string(o.GetU64()), o.cryptor), true
		}
	case model.BytesObjectCode:
		switch o.GetType() {
		case model.StringObjectCode:
			return BytesObject([]byte(o.GetStr()), o.cryptor), true
		}
	case model.AddressObjectCode:
		if o.GetType() == model.StringObjectCode {
			return AddressObject(o.GetStr(), o.cryptor), true
		}
	}
	return nil, false
}

func (o *Object) Marshal() ([]byte, error) {
	return proto.Marshal(o.Object)
}

func (o *Object) Unmarshal(pb []byte) error {
	return proto.Unmarshal(pb, o.Object)
}

func (o *Object) Hash() model.Hash {
	if o.Object == nil {
		return nil
	}
	return o.cryptor.Hash(o)
}
