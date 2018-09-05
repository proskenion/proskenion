package model

type PeerWithPriKey interface {
	GetAddress() string
	GetPublicKey() PublicKey
	GetPrivateKey() PrivateKey
}
