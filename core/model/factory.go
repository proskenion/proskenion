package model

import "github.com/pkg/errors"

var (
	ErrNewBlock    = errors.Errorf("Failed Factory NewBlock")
	ErrNewProposal = errors.Errorf("Failed Factory NewProposal")
)

type ModelFactory interface {
	NewBlock(height int64, preBlockHash Hash, createdTime int64, merkleHash Hash, txsHash Hash, round int32) Block
	NewSignature(pubkey PublicKey, signature []byte) Signature
	NewPeer(address string, pubkey PublicKey) Peer
	NewTxBuilder() TxBuilder
}

type TxBuilder interface {
	CreatedTime(int64) TxBuilder
	Transfer(srcAccountId string, destAccountId string, amount int64) TxBuilder
	Build() Transaction
}
