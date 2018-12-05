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

	Execute(ObjectFinder) error
	Validate(ObjectFinder) error
}

type TransferBalance interface {
	GetDestAccountId() string
	GetBalance() int64
}

type CreateAccount interface{}

type AddBalance interface {
	GetBalance() int64
}

type AddPublicKeys interface {
	GetPublicKeys() [][]byte
}

type RemovePublicKeys interface {
	GetPublicKeys() [][]byte
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
	GetObject() Object
}

type AddObject interface {
	GetObject() Object
}

type TransferObject interface {
	GetDestAccountId() string
}

type AddPeer interface {
	GetAddress() string
	GetPublicKey() []byte
}
