package model

import (
	"github.com/pkg/errors"
)

var ErrInvalidSignature = errors.Errorf("Failed Invalid Signature")

type PublicKey []byte
type PrivateKey []byte
type Hash []byte

type Signature interface {
	GetPublicKey() PublicKey
	GetSignature() []byte
}

func PublicKeysFromBytesSlice(keys [][]byte) []PublicKey {
	ret := make([]PublicKey, len(keys))
	for i, key := range keys {
		ret[i] = key
	}
	return ret
}

func BytesListFromPublicKeys(keys []PublicKey) [][]byte {
	ret := make([][]byte, len(keys))
	for i, key := range keys {
		ret[i] = key
	}
	return ret

}
