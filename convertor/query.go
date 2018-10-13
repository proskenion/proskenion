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
	cryptor   core.Cryptor
	validator core.QueryValidator
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

func (q *Query) Hash() (model.Hash, error) {
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
	return q.cryptor.Verify(q.GetSignature().GetPublicKey(), q.GetPayload(), q.GetSignature().GetSignature())
}

func (q *Query) Validate() error {
	return q.validator.Validate(q)
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

func (p *QueryPaylaod) Marshal() ([]byte, error) {
	return proto.Marshal(p.Query_Payload)
}

func (p *QueryPaylaod) Hash() (model.Hash, error) {
	return p.cryptor.Hash(p)
}

type QueryResponse struct {
	*proskenion.QueryResponse
	cryptor core.Cryptor
}

func (q *QueryResponse) GetPayload() model.QueryResponsePayload {
	if q.QueryResponse == nil {
		return &QueryResponsePayload{}
	}
	return &QueryResponsePayload{q.QueryResponse.GetPayload(), q.cryptor}
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

func (q *QueryResponse) Hash() (model.Hash, error) {
	return q.cryptor.Hash(q)
}

func (q *QueryResponse) Sign(pubkey model.PublicKey, privkey model.PrivateKey) error {
	signature, err := q.cryptor.Sign(q.GetPayload(), privkey)
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
		q.GetPayload(), q.GetSignature().GetSignature())
}

type QueryResponsePayload struct {
	*proskenion.QueryResponse_Payload
	cryptor core.Cryptor
}

func (p *QueryResponsePayload) ResponseCode() model.ObjectCode {
	return model.ObjectCode(p.ResponseCode())
}

func (p *QueryResponsePayload) GetAccount() model.Account {
	if p.QueryResponse_Payload == nil ||
		p.QueryResponse_Payload.GetAccount() == nil {
		return &Account{}
	}
	return &Account{p.cryptor, p.QueryResponse_Payload.GetAccount()}
}

func (p *QueryResponsePayload) GetPeer() model.Peer {
	if p.QueryResponse_Payload == nil ||
		p.QueryResponse_Payload.GetPeer() == nil {
		return &Peer{}
	}
	return &Peer{p.cryptor, p.QueryResponse_Payload.GetPeer()}
}

func (p *QueryResponsePayload) Marshal() ([]byte, error) {
	return proto.Marshal(p.QueryResponse_Payload)
}

func (p *QueryResponsePayload) Hash() (model.Hash, error) {
	return p.cryptor.Hash(p)
}
