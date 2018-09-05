package model

type ObjectCode int

const (
	AccountObjectCode ObjectCode = iota
)

type Account interface {
	GetAccountId() string
	GetAccountName() string
	GetPublicKeys() []PublicKey
	GetAmount() int64
}
