package model

type Query interface {
	GetPayload() QueryPayload
	GetSignature() Signature
	Modelor
	Sign(PublicKey, PrivateKey) error
	Verify() error
	Validate() error
}

type QueryPayload interface {
	GetAuthorizerId() string
	GetTargetId() string
	GetCreatedTime() int64
	GetRequestCode() ObjectCode
	Marshal() ([]byte, error)
	Hash() (Hash, error)
}

type QueryResponse interface {
	GetPayload() QueryResponsePayload
	GetSignature() Signature
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Hash() (Hash, error)
	Sign(PublicKey, PrivateKey) error
	Verify() error
}

type QueryResponsePayload interface {
	ResponseCode() ObjectCode
	GetAccount() Account
	Marshal() ([]byte, error)
	Hash() (Hash, error)
}
