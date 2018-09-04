package model

type Executor interface {
	Execute() error
}

type Validator interface {
	Validate() error
}

type Command interface {
	GetTransfer() Transaction
}

type Transfer interface {
	SrcAccount() string
	DestAccount() string
	Amount() int64
	Execute() error
	Validate() error
}
