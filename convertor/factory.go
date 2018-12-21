package convertor

import (
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
	list := make([]*proskenion.Object, len(value))
	for _, object := range value {
		list = append(list, object.(*Object).Object)
	}

	b.Object[key] = &proskenion.Object{
		Type:   proskenion.ObjectCode_ListObjectCode,
		Object: &proskenion.Object_List{List: &proskenion.ObjectList{List: list}},
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

type AccountBulder struct {
	// TODO
}

type ModelFactory struct {
	model.ObjectFactory
	cryptor          core.Cryptor
	executor         core.CommandExecutor
	commandValidator core.CommandValidator
	queryValidator   core.QueryValidator
}

func NewModelFactory(cryptor core.Cryptor,
	executor core.CommandExecutor,
	cmdValidator core.CommandValidator,
	queryValidator core.QueryValidator) model.ModelFactory {
	factory := &ModelFactory{NewObjectFactory(cryptor), cryptor, executor, cmdValidator, queryValidator}
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
		f.queryValidator,
	}
}

func (f *ModelFactory) NewQueryResponseBuilder() model.QueryResponseBuilder {
	return &QueryResponseBuilder{
		&proskenion.QueryResponse{
			Payload:   &proskenion.QueryResponse_Payload{},
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

func (t *TxBuilder) CreateAccount(authorizerId string, accountId string) model.TxBuilder {
	t.Payload.Commands = append(t.Payload.Commands,
		&proskenion.Command{
			Command: &proskenion.Command_CreateAccount{
				CreateAccount: &proskenion.CreateAccount{},
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

func (t *TxBuilder) AddPublicKey(authorizerId string, accountId string, pubkey model.PublicKey) model.TxBuilder {
	t.Payload.Commands = append(t.Payload.Commands,
		&proskenion.Command{
			Command: &proskenion.Command_AddPublicKeys{
				AddPublicKeys: &proskenion.AddPublicKeys{
					PublicKeys: [][]byte{pubkey},
				},
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
	cryptor   core.Cryptor
	validator core.QueryValidator
}

func (q *QueryBuilder) AuthorizerId(authorizerId string) model.QueryBuilder {
	q.Query.Payload.AuthorizerId = authorizerId
	return q
}

func (q *QueryBuilder) TargetId(targetId string) model.QueryBuilder {
	q.Query.Payload.TargetId = targetId
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
	return &Query{q.Query, q.cryptor, q.validator}
}

type QueryResponseBuilder struct {
	*proskenion.QueryResponse
	cryptor core.Cryptor
}

func (q *QueryResponseBuilder) Account(ac model.Account) model.QueryResponseBuilder {
	q.QueryResponse.Payload.Object = &proskenion.QueryResponse_Payload_Account{
		Account: &proskenion.Account{
			AccountId:   ac.GetAccountId(),
			AccountName: ac.GetAccountName(),
			PublicKeys:  model.BytesListFromPublicKeys(ac.GetPublicKeys()),
			Balance:     ac.GetBalance(),
		},
	}
	q.QueryResponse.Payload.ResponseCode = proskenion.ObjectCode_AccountObjectCode
	return q
}

func (q *QueryResponseBuilder) Peer(p model.Peer) model.QueryResponseBuilder {
	q.QueryResponse.Payload.Object = &proskenion.QueryResponse_Payload_Peer{
		Peer: &proskenion.Peer{
			Address:   p.GetAddress(),
			PublicKey: p.GetPublicKey(),
		},
	}
	q.QueryResponse.Payload.ResponseCode = proskenion.ObjectCode_PeerObjectCode
	return q
}

func (q *QueryResponseBuilder) Build() model.QueryResponse {
	return &QueryResponse{q.QueryResponse, q.cryptor}
}
