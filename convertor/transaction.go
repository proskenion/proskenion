package convertor

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/model"
	"github.com/satellitex/bbft/proto"
	"github.com/proskenion/proskenion/core"
)

var (
	ErrInvalidSignatures = errors.Errorf("Failed Invalid Signatures")
)

type Transaction struct {
	*bbft.Transaction
	cryptor core.Cryptor
}

type TransactionPayload struct {
	*bbft.Transaction_Payload
	executor core.CommandExecutor
	validator core.CommandValidator
}

func (t *Transaction) GetPayload() model.TransactionPayload {
	if t.Transaction == nil {
		return &TransactionPayload{nil}
	}
	return &TransactionPayload{t.Payload}
}


func (t *Transaction) GetHash() ([]byte, error) {
	if t.Transaction == nil {
		return nil, errors.New("*bbft.Transaction is nil")
	}
	return t.cryptor.Hash(t)
}

func (t *Transaction) GetSignatures() []model.Signature {
	ret := make([]model.Signature, len(t.Signatures))
	for i, sig := range t.Signatures {
		ret[i] = &Signature{sig}
	}
	return ret
}

func (t *Transaction) Verify() error {
	hash, err := t.GetHash()
	if err != nil {
		return errors.Wrapf(model.ErrTransactionGetHash, err.Error())
	}
	if len(t.GetSignatures()) == 0 {
		return errors.Wrapf(ErrInvalidSignatures, "Signatures length is 0")
	}
	for i, signature := range t.Signatures {
		if signature == nil {
			return errors.Wrapf(model.ErrInvalidSignature, "%d-th Signature is nil", i)
		}
		if err := Verify(signature.Pubkey, hash, signature.Signature); err != nil {
			return errors.Wrapf(ErrCryptoVerify, err.Error())
		}
	}
	return nil
}

func (p *TransactionPayload) GetMessage() string {
	return p.Todo
}
