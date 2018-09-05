package model

type ObjectCode int

const (
	AccountObjectCode ObjectCode = iota
	PeerObjectCode
)

type Account interface {
	GetAccountId() string
	GetAccountName() string
	GetPublicKeys() []PublicKey
	GetAmount() int64
}

type Peer interface {
	GetAddress() string
	GetPublicKey() PublicKey
}
