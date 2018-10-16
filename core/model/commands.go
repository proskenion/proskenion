package model

type Command interface {
	GetAuthorizerId() string
	GetTargetId() string

	GetTransfer() Transfer
	GetCreateAccount() CreateAccount
	GetAddAsset() AddAsset
	GetAddPublicKey() AddPublicKey

	Execute(ObjectFinder) error
	Validate(ObjectFinder) error
}

type Transfer interface {
	GetDestAccountId() string
	GetAmount() int64
}

type CreateAccount interface{}

type AddAsset interface {
	GetAmount() int64
}

type AddPublicKey interface {
	GetPublicKey() []byte
}
