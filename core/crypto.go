package core

import . "github.com/proskenion/proskenion/core/model"

type Marshaler interface {
	Marshal() ([]byte, error)
}

type Hasher interface {
	GetHash() (Hash, error)
}

type Cryptor interface {
	Hash(marshaler Marshaler) Hash
	Sign(hasher Hasher, privateKey PrivateKey) []byte
}

func MustGetHash(hasher Hasher) Hash {
	hash, err := hasher.GetHash()
	if err != nil {
		panic("Unexpected GetHash : " + err.Error())
	}
	return hash
}
