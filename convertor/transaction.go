package convertor

import (
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
)

var (
	ErrInvalidSignatures = errors.Errorf("Failed Invalid Signatures")
)

type Transaction struct {
	*proskenion.Transaction
	cryptor   core.Cryptor
	executor  core.CommandExecutor
	validator core.CommandValidator
}

type TransactionPayload struct {
	*proskenion.Transaction_Payload
	cryptor   core.Cryptor
	executor  core.CommandExecutor
	validator core.CommandValidator
}

func (t *Transaction) GetPayload() model.TransactionPayload {
	if t.Transaction != nil {
		return &TransactionPayload{t.Payload, t.cryptor, t.executor, t.validator}
	}
	return &TransactionPayload{nil, t.cryptor, t.executor, t.validator}
}

func (t *Transaction) Marshal() ([]byte, error) {
	return proto.Marshal(t.Transaction)
}

func (t *Transaction) Unmarshal(pb []byte) error {
	return proto.Unmarshal(pb, t.Transaction)
}

func (t *Transaction) Hash() (model.Hash, error) {
	if t.Transaction == nil {
		return nil, errors.New("Transaction is nil")
	}
	return t.cryptor.Hash(t)
}

func (t *Transaction) GetSignatures() []model.Signature {
	if t.Transaction == nil {
		return []model.Signature{}
	}
	ret := make([]model.Signature, len(t.Signatures))
	for i, sig := range t.Transaction.GetSignatures() {
		ret[i] = &Signature{sig}
	}
	return ret
}

func (t *Transaction) Sign(pubkey model.PublicKey, privkey model.PrivateKey) error {
	signature, err := t.cryptor.Sign(t.GetPayload(), privkey)
	if err != nil {
		return err
	}
	if t.Transaction == nil {
		return errors.Errorf("proskenion.Transaction is nil")
	}
	t.Transaction.Signatures = append(t.Transaction.Signatures, &proskenion.Signature{PublicKey: []byte(pubkey), Signature: signature})
	return nil
}

func (t *Transaction) Verify() error {
	if len(t.GetSignatures()) == 0 {
		return errors.Wrapf(ErrInvalidSignatures, "Signatures length is 0")
	}
	for i, signature := range t.GetSignatures() {
		if signature == nil {
			return errors.Wrapf(model.ErrInvalidSignature, "%d-th Signature is nil", i)
		}
		if err := t.cryptor.Verify(signature.GetPublicKey(), t.GetPayload(), signature.GetSignature()); err != nil {
			return errors.Wrapf(core.ErrCryptorVerify, err.Error())
		}
	}
	return nil
}

func (t *TransactionPayload) Marshal() ([]byte, error) {
	return proto.Marshal(t.Transaction_Payload)
}

func (t *TransactionPayload) Hash() (model.Hash, error) {
	return t.cryptor.Hash(t)
}

func (t *TransactionPayload) GetCommands() []model.Command {
	if t == nil {
		return []model.Command{}
	}
	ret := make([]model.Command, len(t.Commands))
	for i, cmd := range t.Transaction_Payload.GetCommands() {
		ret[i] = &Command{cmd, t.executor, t.validator}
	}
	return ret
}
