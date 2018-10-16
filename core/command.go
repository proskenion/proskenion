package core

import (
	"fmt"
	. "github.com/proskenion/proskenion/core/model"
)

type Executor interface {
	Execute(ObjectFinder) error
}

type Validator interface {
	Validate(ObjectFinder) error
}

// Transfer Err
var (
	ErrCommandExecutorTransferNotFoundSrcAccountId      = fmt.Errorf("Failed Command Executor Transfer Can Not Load SrcAccounId")
	ErrCommandExecutorTransferNotFoundDestAccountId     = fmt.Errorf("Failed Command Executor Transfer Can Not Load DestAccounId")
	ErrCommandExecutorTransferNotEnoughSrcAccountAmount = fmt.Errorf("Failed Command Executor Transfer Not Enough SrcAccount Amount")
)

// CreateAccount Err
var (
	ErrCommandExecutorCreateAccountAlreadyExistAccount = fmt.Errorf("Failed Command Executor CreateAccount AlreadyExist AccountId")
)

// AddAsset Err
var (
	ErrCommandExecutorAddAssetNotExistAccount = fmt.Errorf("Failed Command Executor AddAsset Not Exist Account")
)

// AddPublicKeys Err
var (
	ErrCommandExecutorAddPublicKeyNotExistAccount = fmt.Errorf("Failed Command Executor AddPublicKey Not Exist Account")
	ErrCommandExecutorAddPublicKeyDuplicatePubkey = fmt.Errorf("Failed Command Executor AddPublicKey Duplicate Add PublicKey")
)

// Transaction Err
var (
	ErrTxValidateNotFoundAuthorizer  = fmt.Errorf("Failed Transaction Validator Authorizer Not Found")
	ErrTxValidateNotSignedAuthorizer = fmt.Errorf("Failed Transaction Validator Authorizer's not signed")
	ErrTxValidateAlreadyExist        = fmt.Errorf("Failed Transaction Validator Already Exists")
)

type CommandExecutor interface {
	SetFactory(factory ModelFactory)
	Transfer(ObjectFinder, Command) error
	CreateAccount(ObjectFinder, Command) error
	AddAsset(ObjectFinder, Command) error
	AddPublicKey(ObjectFinder, Command) error
}

type CommandValidator interface {
	SetFactory(factory ModelFactory)
	Transfer(ObjectFinder, Command) error
	CreateAccount(ObjectFinder, Command) error
	AddAsset(ObjectFinder, Command) error
	AddPublicKey(ObjectFinder, Command) error
	Tx(ObjectFinder, TxFinder, Transaction) error
}
