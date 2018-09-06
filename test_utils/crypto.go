package test_utils

import (
	"encoding/hex"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/crypto"
)

func RandomKeyPairs() (model.PublicKey, model.PrivateKey) {
	return crypto.NewEd25519Sha256Cryptor().NewKeyPairs()
}

func MustHash(hasher core.Hasher) model.Hash {
	hash, _ := hasher.Hash()
	return hash
}

type RandomMockMarshaler struct {
	a string
}

func (m *RandomMockMarshaler) Marshal() ([]byte, error) {
	return hex.DecodeString(m.a)
}

func (m *RandomMockMarshaler) Unmarshal(pb []byte) error {
	m.a = hex.EncodeToString(pb)
	return nil
}

func RandomMarshaler() *RandomMockMarshaler {
	return &RandomMockMarshaler{hex.EncodeToString(RandomByte())}
}
