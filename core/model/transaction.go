package model

import "github.com/pkg/errors"

var (
	ErrInvalidTransaction = errors.Errorf("Failed Invalid Transaction")
	ErrTransactionHash    = errors.Errorf("Failed Transaction Hash")
	ErrTransactionVerify  = errors.Errorf("Failed Transaction Verify")
)

type Transaction interface {
	GetPayload() TransactionPayload
	GetSignatures() []Signature
	Modelor
	Sign(PublicKey, PrivateKey) error
	Verify() error
}

type TransactionPayload interface {
	Marshal() ([]byte, error)
	Hash() (Hash, error)
	GetCreatedTime() int64
	GetCommands() []Command
}
