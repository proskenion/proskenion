package repository

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

// World State の管理
type MerkleParticleTree struct {
	dba  core.KeyValueStore
	root core.MerkleParticleTreeIterator
}

func (t *MerkleParticleTree) Find(key core.Marshaler, value core.Unmarshaler) error {
	return t.Iterator().Find(key)
}
func (t *MerkleParticleTree) Upsert(key core.Marshaler, value core.Unmarshaler) error {
	return t.Iterator().Upser(key, value)
}
func (t *MerkleParticleTree) Hash() (model.Hash, error) {
	return nil, nil
}
func (t *MerkleParticleTree) Marshal() ([]byte, error) {
	return nil, nil
}
func (t *MerkleParticleTree) Unmarshal(b []byte) error {
	return nil
}

type MerkleParticleIterator interface {
	Data(unmarshaler core.Unmarshaler) error
	History() MerkleParticleNodeIterator
	Next() MerkleParticleIterator
	Prev() MerkleParticleIterator
	First() bool
	Last() bool
	Hash() (model.Hash, error)
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

type MerkleParticleNodeIterator struct {
	dba core.KeyValueStore
}

func (t *MerkleParticleNodeIterator) Data(unmarshaler core.Unmarshaler) error {
	return nil
}
func (t *MerkleParticleNodeIterator) Next() MerkleParticleNodeIterator {
	return nil
}
func (t *MerkleParticleNodeIterator) Prev() MerkleParticleNodeIterator {
	return nil
}
func (t *MerkleParticleNodeIterator) First() bool {
	return false
}
func (t *MerkleParticleNodeIterator) Last() bool {
	return false
}
func (t *MerkleParticleNodeIterator) Marshal() ([]byte, error) {
	return nil, nil
}
func (t *MerkleParticleNodeIterator) Unmarshal([]byte) error {
	return nil
}
