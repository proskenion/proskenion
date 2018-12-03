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
	NewAccount(accountId string, accountName string, publicKeys []PublicKey, amount int64) Account
	NewPeer(address string, pubkey PublicKey) Peer

	NewStorageBuilder() StorageBuilder
	NewObjectListBuilder() ObjectListBuilder
	NewObjectDictBuilder() ObjectDictBuilder

	NewEmptyAccount() Account
	NewEmptyPeer() Peer
	NewEmptyStorage() Storage
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

type StorageBuilder interface {
	Int32(key string, value int32) StorageBuilder
	Int64(key string, value int64) StorageBuilder
	Uint32(key string, value uint32) StorageBuilder
	Uint64(key string, value uint32) StorageBuilder
	Str(key string, value string) StorageBuilder
	Data(key string, value []byte) StorageBuilder
	Address(key string, value string) StorageBuilder
	Sig(key string, value Signature) StorageBuilder
	Account(key string, value Account) StorageBuilder
	Peer(key string, value Peer) StorageBuilder
	List(key string, value ObjectList) StorageBuilder
	Dict(key string, value ObjectDict) StorageBuilder
	Build() Storage
}

type ObjectListBuilder interface {
	Int32(key string, value int32) ObjectListBuilder
	Int64(key string, value int64) ObjectListBuilder
	Uint32(key string, value uint32) ObjectListBuilder
	Uint64(key string, value uint32) ObjectListBuilder
	Str(key string, value string) ObjectListBuilder
	Data(key string, value []byte) ObjectListBuilder
	Address(key string, value string) ObjectListBuilder
	Sig(key string, value Signature) ObjectListBuilder
	Account(key string, value Account) ObjectListBuilder
	Peer(key string, value Peer) ObjectListBuilder
	List(key string, value []Object) ObjectListBuilder
	Build() ObjectList
}

type ObjectDictBuilder interface {
	Int32(key string, value int32) ObjectDictBuilder
	Int64(key string, value int64) ObjectDictBuilder
	Uint32(key string, value uint32) ObjectDictBuilder
	Uint64(key string, value uint32) ObjectDictBuilder
	Str(key string, value string) ObjectDictBuilder
	Data(key string, value []byte) ObjectDictBuilder
	Address(key string, value string) ObjectDictBuilder
	Sig(key string, value Signature) ObjectDictBuilder
	Account(key string, value Account) ObjectDictBuilder
	Peer(key string, value Peer) ObjectDictBuilder
	List(key string, value []Object) ObjectDictBuilder
	Build() ObjectDict
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
	Transfer(srcAccountId string, destAccountId string, amount int64) TxBuilder
	CreateAccount(authorizerId string, accountId string) TxBuilder
	AddAsset(accountId string, amount int64) TxBuilder
	AddPublicKey(authorizerId string, accountId string, pubkey PublicKey) TxBuilder
	Build() Transaction
}

type QueryBuilder interface {
	AuthorizerId(string) QueryBuilder
	TargetId(string) QueryBuilder
	CreatedTime(int64) QueryBuilder
	RequestCode(code ObjectCode) QueryBuilder
	Build() Query
}

type QueryResponseBuilder interface {
	Account(Account) QueryResponseBuilder
	Peer(Peer) QueryResponseBuilder
	Build() QueryResponse
}
