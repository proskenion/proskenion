package model

import "github.com/pkg/errors"

var (
	ErrInvalidBlock = errors.Errorf("Failed Invalid Block")
	ErrBlockGetHash = errors.Errorf("Failed Block GetHash")
	ErrBlockVerify  = errors.Errorf("Failed Block Verify")
	ErrBlockSign    = errors.Errorf("Failed Block Sign")

	ErrInvalidBlockHeader = errors.Errorf("Failed Invalid BlockHeader")
	ErrBlockHeaderGetHash = errors.Errorf("Failed BlockHeader GetHash")

	ErrInvalidProposal = errors.Errorf("Failed Invalid Proposal")
)

type Block interface {
	GetHeader() BlockHeader
	GetSignature() Signature
	GetHash() (Hash, error)
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Verify() error
	Sign(pubKey PublicKey, privKey PrivateKey) error
}

type BlockHeader interface {
	GetHeight() int64
	GetPreBlockHash() Hash
	GetCreatedTime() int64
	GetCommitTime() int64
	GetMerkleHash() Hash
	GetTxsHash() Hash
	GetRound() int32
	Marshal() ([]byte, error)
	GetHash() (Hash, error)
}
