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
	GetBalance() int64
	Modelor
}

type Peer interface {
	GetAddress() string
	GetPublicKey() PublicKey
	Modelor
}

type ObjectList interface {
	GetList() []Object
	Modelor
}

type Object interface {
	GetI32() int32
	GetI64() int64
	GetU32() uint32
	GetUint64() uint64
	GetStr() string
	GetData() []byte
	GetAddress() string
	GetSig() Signature
	GetAccount() Account
	GetPeer() Peer
	GetList() ObjectList
	Modelor
}
