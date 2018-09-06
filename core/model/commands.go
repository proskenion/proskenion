package model

type Command interface {
	GetTransfer() Transfer
}

type Transfer interface {
	GetSrcAccountId() string
	GetDestAccountId() string
	GetAmount() int64
	Execute() error
	Validate() error
}
