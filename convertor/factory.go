package convertor

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
)

type ModelFactory struct {
	cryptor   core.Cryptor
	executor  core.CommandExecutor
	validator core.CommandValidator
}

func NewModelFactory(cryptor core.Cryptor,
	executor core.CommandExecutor,
	validator core.CommandValidator) model.ModelFactory {
	return &ModelFactory{cryptor, executor, validator}
}

func (f *ModelFactory) NewBlock(height int64,
	preBlockHash model.Hash, createdTime int64,
	merkleHash model.Hash, txsHash model.Hash,
	round int32) model.Block {
	return &Block{
		&proskenion.Block{
			Payload: &proskenion.Block_Payload{
				Height:       height,
				PreBlockHash: preBlockHash,
				CreatedTime:  createdTime,
				MerkleHash:   merkleHash,
				TxsHash:      txsHash,
				Round:        round,
			},
			Signature: &proskenion.Signature{},
		},
		f.cryptor,
	}
}

func (f *ModelFactory) NewSignature(pubkey model.PublicKey, signature []byte) model.Signature {
	return &Signature{
		&proskenion.Signature{
			PublicKey: []byte(pubkey),
			Signature: signature,
		},
	}
}

func (f *ModelFactory) NewPeer(address string, pubkey model.PublicKey) model.Peer {
	return &Peer{
		&proskenion.Peer{
			Address:   address,
			PublicKey: []byte(pubkey),
		},
	}
}

func (f *ModelFactory) NewTxBuilder() model.TxBuilder {
	return &TxBuilder{
		&proskenion.Transaction{
			Payload:    &proskenion.Transaction_Payload{},
			Signatures: []*proskenion.Signature{},
		},
		f.cryptor,
		f.executor,
		f.validator,
	}
}

type TxBuilder struct {
	*proskenion.Transaction
	cryptor   core.Cryptor
	executor  core.CommandExecutor
	validator core.CommandValidator
}

func (t *TxBuilder) CreatedTime(time int64) model.TxBuilder {
	t.Payload.CreatedTime = time
	return t
}

func (t *TxBuilder) Transfer(srcAccountId string, destAccountId string, amount int64) model.TxBuilder {
	t.Payload.Commands = append(t.Payload.Commands,
		&proskenion.Command{
			Command: &proskenion.Command_Transfer{
				Transfer: &proskenion.Transfer{
					SrcAccountId:  srcAccountId,
					DestAccountId: destAccountId,
					Amount:        amount,
				},
			},
		})
	return t
}

func (t *TxBuilder) Build() model.Transaction {
	return &Transaction{t.Transaction,
		t.cryptor, t.executor, t.validator}
}
