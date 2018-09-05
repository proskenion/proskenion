package model

type Query interface {
	GetPayload() QueryPayload
	GetSignature() Signature
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	GetHash() ([]byte, error)
	Sign(pubKey []byte, privKey []byte) error
	Verify() error
}

type QueryPayload interface {
	GetAuthorizerId() string
	GetTargetId() string
	GetCreatedTime() int64
	GetRequest() ObjectCode
}

type QueryResponse interface {
	GetPayload() QueryResponsePayload
	GetSignature() Signature
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	GetHash() ([]byte, error)
	Sign(pubKey []byte, privKey []byte) error
	Verify() error
}

type QueryResponsePayload interface {
	ResponseCode() ObjectCode
	GetAccount() Account
	Marshal() ([]byte, error)
	GetHash() ([]byte, error)
	Verify() error
}
