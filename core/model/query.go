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
	Modelor
}

type QueryResponse interface {
	GetPayload() QueryResponsePayload
	GetSignature() Signature
	Modelor
	Sign(PublicKey, PrivateKey) error
	Verify() error
}

type QueryResponsePayload interface {
	ResponseCode() ObjectCode
	GetAccount() Account
	GetPeer() Peer
	Modelor
}
