package model

import (
	"github.com/pkg/errors"
)

var (
	ErrNewBlock    = errors.Errorf("Failed Factory NewBlock")
	ErrNewProposal = errors.Errorf("Failed Factory NewProposal")
)

type ModelFactory interface {
	NewSignature(pubkey PublicKey, signature []byte) Signature
	NewAccount(accountId string, accountName string, publicKeys []PublicKey, amount int64) Account
	NewPeer(address string, pubkey PublicKey) Peer
	NewBlockBuilder() BlockBuilder
	NewTxBuilder() TxBuilder
	NewQueryBuilder() QueryBuilder
	NewQueryResponseBuilder() QueryResponseBuilder

	NewEmptyBlock() Block
	NewEmptyAccount() Account
	NewEmptyPeer() Peer
	NewEmptyTx() Transaction
	NewEmptyQuery() Query
	NewEmptyQueryResponse() QueryResponse
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
