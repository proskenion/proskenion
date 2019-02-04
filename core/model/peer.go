package model

import "fmt"

type PeerWithPriKey interface {
	Peer
	GetPrivateKey() PrivateKey
}

func MakeAddressFromHostAndPort(host string, port string) string {
	return fmt.Sprintf("%s:%s", host, port)
}
