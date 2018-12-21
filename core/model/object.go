package model

import (
	"fmt"
	"github.com/proskenion/proskenion/regexp"
)

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

type Address interface {
	Storage() string
	Domain() string
	Account() string
	GetBytes() []byte
}

type AddressConv struct {
	storage string
	domain  string
	account string
}

func NewAddress(id string) (Address, error) {
	if regexp.GetRegexp().VerifyWalletId.MatchString(id) ||
		regexp.GetRegexp().VerifyAccountId.MatchString(id) ||
		regexp.GetRegexp().VerifyDomainId.MatchString(id) ||
		regexp.GetRegexp().VerifyStorageId.MatchString(id) {
		ret := regexp.GetRegexp().SplitAddress.FindStringSubmatch(id)
		return &AddressConv{
			ret[3],
			ret[2],
			ret[1],
		}, nil
	}
	return nil, fmt.Errorf("Failed Parse Address not correct format: %s", id)
}

func MustAddress(id string) Address {
	ret, err := NewAddress(id)
	if err != nil {
		panic(err)
	}
	return ret
}

const dividedChar = "\\"

func (a *AddressConv) Storage() string {
	return a.storage
}

func (a *AddressConv) Domain() string {
	return a.domain
}

func (a *AddressConv) Account() string {
	return a.account
}

func (a *AddressConv) GetBytes() []byte {
	ret := make([]byte, 0)
	if a.domain == "" && a.account == "" {
		ret = append(ret, a.storage...)
	} else if a.account == "" {
		ret = append(ret, (a.storage + dividedChar + a.domain)...)
	} else {
		ret = append(ret, (a.storage + dividedChar + a.domain + dividedChar + a.account)...)
	}
	return ret
}
