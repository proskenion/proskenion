package model

import "github.com/pkg/errors"

var (
	ErrNewBlock    = errors.Errorf("Failed Factory NewBlock")
	ErrNewProposal = errors.Errorf("Failed Factory NewProposal")
)

type ModelFactory interface {
	NewBlock(height int64, preBlockHash Hash, createdTime int64, merkleHash Hash, txsHash Hash, round int32) Block
	NewSignature(pubkey PublicKey, signature []byte) Signature
	NewAccount(accountId string, accountName string, publicKeys []PublicKey, amount int64) Account
	NewPeer(address string, pubkey PublicKey) Peer
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

type TxBuilder interface {
	CreatedTime(int64) TxBuilder
	Transfer(srcAccountId string, destAccountId string, amount int64) TxBuilder
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
	Build() QueryResponse
}
