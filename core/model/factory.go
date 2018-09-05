package model

import "github.com/pkg/errors"

var (
	ErrNewBlock    = errors.Errorf("Failed Factory NewBlock")
	ErrNewProposal = errors.Errorf("Failed Factory NewProposal")
)

type ModelFactory interface {
	NewBlock(height int64, preBlockHash []byte, createdTime int64, txs []Transaction) (Block, error)
	NewSignature(pubkey []byte, signature []byte) Signature
	NewPeer(address string, pubkey []byte) Peer
}
