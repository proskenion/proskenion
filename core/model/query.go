package model

type Query interface {
	GetPayload() QueryPayload
	GetSignature() Signature
	Modelor
	Sign(PublicKey, PrivateKey) error
	Verify() error
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
	GetWhere() string
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
	GetObject() Object
	GetSignature() Signature
	Modelor
	Sign(PublicKey, PrivateKey) error
	Verify() error
}
