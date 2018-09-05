package model

import "github.com/pkg/errors"

var (
	ErrInvalidTransaction = errors.Errorf("Failed Invalid Transaction")
	ErrTransactionGetHash = errors.Errorf("Failed Transaction GetHash")
	ErrTransactionVerify  = errors.Errorf("Failed Transaction Verify")
)

type Transaction interface {
	GetPayload() TransactionPayload
	GetSignatures() []Signature
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	GetHash() ([]byte, error)
	Sign(pubKey []byte, privKey []byte) error
	Verify() error
}

type TransactionPayload interface {
	Marshal() ([]byte, error)
	GetHash() ([]byte, error)
	GetCommands() []Command
}
