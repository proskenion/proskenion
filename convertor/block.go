package convertor

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
	"github.com/satellitex/protobuf/proto"
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

func (b *Block) GetFromKey(key string) model.Object {
	switch key {
	case "height":
		return Int64Object(b.GetPayload().GetHeight(), b.cryptor)
	case "pre_block_hash", "pre_block":
		return BytesObject(b.GetPayload().GetPreBlockHash(), b.cryptor)
	case "created_time", "created_at", "time", "at":
		return Int64Object(b.GetPayload().GetCreatedTime(), b.cryptor)
	case "wsv_hash", "wsv":
		return BytesObject(b.GetPayload().GetWSVHash(), b.cryptor)
	case "tx_history_hash", "tx_history":
		return BytesObject(b.GetPayload().GetTxHistoryHash(), b.cryptor)
	case "txs_hash", "txs":
		return BytesObject(b.GetPayload().GetTxsHash(), b.cryptor)
	case "round":
		return Int32Object(b.GetPayload().GetRound(), b.cryptor)
	}
	return &Object{b.cryptor,nil,nil,&proskenion.Object{}}
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
