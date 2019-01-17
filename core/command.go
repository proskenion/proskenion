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
	ErrCommandExecutorTransferBalanceNotFoundSrcAccountId       = fmt.Errorf("Failed Command Executor TransferBalance Can Not Load SrcAccounId")
	ErrCommandExecutorTransferBalanceNotFoundDestAccountId      = fmt.Errorf("Failed Command Executor TransferBalance Can Not Load DestAccounId")
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

// CreateStorage Err
var (
	ErrCommandExecutorCreateStorageNotDefinedStorage = fmt.Errorf("Failed Command Executor CreateStorage Not Defined Storage")
)

// UpdateObject Err
var (
	ErrCommandExecutorUpdateObjectNotExistWallet = fmt.Errorf("Failed Command Executor UpdateObject Not Exist Wallet")
)

// AddObject Err
var (
	ErrCommandExecutorAddObjectNotExistWallet = fmt.Errorf("Failed Command Executor AddObject Not Exist Wallet")
)

// TransferObject Err
var (
	ErrCommandExecutorTransferObjectNotExistSrcWallet  = fmt.Errorf("Failed Command Executor TransferObject Not Exist Source Wallet")
	ErrCommandExecutorTransferObjectNotExistDestWallet = fmt.Errorf("Failed Command Executor TransferObject Not Exist Dest Wallet")
)

// Consign Err
var (
	ErrCommandExecutorConsignNotFoundAccount = fmt.Errorf("Failed Command Executor Consign Not Found Account")
)

// CheckAndCommitProsl Err
var (
	ErrCommandExecutorCheckAndCommitProslInvalid  = fmt.Errorf("Failed Check And Commit Prosl invalid change rule: false")
	ErrCommandExecutorCheckAndCommitProslNotFound = fmt.Errorf("Failed Check And Commit Prosl not found target prosl")
)

// Transaction Err
var (
	ErrTxValidateNotFoundAuthorizer  = fmt.Errorf("Failed Transaction Validator Authorizer Not Found")
	ErrTxValidateNotSignedAuthorizer = fmt.Errorf("Failed Transaction Validator Authorizer's not signed")
	ErrTxValidateAlreadyExist        = fmt.Errorf("Failed Transaction Validator Already Exists")
)

const (
	TargetIdKey   = "target_id"
	ProslKey      = "prosl"
	ProslTypeKey  = "prosl_type"
	IncentiveKey  = "incentive"
	ConsensusKey  = "consensus"
	ChangeRuleLey = "rule"
)

type CommandExecutor interface {
	SetField(factory ModelFactory, prosl Prosl)
	TransferBalance(ObjectFinder, Command) error
	CreateAccount(ObjectFinder, Command) error
	AddBalance(ObjectFinder, Command) error
	AddPublicKeys(ObjectFinder, Command) error
	DefineStorage(ObjectFinder, Command) error
	CreateStorage(ObjectFinder, Command) error
	UpdateObject(ObjectFinder, Command) error
	AddObject(ObjectFinder, Command) error
	TransferObject(ObjectFinder, Command) error
	AddPeer(ObjectFinder, Command) error
	Consign(ObjectFinder, Command) error
	CheckAndCommitProsl(ObjectFinder, Command) error
}

type CommandValidator interface {
	SetField(factory ModelFactory, prosl Prosl)
	TransferBalance(ObjectFinder, Command) error
	CreateAccount(ObjectFinder, Command) error
	AddBalance(ObjectFinder, Command) error
	AddPublicKeys(ObjectFinder, Command) error
	DefineStorage(ObjectFinder, Command) error
	CreateStorage(ObjectFinder, Command) error
	UpdateObject(ObjectFinder, Command) error
	AddObject(ObjectFinder, Command) error
	TransferObject(ObjectFinder, Command) error
	AddPeer(ObjectFinder, Command) error
	Consign(ObjectFinder, Command) error
	Tx(ObjectFinder, TxFinder, Transaction) error
	CheckAndCommitProsl(ObjectFinder, Command) error
}
