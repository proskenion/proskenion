package model

type Command interface {
	GetTransfer() Transfer
	Execute(ObjectFinder) error
	Validate(ObjectFinder) error
}

type Transfer interface {
	GetSrcAccountId() string
	GetDestAccountId() string
	GetAmount() int64
	Execute(ObjectFinder) error
	Validate(ObjectFinder) error
}

type CreateAccount interface {
	GetAccountId() string
	Execute(ObjectFinder) error
	Validate(ObjectFinder) error
}

type AddAsset interface {
	GetAccountId() string
	GetAmount() int64
	Execute(ObjectFinder) error
	Validate(ObjectFinder) error
}
