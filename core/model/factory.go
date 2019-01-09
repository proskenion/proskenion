package model

import (
	"github.com/pkg/errors"
)

var (
	ErrNewBlock    = errors.Errorf("Failed Factory NewBlock")
	ErrNewProposal = errors.Errorf("Failed Factory NewProposal")
)

type ObjectFactory interface {
	NewSignature(pubkey PublicKey, signature []byte) Signature
	NewAccount(accountId string, accountName string, publicKeys []PublicKey, quorum int32, amount int64, peerId string) Account
	NewPeer(peerId string, address string, pubkey PublicKey) Peer

	NewStorageBuilder() StorageBuilder
	NewAccountBuilder() AccountBuilder
	NewObjectBuilder() ObjectBuilder

	NewEmptySignature() Signature
	NewEmptyAccount() Account
	NewEmptyPeer() Peer
	NewEmptyStorage() Storage
	NewEmptyObject() Object
}

type ModelFactory interface {
	ObjectFactory

	NewBlockBuilder() BlockBuilder
	NewTxBuilder() TxBuilder
	NewQueryBuilder() QueryBuilder
	NewQueryResponseBuilder() QueryResponseBuilder

	NewEmptyBlock() Block
	NewEmptyTx() Transaction
	NewEmptyQuery() Query
	NewEmptyQueryResponse() QueryResponse
}

type ObjectBuilder interface {
	Int32(value int32) ObjectBuilder
	Int64(value int64) ObjectBuilder
	Uint32(value uint32) ObjectBuilder
	Uint64(value uint64) ObjectBuilder
	Str(value string) ObjectBuilder
	Data(value []byte) ObjectBuilder
	Address(value string) ObjectBuilder
	Sig(value Signature) ObjectBuilder
	Account(value Account) ObjectBuilder
	Peer(value Peer) ObjectBuilder
	List(value []Object) ObjectBuilder
	Dict(value map[string]Object) ObjectBuilder
	Storage(value Storage) ObjectBuilder
	Build() Object
}

type StorageBuilder interface {
	From(Storage) StorageBuilder
	Int32(key string, value int32) StorageBuilder
	Int64(key string, value int64) StorageBuilder
	Uint32(key string, value uint32) StorageBuilder
	Uint64(key string, value uint64) StorageBuilder
	Str(key string, value string) StorageBuilder
	Data(key string, value []byte) StorageBuilder
	Address(key string, value string) StorageBuilder
	Sig(key string, value Signature) StorageBuilder
	Account(key string, value Account) StorageBuilder
	Peer(key string, value Peer) StorageBuilder
	List(key string, value []Object) StorageBuilder
	Dict(key string, value map[string]Object) StorageBuilder
	Set(key string, value Object) StorageBuilder
	Build() Storage
}

type AccountBuilder interface {
	From(Account) AccountBuilder
	AccountId(string) AccountBuilder
	AccountName(string) AccountBuilder
	Balance(int64) AccountBuilder
	PublicKeys([]PublicKey) AccountBuilder
	Quorum(int32) AccountBuilder
	DelegatePeerId(string) AccountBuilder
	Build() Account
}

type BlockBuilder interface {
	Height(int64) BlockBuilder
	PreBlockHash(Hash) BlockBuilder
	CreatedTime(int64) BlockBuilder
	WSVHash(Hash) BlockBuilder
	TxHistoryHash(Hash) BlockBuilder
	TxsHash(Hash) BlockBuilder
	Round(int32) BlockBuilder
	Build() Block
}

type TxBuilder interface {
	CreatedTime(int64) TxBuilder
	TransferBalance(srcAccountId string, destAccountId string, amount int64) TxBuilder
	CreateAccount(authorizerId string, accountId string, publicKeys []PublicKey, quorum int32) TxBuilder
	AddBalance(accountId string, amount int64) TxBuilder
	AddPublicKeys(authorizerId string, accountId string, pubkeys []PublicKey) TxBuilder
	RemovePublicKeys(authorizerId string, accountId string, pubkeys []PublicKey) TxBuilder
	SetQuorum(authorizerId string, accountId string, quorum int32) TxBuilder
	DefineStorage(authorizerId string, storageId string, storage Storage) TxBuilder
	CreateStorage(authorizerId string, storageId string) TxBuilder
	UpdateObject(authorizerId string, walletId string, key string, object Object) TxBuilder
	AddObject(authorizerId string, walletId string, key string, object Object) TxBuilder
	TransferObject(authorizerId string, walletId string, destAccountId string, key string, object Object) TxBuilder
	AddPeer(authorizerId string, peerId string, address string, pubkey PublicKey) TxBuilder
	Consign(authorizerId string, accountId string, peerId string) TxBuilder
	Build() Transaction
}

type QueryBuilder interface {
	AuthorizerId(string) QueryBuilder
	Select(string) QueryBuilder
	FromId(string) QueryBuilder
	Where([]byte) QueryBuilder
	OrderBy(key string, order OrderCode) QueryBuilder
	Limit(int32) QueryBuilder
	CreatedTime(int64) QueryBuilder
	RequestCode(code ObjectCode) QueryBuilder
	Build() Query
}

type QueryResponseBuilder interface {
	Account(Account) QueryResponseBuilder
	Peer(Peer) QueryResponseBuilder
	Storage(Storage) QueryResponseBuilder
	List([]Object) QueryResponseBuilder
	Object(Object) QueryResponseBuilder
	Build() QueryResponse
}
