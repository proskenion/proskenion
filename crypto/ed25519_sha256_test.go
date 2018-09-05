package crypto_test

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockMarshaler struct {
	v string
}

type MockErrMarshaler struct {
	v string
}

func (m *MockMarshaler) Marshal() ([]byte, error) {
	return []byte(m.v), nil
}

func (m *MockErrMarshaler) Marshal() ([]byte, error) {
	return nil, errors.Errorf("Mock Error Marshal")
}

func TestEd25519Sha256Cryptor_Hash(t *testing.T) {
	cryptor := NewEd25519Sha256Cryptor()

	for _, c := range []struct {
		name      string
		marshaler core.Marshaler
		exp       string
		expErr    error
	}{
		{
			"success case 1",
			&MockMarshaler{""},
			"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			nil,
		},
		{
			"success case 2",
			&MockMarshaler{"a"},
			"ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb",
			nil,
		},
		{
			"failed err marshal",
			&MockErrMarshaler{""},
			"",
			core.ErrMarshal,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			hash, err := cryptor.Hash(c.marshaler)
			if c.expErr != nil {
				assert.EqualError(t, errors.Cause(err), c.expErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, c.exp, fmt.Sprintf("%x", hash))
			}
		})
	}
}

func (m *MockMarshaler) Hash() (model.Hash, error) {
	return NewEd25519Sha256Cryptor().Hash(m)
}

func (m *MockErrMarshaler) Hash() (model.Hash, error) {
	return nil, errors.Errorf("Mock Error Hash")
}

func TestEd25519Sha256Cryptor_SignAndVerify(t *testing.T) {
	cryptor := NewEd25519Sha256Cryptor()
	pub, pri := cryptor.NewKeyPairs()
	fpub, fpri := cryptor.NewKeyPairs()

	for _, c := range []struct {
		name         string
		hasher       core.Hasher
		publicKey    model.PublicKey
		privateKey   model.PrivateKey
		expSignErr   error
		expVerifyErr bool
	}{
		{
			"success case 1",
			&MockMarshaler{""},
			pub,
			pri,
			nil,
			false,
		},
		{
			"success case 2",
			&MockMarshaler{"a"},
			pub,
			pri,
			nil,
			false,
		},
		{
			"failed err marshal",
			&MockErrMarshaler{""},
			pub, pri,
			core.ErrHash,
			false,
		},
		{
			"failed Verify case 1",
			&MockMarshaler{"a"},
			pub, fpri,
			nil,
			true,
		},
		{
			"failed Verify case 2",
			&MockMarshaler{"a"},
			fpub, pri,
			nil,
			true,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			signature, err := cryptor.Sign(c.hasher, c.privateKey)
			if c.expSignErr != nil {
				assert.EqualError(t, errors.Cause(err), c.expSignErr.Error())
			} else {
				assert.NoError(t, err)
				if c.expVerifyErr {
					assert.Error(t, cryptor.Verify(c.publicKey, c.hasher, signature))
				} else {
					assert.NoError(t, cryptor.Verify(c.publicKey, c.hasher, signature))
				}
			}
		})
	}
}
