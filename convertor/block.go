package convertor

import (
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
)

type Block struct {
	*proskenion.Block
	cryptor core.Cryptor
}

func (b *Block) GetPayload() model.BlockPayload {
	if b.Block != nil {
		return &BlockPayload{b.Payload, b.cryptor}
	}
	return &BlockPayload{nil, b.cryptor}
}

func (b *Block) GetSignature() model.Signature {
	if b.Block != nil {
		return &Signature{b.Signature}
	}
	return &Signature{nil}
}

func (b *Block) Marshal() ([]byte, error) {
	return proto.Marshal(b.Block)
}

func (b *Block) Unmarshal(pb []byte) error {
	return proto.Unmarshal(pb, b.Block)
}

func (b *Block) Hash() model.Hash {
	return b.cryptor.Hash(b)
}

func (b *Block) Verify() error {
	return b.cryptor.Verify(b.GetSignature().GetPublicKey(), b.GetPayload(), b.GetSignature().GetSignature())
}

func (b *Block) Sign(pubKey model.PublicKey, privKey model.PrivateKey) error {
	signature, err := b.cryptor.Sign(b.GetPayload(), privKey)
	if err != nil {
		return errors.Wrapf(core.ErrCryptorSign, err.Error())
	}
	b.Signature = &proskenion.Signature{PublicKey: pubKey, Signature: signature}
	return nil
}

type BlockPayload struct {
	*proskenion.Block_Payload
	cryptor core.Cryptor
}

func (p *BlockPayload) Marshal() ([]byte, error) {
	return proto.Marshal(p.Block_Payload)
}

func (p *BlockPayload) Unmarshal(pb []byte) error {
	return proto.Unmarshal(pb, p.Block_Payload)
}

func (p *BlockPayload) Hash() model.Hash {
	return p.cryptor.Hash(p)
}

func (p *BlockPayload) GetPreBlockHash() model.Hash {
	if p.Block_Payload == nil {
		return nil
	}
	return p.Block_Payload.GetPreBlockHash()
}

func (p *BlockPayload) GetWSVHash() model.Hash {
	if p.Block_Payload == nil {
		return nil
	}
	return p.Block_Payload.GetWsvHash()
}

func (p *BlockPayload) GetTxHistoryHash() model.Hash {
	if p.Block_Payload == nil {
		return nil
	}
	return p.Block_Payload.GetTxHistoryHash()
}

func (p *BlockPayload) GetTxsHash() model.Hash {
	if p.Block_Payload == nil {
		return nil
	}
	return p.Block_Payload.GetTxsHash()
}
