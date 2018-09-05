package model

import "github.com/pkg/errors"

var (
	ErrInvalidTransaction = errors.Errorf("Failed Invalid Transaction")
	ErrTransactionHash = errors.Errorf("Failed Transaction Hash")
	ErrTransactionVerify  = errors.Errorf("Failed Transaction Verify")
)

type Transaction interface {
	GetPayload() TransactionPayload
	GetSignatures() []Signature
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Hash() ([]byte, error)
	Sign(pubKey []byte, privKey []byte) error
	Verify() error
}

type TransactionPayload interface {
	Marshal() ([]byte, error)
	Hash() ([]byte, error)
	GetCommands() []Command
}
