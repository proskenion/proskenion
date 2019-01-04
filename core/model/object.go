package model

import (
	"bytes"
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
	StorageObjectCode
)

type Account interface {
	GetAccountId() string
	GetAccountName() string
	GetPublicKeys() []PublicKey
	GetBalance() int64
	GetQuorum() int32
	GetDelegatePeerId() string
	Modelor
}

type Peer interface {
	GetPeerId() string
	GetAddress() string
	GetPublicKey() PublicKey
	Modelor
}

func HasherLess(a, b Hasher) bool {
	return bytes.Compare(a.Hash(), b.Hash()) > 0
}

func HasherEqual(a, b Hasher) bool {
	return bytes.Equal(a.Hash(), b.Hash())
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
	GetStorage() Storage
	Modelor
}

// a < b
func ObjectLess(a, b Object) bool {
	switch a.GetType() {
	case Int32ObjectCode:
		return a.GetI32() < b.GetI32()
	case Int64ObjectCode:
		return a.GetI64() < b.GetI64()
	case Uint32ObjectCode:
		return a.GetU32() < b.GetU32()
	case Uint64ObjectCode:
		return a.GetU64() < b.GetU64()
	case StringObjectCode:
		return a.GetStr() < b.GetStr()
	case BytesObjectCode:
		return bytes.Compare(a.GetData(), b.GetData()) > 0
	case AddressObjectCode:
		return a.GetAddress() < b.GetAddress()
	case SignatureObjectCode:
		if bytes.Equal(a.GetSig().GetPublicKey(), a.GetSig().GetPublicKey()) {
			return bytes.Compare(a.GetSig().GetSignature(), b.GetSig().GetSignature()) > 0
		}
		return bytes.Compare(a.GetSig().GetPublicKey(), b.GetSig().GetPublicKey()) > 0
	case AccountObjectCode:
		return HasherLess(a.GetAccount(), b.GetAccount())
	case PeerObjectCode:
		return HasherLess(a.GetPeer(), b.GetPeer())
	case ListObjectCode:
		for i := 0; i < len(a.GetList()) && i < len(b.GetList()); i++ {
			if HasherEqual(a.GetList()[i], b.GetList()[i]) {
				continue
			}
			return ObjectLess(a.GetList()[i], b.GetList()[i])
		}
		return len(a.GetList()) < len(b.GetList())
	case DictObjectCode:
		for k, v := range a.GetDict() {
			if v2, ok := b.GetDict()[k]; ok {
				if HasherEqual(a.GetDict()[k], a.GetDict()[k]) {
					continue
				}
				return ObjectLess(v, v2)
			} else {
				return len(a.GetDict()) < len(b.GetDict())
			}
		}
	case StorageObjectCode:
		return HasherLess(a.GetStorage(), b.GetStorage())
	}
	return true
}

const (
	WallettAddressType = iota
	AccountAddressType
	DomainAddressType
	StorageAddressType
)

type Address interface {
	Storage() string
	Domain() string
	Account() string
	GetBytes() []byte
	Type() int
	Id() string
	AccountId() string
	PeerId() string
}

type AddressConv struct {
	storage string
	domain  string
	account string
	t       int
}

func NewAddress(id string) (Address, error) {
	t := -1
	if regexp.GetRegexp().VerifyWalletId.MatchString(id) {
		t = WallettAddressType
	} else if regexp.GetRegexp().VerifyAccountId.MatchString(id) {
		t = AccountAddressType
	} else if regexp.GetRegexp().VerifyDomainId.MatchString(id) {
		t = DomainAddressType
	} else if regexp.GetRegexp().VerifyStorageId.MatchString(id) {
		t = StorageAddressType
	}
	if t != -1 {
		ret := regexp.GetRegexp().SplitAddress.FindStringSubmatch(id)
		return &AddressConv{
			ret[3],
			ret[2],
			ret[1],
			t,
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

func (a *AddressConv) Type() int {
	return a.t
}

func (a *AddressConv) Id() string {
	if a.account == "" {
		return a.domain + "/" + a.storage
	}
	if a.storage == "" {
		return a.account + "@" + a.domain
	}
	return a.account + "@" + a.domain + "/" + a.storage
}

func (a *AddressConv) AccountId() string {
	return a.account + "@" + a.domain + "/account"
}

func (a *AddressConv) PeerId() string {
	return a.account + "@" + a.domain + "/peer"
}
