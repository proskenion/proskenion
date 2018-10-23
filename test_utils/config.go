package test_utils

import (
	"encoding/hex"
	"github.com/proskenion/proskenion/config"
)

func NewTestConfig() *config.Config {
	config := config.NewConfig("../config/config.yaml")

	pub, pri := RandomCryptor().NewKeyPairs()
	config.Peer.PublicKey = hex.EncodeToString(pub)
	config.Peer.PrivateKey = hex.EncodeToString(pri)

	config.DB.Path = "../database"
	return config
}
