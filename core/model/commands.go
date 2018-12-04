package model

type Command interface {
	GetAuthorizerId() string
	GetTargetId() string

	GetTransferBalance() TransferBalance
	GetCreateAccount() CreateAccount
	GetAddBalance() AddBalance
	GetAddPublicKeys() AddPublicKeys

	Execute(ObjectFinder) error
	Validate(ObjectFinder) error
}

type TransferBalance interface {
	GetDestAccountId() string
	GetBalance() int64
}

type CreateAccount interface{}

type AddBalance interface {
	GetBalance() int64
}

type AddPublicKeys interface {
	GetPublicKeys() [][]byte
}
