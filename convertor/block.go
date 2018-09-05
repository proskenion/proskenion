package convertor

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
)

type Block struct {
	*proskenion.Block
}

type BlockPayload struct {
	*proskenion.Block_Payload
}

func (b *Block) GetPayload() model.BlockPayload {
	if b.Block != nil {
		return &BlockPayload{b.Payload}
	}
	return &BlockPayload{nil}
}

func (b *Block) GetTransactions() []model.Transaction {
	ret := make([]model.Transaction, len(b.Transactions))
	for id, tx := range b.Transactions {
		ret[id] = &Transaction{tx}
	}
	return ret
}

func (b *Block) GetSignature() model.Signature {
	if b.Block != nil {
		return &Signature{b.Signature}
	}
	return &Signature{nil}
}

func (b *Block) GetHash() ([]byte, error) {
	//TODO 毎回 sha256計算したほうが一気にやるよりはやそう？
	result, err := b.GetPayload().GetHash()
	if err != nil {
		return nil, errors.Wrapf(model.ErrBlockPayloadGetHash, err.Error())
	}
	for _, tx := range b.GetTransactions() {
		// Unexpeted Can not cast Transaction of GetTransactions()
		hash, err := CalcHashFromProto(tx.(*Transaction))
		if err != nil {
			return nil, errors.Wrapf(model.ErrTransactionGetHash, err.Error())
		}
		result = append(result, hash...)
	}
	return CalcHash(result), nil
}

func (b *Block) Verify() error {
	hash, err := b.GetHash()
	if err != nil {
		return errors.Wrapf(model.ErrBlockGetHash, err.Error())
	}
	if b.Signature == nil {
		return errors.Wrapf(model.ErrInvalidSignature, "Signature is nil")
	}
	if err = Verify(b.Signature.Pubkey, hash, b.Signature.Signature); err != nil {
		return errors.Wrapf(ErrCryptoVerify, err.Error())
	}
	return nil
}

func (b *Block) Sign(pubKey []byte, privKey []byte) error {
	hash, err := b.GetHash()
	if err != nil {
		return errors.Wrapf(model.ErrBlockGetHash, err.Error())
	}
	signature, err := Sign(privKey, hash)
	if err != nil {
		return errors.Wrapf(ErrCryptoSign, err.Error())
	}
	if err := Verify(pubKey, hash, signature); err != nil {
		return errors.Wrapf(ErrCryptoVerify, err.Error())
	}
	b.Signature = &proskenion.Signature{Pubkey: pubKey, Signature: signature}
	return nil
}

func (h *BlockPayload) GetHash() ([]byte, error) {
	return CalcHashFromProto(h)
}

func (p *Proposal) GetBlock() model.Block {
	return &Block{p.Block}
}
