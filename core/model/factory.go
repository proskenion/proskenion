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

type StorageBuilder interface {
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
	Build() Storage
}

type AccountBuilder interface {
	From(Account) AccountBuilder
	AccountId(string) AccountBuilder
	AccountName(string) AccountBuilder
	PublicKeys([]PublicKey) AccountBuilder
	Quroum(int32) AccountBuilder
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
	CreateAccount(authorizerId string, accountId string) TxBuilder
	AddBalance(accountId string, amount int64) TxBuilder
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
