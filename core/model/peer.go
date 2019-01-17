package model

type PeerWithPriKey interface {
	GetPeerId() string
	GetAddress() string
	GetPublicKey() PublicKey
	GetPrivateKey() PrivateKey
}
