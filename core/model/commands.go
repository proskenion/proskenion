package model

type Command interface {
	GetAuthorizerId() string
	GetTargetId() string

	GetTransferBalance() TransferBalance
	GetCreateAccount() CreateAccount
	GetAddBalance() AddBalance
	GetAddPublicKeys() AddPublicKeys
	GetRemovePublicKeys() RemovePublicKeys
	GetSetQuorum() SetQuroum
	GetDefineStorage() DefineStorage
	GetCreateStorage() CreateStorage
	GetUpdateObject() UpdateObject
	GetAddObject() AddObject
	GetTransferObject() TransferObject
	GetAddPeer() AddPeer
	GetConsign() Consign

	Execute(ObjectFinder) error
	Validate(ObjectFinder) error
}

type TransferBalance interface {
	GetDestAccountId() string
	GetBalance() int64
}

type CreateAccount interface {
	GetPublicKeys() []PublicKey
	GetQuorum() int32
}

type AddBalance interface {
	GetBalance() int64
}

type AddPublicKeys interface {
	GetPublicKeys() []PublicKey
}

type RemovePublicKeys interface {
	GetPublicKeys() []PublicKey
}

type SetQuroum interface {
	GetQuorum() int32
}

type DefineStorage interface {
	GetStorage() Storage
}

type CreateStorage interface {
}

type UpdateObject interface {
	GetKey() string
	GetObject() Object
}

type AddObject interface {
	GetKey() string
	GetObject() Object
}

type TransferObject interface {
	GetKey() string
	GetDestAccountId() string
	GetObject() Object
}

type AddPeer interface {
	GetAddress() string
	GetPublicKey() []byte
}

type Consign interface {
	GetPeerId() string
}
