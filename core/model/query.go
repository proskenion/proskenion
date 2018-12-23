package model

type Query interface {
	GetPayload() QueryPayload
	GetSignature() Signature
	Modelor
	Sign(PublicKey, PrivateKey) error
	Verify() error
	Validate() error
}

type OrderCode int

const (
	DESC = iota
	ASC
)

/**
  string authorizerId = 1;
        string select = 2;
        ObjectCode requstCode = 3;
        string fromId = 4;
        ConditionalFormula where = 5;
        OrderBy orederBy = 6;
        int32 limit = 7;
        int64 createdTime = 8;
*/
type QueryPayload interface {
	GetAuthorizerId() string
	GetSelect() string
	GetRequestCode() ObjectCode
	GetFromId() string
	GetWhere() []byte
	GetOrderBy() OrderBy
	GetLimit() int32
	GetCreatedTime() int64
	Modelor
}

type OrderBy interface {
	GetKey() string
	GetOrder() OrderCode
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
