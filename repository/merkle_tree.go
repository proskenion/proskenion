package repository

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

// 累積Hash で Transaction 列の検証には十分
type AccumulateHash struct {
	cryptor core.Cryptor
	hashes  []model.Hash
}

func NewAccumulateHash(cryptor core.Cryptor) core.MerkleTree {
	return &AccumulateHash{cryptor: cryptor, hashes: make([]model.Hash, 0)}
}

type DefaultMarshaler struct {
	marshal []byte
}

func (m *DefaultMarshaler) Marshal() ([]byte, error) {
	return m.marshal, nil
}

func newDefaultMarshaler(a []byte, b []byte) *DefaultMarshaler {
	return &DefaultMarshaler{append(a, b...)}
}

func (t *AccumulateHash) Push(hasher model.Hasher) error {
	rh := t.cryptor.Hash(newDefaultMarshaler(t.Top(), hasher.Hash()))
	t.hashes = append(t.hashes, rh)
	return nil
}

func (t *AccumulateHash) Top() model.Hash {
	if len(t.hashes) == 0 {
		return nil
	}
	return t.hashes[len(t.hashes)-1]
}
