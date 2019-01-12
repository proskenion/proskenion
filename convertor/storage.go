package convertor

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
	"github.com/satellitex/protobuf/proto"
)

type Storage struct {
	cryptor core.Cryptor
	*proskenion.Storage
}

func ProslObjectMapsFromObjectMaps(objects map[string]model.Object) map[string]*proskenion.Object {
	ret := make(map[string]*proskenion.Object)
	for key, value := range objects {
		ret[key] = value.(*Object).Object
	}
	return ret
}

func ProslObjectListFromObjectList(objects []model.Object) []*proskenion.Object {
	ret := make([]*proskenion.Object, 0)
	for _, value := range objects {
		ret = append(ret, value.(*Object).Object)
	}
	return ret
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

func (s *Storage) GetFromKey(key string) model.Object {
	return s.GetObject()[key]
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
