package model

type Peer interface {
	GetAddress() string
	GetPublicKey() PublicKey
}

type PeerWithPriKey interface {
	GetAddress() string
	GetPublicKey() PublicKey
	GetPrivateKey() PrivateKey
}