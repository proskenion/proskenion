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
