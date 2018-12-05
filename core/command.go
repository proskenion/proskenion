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

// TransferBalance Err
var (
	ErrCommandExecutorTransferBalanceNotFoundSrcAccountId      = fmt.Errorf("Failed Command Executor TransferBalance Can Not Load SrcAccounId")
	ErrCommandExecutorTransferBalanceNotFoundDestAccountId     = fmt.Errorf("Failed Command Executor TransferBalance Can Not Load DestAccounId")
	ErrCommandExecutorTransferBalanceNotEnoughSrcAccountBalance = fmt.Errorf("Failed Command Executor TransferBalance Not Enough SrcAccount Balance")
)

// CreateAccount Err
var (
	ErrCommandExecutorCreateAccountAlreadyExistAccount = fmt.Errorf("Failed Command Executor CreateAccount AlreadyExist AccountId")
)

// AddBalance Err
var (
	ErrCommandExecutorAddBalanceNotExistAccount = fmt.Errorf("Failed Command Executor AddBalance Not Exist Account")
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
	TransferBalance(ObjectFinder, Command) error
	CreateAccount(ObjectFinder, Command) error
	AddBalance(ObjectFinder, Command) error
	AddPublicKeys(ObjectFinder, Command) error
}

type CommandValidator interface {
	SetFactory(factory ModelFactory)
	TransferBalance(ObjectFinder, Command) error
	CreateAccount(ObjectFinder, Command) error
	AddBalance(ObjectFinder, Command) error
	AddPublicKeys(ObjectFinder, Command) error
	Tx(ObjectFinder, TxFinder, Transaction) error
}
