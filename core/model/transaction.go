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
	Validate(ObjectFinder, TxFinder) error
}

type TransactionPayload interface {
	GetCreatedTime() int64
	GetCommands() []Command
	Modelor
}
