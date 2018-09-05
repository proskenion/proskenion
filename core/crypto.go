package core

import (
	"github.com/pkg/errors"
	. "github.com/proskenion/proskenion/core/model"
)

var (
	ErrMarshal       = errors.Errorf("Failed Marshal")
	ErrHash       = errors.Errorf("Failed Hash")
	ErrCryptorHash   = errors.Errorf("Failed Cryptor Hash")
	ErrCryptorSign   = errors.Errorf("Failed Cryptor Sign")
	ErrCryptorVerify = errors.Errorf("Failed Cryptor Verify")
)

type Marshaler interface {
	Marshal() ([]byte, error)
}

type Hasher interface {
	Hash() (Hash, error)
}

type Cryptor interface {
	Hash(marshaler Marshaler) (Hash, error)
	Sign(hasher Hasher, privateKey PrivateKey) ([]byte, error)
	Verify(publicKey PublicKey, hasher Hasher, signature []byte) error
	NewKeyPairs() (PublicKey, PrivateKey)
}

func MustHash(hasher Hasher) Hash {
	hash, err := hasher.Hash()
	if err != nil {
		panic("Unexpected Hash : " + err.Error())
	}
	return hash
}
