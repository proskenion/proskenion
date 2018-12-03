package model

type ObjectCode int

const (
	AnythingObjectCode = iota
	BoolObjectCode
	Int32ObjectCode
	Int64ObjectCode
	Uint32ObjectCode
	Uint64ObjectCode
	StringObjectCode
	BytesObjectCode
	AddressObjectCode
	SignatureObjectCode
	AccountObjectCode
	PeerObjectCode
	ListObjectCode
	DictObjectCode
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

type Object interface {
	GetType() ObjectCode
	GetI32() int32
	GetI64() int64
	GetU32() uint32
	GetU64() uint64
	GetStr() string
	GetData() []byte
	GetAddress() string
	GetSig() Signature
	GetAccount() Account
	GetPeer() Peer
	GetList() []Object
	GetDict() map[string]Object
	Modelor
}
