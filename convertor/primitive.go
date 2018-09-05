package convertor

import (
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
)

type Signature struct {
	*proskenion.Signature
}

func (s *Signature) GetPublicKey() model.PublicKey {
	return s.Signature.GetPublicKey()
}