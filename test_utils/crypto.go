package test_utils

import (
	"encoding/hex"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type Randomer interface {
	model.Hasher
	model.Marshaler
	model.Unmarshaler
}

func RandomCryptor() core.Cryptor {
	return crypto.NewEd25519Sha256Cryptor()
}

func RandomVerify(t *testing.T, pubkey model.PublicKey, hasher model.Hasher, sig []byte) {
	assert.NoError(t, RandomCryptor().Verify(pubkey, hasher, sig))
}

func RandomKeyPairs() (model.PublicKey, model.PrivateKey) {
	return crypto.NewEd25519Sha256Cryptor().NewKeyPairs()
}

func RandomPublicKey() model.PublicKey {
	p, _ := crypto.NewEd25519Sha256Cryptor().NewKeyPairs()
	return p
}

func RandomPurivateKey() model.PrivateKey {
	_, p := crypto.NewEd25519Sha256Cryptor().NewKeyPairs()
	return p
}

func MustHash(hasher model.Hasher) model.Hash {
	return hasher.Hash()
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

func (m *RandomMockMarshaler) Hash() model.Hash {
	ret, _ := m.Marshal()
	return ret
}

func RandomMarshaler() *RandomMockMarshaler {
	return &RandomMockMarshaler{hex.EncodeToString(RandomByte())}
}

func RandomMarshalerFromStr(s string) *RandomMockMarshaler {
	return &RandomMockMarshaler{s}
}

func DecodeMustString(t *testing.T, s string) []byte {
	h, err := hex.DecodeString(s)
	require.NoError(t, err)
	return h
}
