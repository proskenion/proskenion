package convertor

import (
	"github.com/golang/protobuf/proto"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
)

type ObjectFactory struct {
	cryptor core.Cryptor
}

func NewObjectFactory(cryptor core.Cryptor) model.ObjectFactory {
	return &ObjectFactory{cryptor}
}

func (f *ObjectFactory) NewEmptySignature() model.Signature {
	return &Signature{
		&proskenion.Signature{},
	}
}

func (f *ObjectFactory) NewEmptyAccount() model.Account {
	return &Account{
		f.cryptor,
		&proskenion.Account{},
	}
}

func (f *ObjectFactory) NewEmptyPeer() model.Peer {
	return &Peer{
		f.cryptor,
		&proskenion.Peer{},
	}
}

func (f *ObjectFactory) NewSignature(pubkey model.PublicKey, signature []byte) model.Signature {
	return &Signature{
		&proskenion.Signature{
			PublicKey: []byte(pubkey),
			Signature: signature,
		},
	}
}

func (f *ObjectFactory) NewAccount(accountId string, accountName string, publicKeys []model.PublicKey, quorum int32, balance int64, peerId string) model.Account {
	return &Account{
		f.cryptor,
		&proskenion.Account{
			AccountId:      accountId,
			AccountName:    accountName,
			PublicKeys:     model.BytesListFromPublicKeys(publicKeys),
			Balance:        balance,
			Quorum:         quorum,
			DelegatePeerId: peerId,
		},
	}
}

func (f *ObjectFactory) NewAccountBuilder() model.AccountBuilder {
	return &AccountBuilder{
		cryptor: f.cryptor,
		Account: &proskenion.Account{},
	}
}

func (f *ObjectFactory) NewPeer(peerId string, address string, pubkey model.PublicKey) model.Peer {
	return &Peer{
		f.cryptor,
		&proskenion.Peer{
			PeerId:    peerId,
			Address:   address,
			PublicKey: []byte(pubkey),
		},
	}
}

func (f *ObjectFactory) NewObjectBuilder() model.ObjectBuilder {
	return &ObjectBuilder{f.cryptor, &proskenion.Object{}}
}

func (f *ObjectFactory) NewStorageBuilder() model.StorageBuilder {
	return &StorageBuilder{
		f.cryptor,
		&proskenion.Storage{Object: make(map[string]*proskenion.Object)},
	}
}

func (f *ObjectFactory) NewEmptyStorage() model.Storage {
	return &Storage{
		f.cryptor,
		&proskenion.Storage{Object: make(map[string]*proskenion.Object)},
	}
}

func (f *ObjectFactory) NewEmptyObject() model.Object {
	return &Object{
		f.cryptor,
		&proskenion.Object{},
	}
}

type ObjectBuilder struct {
	cryptor core.Cryptor
	*proskenion.Object
}

func (f *ObjectBuilder) Int32(value int32) model.ObjectBuilder {
	f.Object = &proskenion.Object{
		Type:   proskenion.ObjectCode_Int32ObjectCode,
		Object: &proskenion.Object_I32{I32: value},
	}
	return f
}

func (f *ObjectBuilder) Int64(value int64) model.ObjectBuilder {
	f.Object = &proskenion.Object{
		Type:   proskenion.ObjectCode_Int64ObjectCode,
		Object: &proskenion.Object_I64{I64: value},
	}
	return f
}

func (f *ObjectBuilder) Uint32(value uint32) model.ObjectBuilder {
	f.Object = &proskenion.Object{
		Type:   proskenion.ObjectCode_Uint32ObjectCode,
		Object: &proskenion.Object_U32{U32: value},
	}
	return f
}

func (f *ObjectBuilder) Uint64(value uint64) model.ObjectBuilder {
	f.Object = &proskenion.Object{
		Type:   proskenion.ObjectCode_Uint64ObjectCode,
		Object: &proskenion.Object_U64{U64: value},
	}
	return f
}

func (f *ObjectBuilder) Str(value string) model.ObjectBuilder {
	f.Object = &proskenion.Object{
		Type:   proskenion.ObjectCode_StringObjectCode,
		Object: &proskenion.Object_Str{Str: value},
	}
	return f
}

func (f *ObjectBuilder) Data(value []byte) model.ObjectBuilder {
	f.Object = &proskenion.Object{
		Type:   proskenion.ObjectCode_BytesObjectCode,
		Object: &proskenion.Object_Data{Data: value},
	}
	return f
}

func (f *ObjectBuilder) Address(value string) model.ObjectBuilder {
	f.Object = &proskenion.Object{
		Type:   proskenion.ObjectCode_AddressObjectCode,
		Object: &proskenion.Object_Address{Address: value},
	}
	return f
}

func (f *ObjectBuilder) Sig(value model.Signature) model.ObjectBuilder {
	f.Object = &proskenion.Object{
		Type: proskenion.ObjectCode_SignatureObjectCode,
		Object: &proskenion.Object_Sig{Sig: &proskenion.Signature{
			PublicKey: value.GetPublicKey(),
			Signature: value.GetSignature()}},
	}
	return f
}

func (f *ObjectBuilder) Account(value model.Account) model.ObjectBuilder {
	f.Object = &proskenion.Object{
		Type:   proskenion.ObjectCode_AccountObjectCode,
		Object: &proskenion.Object_Account{Account: value.(*Account).Account},
	}
	return f
}

func (f *ObjectBuilder) Peer(value model.Peer) model.ObjectBuilder {
	f.Object = &proskenion.Object{
		Type:   proskenion.ObjectCode_PeerObjectCode,
		Object: &proskenion.Object_Peer{Peer: value.(*Peer).Peer},
	}
	return f
}

func (f *ObjectBuilder) List(value []model.Object) model.ObjectBuilder {
	f.Object = &proskenion.Object{
		Type:   proskenion.ObjectCode_ListObjectCode,
		Object: &proskenion.Object_List{List: &proskenion.ObjectList{List: ProslObjectListFromObjectList(value)}},
	}
	return f
}

func (f *ObjectBuilder) Dict(value map[string]model.Object) model.ObjectBuilder {
	f.Object = &proskenion.Object{
		Type:   proskenion.ObjectCode_DictObjectCode,
		Object: &proskenion.Object_Dict{Dict: &proskenion.ObjectDict{Dict: ProslObjectMapsFromObjectMaps(value)}},
	}
	return f
}

func (f *ObjectBuilder) Storage(value model.Storage) model.ObjectBuilder {
	f.Object = &proskenion.Object{
		Type:   proskenion.ObjectCode_StorageObjectCode,
		Object: &proskenion.Object_Storage{Storage: value.(*Storage).Storage},
	}
	return f
}

func (f *ObjectBuilder) Build() model.Object {
	return &Object{f.cryptor, f.Object}
}

type StorageBuilder struct {
	cryptor core.Cryptor
	*proskenion.Storage
}

func (b *StorageBuilder) Int32(key string, value int32) model.StorageBuilder {
	b.Object[key] = &proskenion.Object{
		Type:   proskenion.ObjectCode_Int32ObjectCode,
		Object: &proskenion.Object_I32{I32: value},
	}
	return b
}

func (b *StorageBuilder) Int64(key string, value int64) model.StorageBuilder {
	b.Object[key] = &proskenion.Object{
		Type:   proskenion.ObjectCode_Int64ObjectCode,
		Object: &proskenion.Object_I64{I64: value},
	}
	return b
}

func (b *StorageBuilder) Uint32(key string, value uint32) model.StorageBuilder {
	b.Object[key] = &proskenion.Object{
		Type:   proskenion.ObjectCode_Uint32ObjectCode,
		Object: &proskenion.Object_U32{U32: value},
	}
	return b
}

func (b *StorageBuilder) Uint64(key string, value uint64) model.StorageBuilder {
	b.Object[key] = &proskenion.Object{
		Type:   proskenion.ObjectCode_Uint64ObjectCode,
		Object: &proskenion.Object_U64{U64: value},
	}
	return b
}

func (b *StorageBuilder) Str(key string, value string) model.StorageBuilder {
	b.Object[key] = &proskenion.Object{
		Type:   proskenion.ObjectCode_StringObjectCode,
		Object: &proskenion.Object_Str{Str: value},
	}
	return b
}

func (b *StorageBuilder) Data(key string, value []byte) model.StorageBuilder {
	b.Object[key] = &proskenion.Object{
		Type:   proskenion.ObjectCode_BytesObjectCode,
		Object: &proskenion.Object_Data{Data: value},
	}
	return b
}

func (b *StorageBuilder) Address(key string, value string) model.StorageBuilder {
	b.Object[key] = &proskenion.Object{
		Type:   proskenion.ObjectCode_AddressObjectCode,
		Object: &proskenion.Object_Address{Address: value},
	}
	return b
}

func (b *StorageBuilder) Sig(key string, value model.Signature) model.StorageBuilder {
	b.Object[key] = &proskenion.Object{
		Type:   proskenion.ObjectCode_SignatureObjectCode,
		Object: &proskenion.Object_Sig{Sig: value.(*Signature).Signature},
	}
	return b
}

func (b *StorageBuilder) Account(key string, value model.Account) model.StorageBuilder {
	b.Object[key] = &proskenion.Object{
		Type:   proskenion.ObjectCode_AccountObjectCode,
		Object: &proskenion.Object_Account{Account: value.(*Account).Account},
	}
	return b
}

func (b *StorageBuilder) Peer(key string, value model.Peer) model.StorageBuilder {
	b.Object[key] = &proskenion.Object{
		Type:   proskenion.ObjectCode_PeerObjectCode,
		Object: &proskenion.Object_Peer{Peer: value.(*Peer).Peer},
	}
	return b
}

func (b *StorageBuilder) List(key string, value []model.Object) model.StorageBuilder {
	b.Object[key] = &proskenion.Object{
		Type:   proskenion.ObjectCode_ListObjectCode,
		Object: &proskenion.Object_List{List: &proskenion.ObjectList{List: ProslObjectListFromObjectList(value)}},
	}
	return b
}

func (b *StorageBuilder) Dict(key string, value map[string]model.Object) model.StorageBuilder {
	dict := make(map[string]*proskenion.Object)
	for key, object := range value {
		dict[key] = object.(*Object).Object
	}
	b.Object[key] = &proskenion.Object{
		Type:   proskenion.ObjectCode_DictObjectCode,
		Object: &proskenion.Object_Dict{Dict: &proskenion.ObjectDict{Dict: dict}},
	}
	return b
}

func (b *StorageBuilder) Build() model.Storage {
	return &Storage{
		b.cryptor,
		b.Storage,
	}
}

type AccountBuilder struct {
	cryptor core.Cryptor
	*proskenion.Account
}

func (b *AccountBuilder) From(a model.Account) model.AccountBuilder {
	b.Account = &proskenion.Account{
		AccountId:      a.GetAccountId(),
		AccountName:    a.GetAccountName(),
		PublicKeys:     model.BytesListFromPublicKeys(a.GetPublicKeys()),
		Quorum:         a.GetQuorum(),
		Balance:        a.GetBalance(),
		DelegatePeerId: a.GetDelegatePeerId(),
	}
	return b
}

func (b *AccountBuilder) AccountId(id string) model.AccountBuilder {
	b.Account.AccountId = id
	return b
}

func (b *AccountBuilder) AccountName(name string) model.AccountBuilder {
	b.Account.AccountName = name
	return b
}

func (b *AccountBuilder) Balance(balance int64) model.AccountBuilder {
	b.Account.Balance = balance
	return b
}

func (b *AccountBuilder) PublicKeys(keys []model.PublicKey) model.AccountBuilder {
	b.Account.PublicKeys = model.BytesListFromPublicKeys(keys)
	return b
}

func (b *AccountBuilder) Quorum(quorum int32) model.AccountBuilder {
	b.Account.Quorum = quorum
	return b
}

func (b *AccountBuilder) DelegatePeerId(dpeerId string) model.AccountBuilder {
	b.Account.DelegatePeerId = dpeerId
	return b
}

func (b *AccountBuilder) Build() model.Account {
	return &Account{
		cryptor: b.cryptor,
		Account: b.Account,
	}
}

type ModelFactory struct {
	model.ObjectFactory
	cryptor          core.Cryptor
	executor         core.CommandExecutor
	commandValidator core.CommandValidator
	queryVerifier    core.QueryVerifier
}

func NewModelFactory(cryptor core.Cryptor,
	executor core.CommandExecutor,
	cmdValidator core.CommandValidator,
	queryVerifier core.QueryVerifier) model.ModelFactory {
	factory := &ModelFactory{NewObjectFactory(cryptor), cryptor, executor, cmdValidator, queryVerifier}
	executor.SetFactory(factory)
	cmdValidator.SetFactory(factory)
	return factory
}

func (f *ModelFactory) NewEmptyBlock() model.Block {
	return &Block{
		&proskenion.Block{
			Payload:   &proskenion.Block_Payload{},
			Signature: &proskenion.Signature{},
		},
		f.cryptor,
	}
}

func (f *ModelFactory) NewEmptyTx() model.Transaction {
	return f.NewTxBuilder().Build()
}

func (f *ModelFactory) NewEmptyQuery() model.Query {
	return f.NewQueryBuilder().Build()
}

func (f *ModelFactory) NewEmptyQueryResponse() model.QueryResponse {
	return f.NewQueryResponseBuilder().Build()
}

func (f *ModelFactory) NewBlockBuilder() model.BlockBuilder {
	return &BlockBuilder{
		&proskenion.Block{
			Payload:   &proskenion.Block_Payload{},
			Signature: &proskenion.Signature{},
		},
		f.cryptor,
	}
}

func (f *ModelFactory) NewTxBuilder() model.TxBuilder {
	return &TxBuilder{
		&proskenion.Transaction{
			Payload:    &proskenion.Transaction_Payload{},
			Signatures: make([]*proskenion.Signature, 0),
		},
		f.cryptor,
		f.executor,
		f.commandValidator,
	}
}

func (f *ModelFactory) NewQueryBuilder() model.QueryBuilder {
	return &QueryBuilder{
		&proskenion.Query{
			Payload:   &proskenion.Query_Payload{},
			Signature: &proskenion.Signature{},
		},
		f.cryptor,
		f.queryVerifier,
	}
}

func (f *ModelFactory) NewQueryResponseBuilder() model.QueryResponseBuilder {
	return &QueryResponseBuilder{
		&proskenion.QueryResponse{
			Object:    &proskenion.Object{},
			Signature: &proskenion.Signature{},
		},
		f.cryptor,
	}
}

type BlockBuilder struct {
	*proskenion.Block
	cryptor core.Cryptor
}

func (b *BlockBuilder) Height(height int64) model.BlockBuilder {
	b.Block.Payload.Height = height
	return b
}
func (b *BlockBuilder) PreBlockHash(hash model.Hash) model.BlockBuilder {
	b.Block.Payload.PreBlockHash = hash
	return b
}
func (b *BlockBuilder) CreatedTime(time int64) model.BlockBuilder {
	b.Block.Payload.CreatedTime = time
	return b
}
func (b *BlockBuilder) WSVHash(hash model.Hash) model.BlockBuilder {
	b.Block.Payload.WsvHash = hash
	return b
}
func (b *BlockBuilder) TxHistoryHash(hash model.Hash) model.BlockBuilder {
	b.Block.Payload.TxHistoryHash = hash
	return b
}
func (b *BlockBuilder) TxsHash(hash model.Hash) model.BlockBuilder {
	b.Block.Payload.TxsHash = hash
	return b
}

func (b *BlockBuilder) Round(round int32) model.BlockBuilder {
	b.Block.Payload.Round = round
	return b
}

func (b *BlockBuilder) Build() model.Block {
	return &Block{b.Block, b.cryptor}
}

type TxBuilder struct {
	*proskenion.Transaction
	cryptor   core.Cryptor
	executor  core.CommandExecutor
	validator core.CommandValidator
}

func (t *TxBuilder) CreatedTime(time int64) model.TxBuilder {
	t.Payload.CreatedTime = time
	return t
}

func (t *TxBuilder) TransferBalance(srcAccountId string, destAccountId string, balance int64) model.TxBuilder {
	t.Payload.Commands = append(t.Payload.Commands,
		&proskenion.Command{
			Command: &proskenion.Command_TransferBalance{
				TransferBalance: &proskenion.TransferBalance{
					DestAccountId: destAccountId,
					Balance:       balance,
				},
			},
			TargetId:     srcAccountId,
			AuthorizerId: srcAccountId,
		})
	return t
}

func (t *TxBuilder) CreateAccount(authorizerId string, accountId string, publicKeys []model.PublicKey, quorum int32) model.TxBuilder {
	t.Payload.Commands = append(t.Payload.Commands,
		&proskenion.Command{
			Command: &proskenion.Command_CreateAccount{
				CreateAccount: &proskenion.CreateAccount{
					PublicKeys: model.BytesListFromPublicKeys(publicKeys),
					Quorum:     quorum,
				},
			},
			TargetId:     accountId,
			AuthorizerId: authorizerId,
		})
	return t
}

func (t *TxBuilder) AddBalance(accountId string, balance int64) model.TxBuilder {
	t.Payload.Commands = append(t.Payload.Commands,
		&proskenion.Command{
			Command: &proskenion.Command_AddBalance{
				AddBalance: &proskenion.AddBalance{
					Balance: balance,
				},
			},
			TargetId:     accountId,
			AuthorizerId: accountId,
		})
	return t
}

func (t *TxBuilder) AddPublicKeys(authorizerId string, accountId string, pubkey []model.PublicKey) model.TxBuilder {
	t.Payload.Commands = append(t.Payload.Commands,
		&proskenion.Command{
			Command: &proskenion.Command_AddPublicKeys{
				AddPublicKeys: &proskenion.AddPublicKeys{
					PublicKeys: model.BytesListFromPublicKeys(pubkey),
				},
			},
			TargetId:     accountId,
			AuthorizerId: authorizerId,
		})
	return t
}

func (t *TxBuilder) RemovePublicKeys(authorizerId string, accountId string, pubkeys []model.PublicKey) model.TxBuilder {
	t.Payload.Commands = append(t.Payload.Commands,
		&proskenion.Command{
			Command: &proskenion.Command_RemovePublicKeys{
				RemovePublicKeys: &proskenion.RemovePublicKeys{
					PublicKeys: model.BytesListFromPublicKeys(pubkeys),
				},
			},
			TargetId:     accountId,
			AuthorizerId: authorizerId,
		})
	return t
}

func (t *TxBuilder) SetQuorum(authorizerId string, accountId string, quorum int32) model.TxBuilder {
	t.Payload.Commands = append(t.Payload.Commands,
		&proskenion.Command{
			Command: &proskenion.Command_SetQurum{
				SetQurum: &proskenion.SetQuorum{
					Quorum: quorum,
				},
			},
			TargetId:     accountId,
			AuthorizerId: authorizerId,
		})
	return t
}

func (t *TxBuilder) DefineStorage(authorizerId string, storageId string, storage model.Storage) model.TxBuilder {
	t.Payload.Commands = append(t.Payload.Commands,
		&proskenion.Command{
			Command: &proskenion.Command_DefineStorage{
				DefineStorage: &proskenion.DefineStorage{
					Storage: storage.(*Storage).Storage,
				},
			},
			TargetId:     storageId,
			AuthorizerId: authorizerId,
		})
	return t
}

func (t *TxBuilder) CreateStorage(authorizerId string, storageId string) model.TxBuilder {
	t.Payload.Commands = append(t.Payload.Commands,
		&proskenion.Command{
			Command: &proskenion.Command_CreateStorage{
				CreateStorage: &proskenion.CreateStorage{},
			},
			TargetId:     storageId,
			AuthorizerId: authorizerId,
		})
	return t
}

func (t *TxBuilder) UpdateObject(authorizerId string, walletId string, key string, object model.Object) model.TxBuilder {
	t.Payload.Commands = append(t.Payload.Commands,
		&proskenion.Command{
			Command: &proskenion.Command_UpdateObject{
				UpdateObject: &proskenion.UpdateObject{
					Key:    key,
					Object: object.(*Object).Object,
				},
			},
			TargetId:     walletId,
			AuthorizerId: authorizerId,
		})
	return t
}

func (t *TxBuilder) AddObject(authorizerId string, walletId string, key string, object model.Object) model.TxBuilder {
	t.Payload.Commands = append(t.Payload.Commands,
		&proskenion.Command{
			Command: &proskenion.Command_AddObject{
				AddObject: &proskenion.AddObject{
					Key:    key,
					Object: object.(*Object).Object,
				},
			},
			TargetId:     walletId,
			AuthorizerId: authorizerId,
		})
	return t
}

func (t *TxBuilder) TransferObject(authorizerId string, walletId string, destAccountId string, key string, object model.Object) model.TxBuilder {
	t.Payload.Commands = append(t.Payload.Commands,
		&proskenion.Command{
			Command: &proskenion.Command_TransferObject{
				TransferObject: &proskenion.TransferObject{
					Key:           key,
					DestAccountId: destAccountId,
					Object:        object.(*Object).Object,
				},
			},
			TargetId:     walletId,
			AuthorizerId: authorizerId,
		})
	return t
}

func (t *TxBuilder) AddPeer(authorizerId string, accountId string, address string, pubkey model.PublicKey) model.TxBuilder {
	t.Payload.Commands = append(t.Payload.Commands,
		&proskenion.Command{
			Command: &proskenion.Command_AddPeer{
				AddPeer: &proskenion.AddPeer{
					Address:   address,
					PublicKey: pubkey,
				},
			},
			TargetId:     accountId,
			AuthorizerId: authorizerId,
		})
	return t
}

func (t *TxBuilder) Consign(authorizerId string, accountId string, peerId string) model.TxBuilder {
	t.Payload.Commands = append(t.Payload.Commands,
		&proskenion.Command{
			Command: &proskenion.Command_Consign{
				Consign: &proskenion.Consign{PeerId: peerId},
			},
			TargetId:     accountId,
			AuthorizerId: authorizerId,
		})
	return t
}

func (t *TxBuilder) Build() model.Transaction {
	return &Transaction{t.Transaction,
		t.cryptor, t.executor, t.validator}
}

type QueryBuilder struct {
	*proskenion.Query
	cryptor  core.Cryptor
	verifier core.QueryVerifier
}

func (q *QueryBuilder) AuthorizerId(authorizerId string) model.QueryBuilder {
	q.Query.Payload.AuthorizerId = authorizerId
	return q
}

func (q *QueryBuilder) Select(selectId string) model.QueryBuilder {
	q.Query.Payload.Select = selectId
	return q
}

func (q *QueryBuilder) FromId(fromId string) model.QueryBuilder {
	q.Query.Payload.FromId = fromId
	return q
}

func (q *QueryBuilder) Where(where []byte) model.QueryBuilder {
	cond := &proskenion.ConditionalFormula{}
	proto.Unmarshal(where, cond)
	q.Payload.Where = cond
	return q
}

func (q *QueryBuilder) OrderBy(key string, order model.OrderCode) model.QueryBuilder {
	q.Payload.OrderBy = &proskenion.Query_OrderBy{
		Key:   key,
		Order: proskenion.Query_Order(order),
	}
	return q
}

func (q *QueryBuilder) Limit(limit int32) model.QueryBuilder {
	q.Payload.Limit = limit
	return q
}

func (q *QueryBuilder) CreatedTime(time int64) model.QueryBuilder {
	q.Query.Payload.CreatedTime = time
	return q
}

func (q *QueryBuilder) RequestCode(code model.ObjectCode) model.QueryBuilder {
	q.Query.Payload.RequstCode = proskenion.ObjectCode(code)
	return q
}

func (q *QueryBuilder) Build() model.Query {
	return &Query{q.Query, q.cryptor, q.verifier}
}

type QueryResponseBuilder struct {
	*proskenion.QueryResponse
	cryptor core.Cryptor
}

func (q *QueryResponseBuilder) Account(ac model.Account) model.QueryResponseBuilder {
	q.QueryResponse.Object = &proskenion.Object{
		Type: proskenion.ObjectCode_AccountObjectCode,
		Object: &proskenion.Object_Account{Account: &proskenion.Account{
			AccountId:   ac.GetAccountId(),
			AccountName: ac.GetAccountName(),
			PublicKeys:  model.BytesListFromPublicKeys(ac.GetPublicKeys()),
			Balance:     ac.GetBalance(),
		}}}
	return q
}

func (q *QueryResponseBuilder) Peer(p model.Peer) model.QueryResponseBuilder {
	q.QueryResponse.Object = &proskenion.Object{
		Type: proskenion.ObjectCode_PeerObjectCode,
		Object: &proskenion.Object_Peer{Peer: &proskenion.Peer{
			Address:   p.GetAddress(),
			PublicKey: p.GetPublicKey(),
		}}}
	return q
}

func (q *QueryResponseBuilder) Storage(s model.Storage) model.QueryResponseBuilder {
	q.QueryResponse.Object = &proskenion.Object{
		Type: proskenion.ObjectCode_StorageObjectCode,
		Object: &proskenion.Object_Storage{Storage: &proskenion.Storage{
			Object: ProslObjectMapsFromObjectMaps(s.GetObject()),
		}}}
	return q
}

func (q *QueryResponseBuilder) List(os []model.Object) model.QueryResponseBuilder {
	q.QueryResponse.Object = &proskenion.Object{
		Type:   proskenion.ObjectCode_ListObjectCode,
		Object: &proskenion.Object_List{List: &proskenion.ObjectList{List: ProslObjectListFromObjectList(os)}}}
	return q
}

func (q *QueryResponseBuilder) Build() model.QueryResponse {
	return &QueryResponse{q.QueryResponse, q.cryptor}
}
