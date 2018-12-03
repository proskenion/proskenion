package convertor

import (
	"github.com/golang/protobuf/proto"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
)

type Storage struct {
	cryptor core.Cryptor
	*proskenion.Storage
}

func (s *Storage) GetObject() map[string]model.Object {
	if s.Storage == nil {
		return nil
	}
	dict := make(map[string]model.Object)
	for k, v := range s.Storage.GetObject() {
		dict[k] = &Object{
			s.cryptor,
			v,
		}
	}
	return dict
}

func (s *Storage) Marshal() ([]byte, error) {
	return proto.Marshal(s.Storage)
}

func (s *Storage) Unmarshal(pb []byte) error {
	return proto.Unmarshal(pb, s.Storage)
}

func (s *Storage) Hash() model.Hash {
	return s.cryptor.Hash(s)
}
