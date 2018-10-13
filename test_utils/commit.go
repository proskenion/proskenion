package test_utils

import "github.com/proskenion/proskenion/commit"

func RandomCommitProperty() *commit.CommitProperty {
	validPub, validPri := RandomCryptor().NewKeyPairs()
	return &commit.CommitProperty{
		NumTxInBlock: 100,
		PublicKey:    validPub,
		PrivateKey:   validPri,
	}
}

func RandomNow() int64 {
	return commit.Now()
}