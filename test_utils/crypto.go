package test_utils

import (
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/crypto"
)

func RandomKeyPairs() (model.PublicKey, model.PrivateKey) {
	return crypto.NewEd25519Sha256Cryptor().NewKeyPairs()
}
