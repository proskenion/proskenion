package model

import "github.com/pkg/errors"

var (
	ErrInvalidBlock = errors.Errorf("Failed Invalid Block")
	ErrBlockHash = errors.Errorf("Failed Block Hash")
	ErrBlockVerify  = errors.Errorf("Failed Block Verify")
	ErrBlockSign    = errors.Errorf("Failed Block Sign")

	ErrInvalidBlockHeader = errors.Errorf("Failed Invalid BlockHeader")
	ErrBlockHeaderHash = errors.Errorf("Failed BlockHeader Hash")

	ErrInvalidProposal = errors.Errorf("Failed Invalid Proposal")
)

type Block interface {
	GetPayload() BlockPayload
	GetSignature() Signature
	Hash() (Hash, error)
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Verify() error
	Sign(PublicKey, PrivateKey) error
}

type BlockPayload interface {
	GetHeight() int64
	GetPreBlockHash() Hash
	GetCreatedTime() int64
	GetMerkleHash() Hash
	GetTxsHash() Hash
	GetRound() int32
	Marshal() ([]byte, error)
	Hash() (Hash, error)
}
