package crypto

import (
	"crypto/sha256"
	"github.com/pkg/errors"
	. "github.com/proskenion/proskenion/core"
	. "github.com/proskenion/proskenion/core/model"
	"golang.org/x/crypto/ed25519"
	"crypto/rand"
)

type Ed25519Sha256Cryptor struct{}

func NewEd25519Sha256Cryptor() Cryptor {
	return Ed25519Sha256Cryptor{}
}

func (c Ed25519Sha256Cryptor) Hash(marshaler Marshaler) (Hash, error) {
	marshal, err := marshaler.Marshal()
	if err != nil {
		return nil, errors.Wrap(ErrMarshal, err.Error())
	}
	sha := sha256.New()
	sha.Write(marshal)
	return sha.Sum(nil), nil
}

func (c Ed25519Sha256Cryptor) Sign(hasher Hasher, privateKey PrivateKey) ([]byte, error) {
	hash, err := hasher.Hash()
	if err != nil {
		return nil, errors.Wrap(ErrHash, err.Error())
	}
	return ed25519.Sign(ed25519.PrivateKey(privateKey), hash), nil
}

func (c Ed25519Sha256Cryptor) Verify(publicKey PublicKey, hasher Hasher, signature []byte) error {
	hash, err := hasher.Hash()
	if err != nil {
		return err
	}
	if ok := ed25519.Verify(ed25519.PublicKey(publicKey), hash, signature); !ok {
		return errors.Errorf("ed25519.Verify is invalid\n"+
			"pubkey: %x\n"+
			"message: %x\n"+
			"signature: %x",
			publicKey, hash, signature)
	}
	return nil
}

func (c Ed25519Sha256Cryptor) NewKeyPairs() (PublicKey, PrivateKey) {
	a, b, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(errors.Errorf("ed25519.GenerateKey(rand.Reader) failed: %v", err))
	}
	return PublicKey(a), PrivateKey(b)
}
