package model

import "github.com/pkg/errors"

var (
	ErrInvalidBlock = errors.Errorf("Failed Invalid Block")
	ErrBlockHash    = errors.Errorf("Failed Block Hash")
	ErrBlockVerify  = errors.Errorf("Failed Block Verify")
	ErrBlockSign    = errors.Errorf("Failed Block Sign")

	ErrInvalidBlockHeader = errors.Errorf("Failed Invalid BlockHeader")
	ErrBlockHeaderHash    = errors.Errorf("Failed BlockHeader Hash")

	ErrInvalidProposal = errors.Errorf("Failed Invalid Proposal")
)

type Block interface {
	GetPayload() BlockPayload
	GetSignature() Signature
	Modelor
	Verify() error
	Sign(PublicKey, PrivateKey) error
}

type BlockPayload interface {
	GetHeight() int64
	GetPreBlockHash() Hash
	GetCreatedTime() int64
	GetWSVHash() Hash
	GetTxHistoryHash() Hash
	GetTxsHash() Hash
	GetRound() int32
	Modelor
}
