package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"github.com/pkg/errors"
	. "github.com/proskenion/proskenion/core"
	. "github.com/proskenion/proskenion/core/model"
	"golang.org/x/crypto/ed25519"
	"strconv"
)

type Ed25519Sha256Cryptor struct{}

func NewEd25519Sha256Cryptor() Cryptor {
	return Ed25519Sha256Cryptor{}
}

func (c Ed25519Sha256Cryptor) Hash(marshaler Marshaler) Hash {
	if marshaler == nil {
		return nil
	}
	marshal, err := marshaler.Marshal()
	if err != nil {
		return nil
	}
	sha := sha256.New()
	sha.Write(marshal)
	return sha.Sum(nil)
}

func (c Ed25519Sha256Cryptor) ConcatHash(hashes ...Hash) Hash {
	if len(hashes) == 0 {
		return Hash(nil)
	}
	if len(hashes) == 1 {
		return hashes[0]
	}
	thash := hashes[0]
	for _, hash := range hashes[1:] {
		thash = append(thash, hash...)
	}
	sha := sha256.New()
	sha.Write(thash)
	return sha.Sum(nil)
}

func (c Ed25519Sha256Cryptor) Sign(hasher Hasher, privateKey PrivateKey) ([]byte, error) {
	hash := hasher.Hash()
	if hash == nil {
		return nil, ErrHash
	}
	return ed25519.Sign(ed25519.PrivateKey(privateKey), hash), nil
}

func (c Ed25519Sha256Cryptor) Verify(publicKey PublicKey, hasher Hasher, signature []byte) error {
	hash := hasher.Hash()
	if l := len(publicKey); l != ed25519.PublicKeySize {
		return errors.Errorf("ed25519: bad public key length: " + strconv.Itoa(l))
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
