package convertor

import (
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
)

type Query struct {
	*proskenion.Query
	cryptor  core.Cryptor
	verifier core.QueryVerifier
}

func (q *Query) GetPayload() model.QueryPayload {
	if q.Query == nil {
		return &QueryPaylaod{}
	}
	return &QueryPaylaod{q.Query.GetPayload(), q.cryptor}
}

func (q *Query) GetSignature() model.Signature {
	if q.Query == nil {
		return &Signature{}
	}
	return &Signature{q.Query.GetSignature()}
}

func (q *Query) Marshal() ([]byte, error) {
	return proto.Marshal(q.Query)
}

func (q *Query) Unmarshal(pb []byte) error {
	return proto.Unmarshal(pb, q.Query)
}

func (q *Query) Hash() model.Hash {
	return q.cryptor.Hash(q)
}

func (q *Query) Sign(pubkey model.PublicKey, privkey model.PrivateKey) error {
	signature, err := q.cryptor.Sign(q.GetPayload(), privkey)
	if err != nil {
		return errors.Wrap(core.ErrCryptorSign, err.Error())
	}
	if q.Query == nil {
		return errors.Errorf("proskenion.Query is nil")
	}
	q.Query.Signature = &proskenion.Signature{
		PublicKey: []byte(pubkey),
		Signature: signature,
	}
	return nil
}

func (q *Query) Verify() error {
	if err := q.cryptor.Verify(q.GetSignature().GetPublicKey(), q.GetPayload(), q.GetSignature().GetSignature()); err != nil {
		return err
	}
	return q.verifier.Verify(q)
}

type QueryPaylaod struct {
	*proskenion.Query_Payload
	cryptor core.Cryptor
}

func (p *QueryPaylaod) GetRequestCode() model.ObjectCode {
	if p.Query_Payload == nil {
		return -1
	}
	return model.ObjectCode(p.RequstCode)
}

func (p *QueryPaylaod) GetWhere() []byte {
	if p.Query_Payload == nil {
		return nil
	}
	ret, _ := proto.Marshal(p.Where)
	return ret
}

func (p *QueryPaylaod) GetOrderBy() model.OrderBy {
	if p.Query_Payload == nil {
		return &OrderBy{&proskenion.Query_OrderBy{}}
	}
	return &OrderBy{p.OrderBy}
}

func (p *QueryPaylaod) Marshal() ([]byte, error) {
	return proto.Marshal(p.Query_Payload)
}

func (q *QueryPaylaod) Unmarshal(pb []byte) error {
	return proto.Unmarshal(pb, q.Query_Payload)
}

func (p *QueryPaylaod) Hash() model.Hash {
	return p.cryptor.Hash(p)
}

type OrderBy struct {
	*proskenion.Query_OrderBy
}

func (o *OrderBy) GetOrder() model.OrderCode {
	return model.OrderCode(o.Order)
}

type QueryResponse struct {
	*proskenion.QueryResponse
	cryptor core.Cryptor
}

func (q *QueryResponse) GetObject() model.Object {
	if q.Object == nil {
		return &Object{}
	}
	return &Object{q.cryptor, q.QueryResponse.GetObject()}
}

func (q *QueryResponse) GetSignature() model.Signature {
	if q.QueryResponse == nil {
		return &Signature{}
	}
	return &Signature{q.QueryResponse.GetSignature()}
}

func (q *QueryResponse) Marshal() ([]byte, error) {
	return proto.Marshal(q.QueryResponse)
}

func (q *QueryResponse) Unmarshal(pb []byte) error {
	return proto.Unmarshal(pb, q.QueryResponse)
}

func (q *QueryResponse) Hash() model.Hash {
	return q.cryptor.Hash(q)
}

func (q *QueryResponse) Sign(pubkey model.PublicKey, privkey model.PrivateKey) error {
	signature, err := q.cryptor.Sign(q.GetObject(), privkey)
	if err != nil {
		return errors.Wrap(core.ErrCryptorSign, err.Error())
	}
	if q.QueryResponse == nil {
		return errors.Errorf("proskenion.Query is nil")
	}
	q.QueryResponse.Signature = &proskenion.Signature{
		PublicKey: []byte(pubkey),
		Signature: signature,
	}
	return nil
}

func (q *QueryResponse) Verify() error {
	return q.cryptor.Verify(q.GetSignature().GetPublicKey(),
		q.GetObject(), q.GetSignature().GetSignature())
}
