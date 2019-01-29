package model

import (
	"bytes"
	"encoding/gob"
)

type Marshaler interface {
	Marshal() ([]byte, error)
}

type Unmarshaler interface {
	Unmarshal([]byte) error
}

type Hasher interface {
	Hash() Hash
}

type Modelor interface {
	Marshaler
	Unmarshaler
	Hasher
}

type UnmarshalerFactory interface {
	CreateUnmarshaler() Unmarshaler
}

func GobMarshal(e interface{}) ([]byte, error) {
	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.
	err := enc.Encode(e)
	if err != nil {
		return nil, err
	}
	return network.Bytes(), nil
}

func GobUnmarshal(b []byte, e interface{}) error {
	network := bytes.NewBuffer(b)
	dec := gob.NewDecoder(network) // Will read from network.
	err := dec.Decode(e)
	if err != nil {
		return err
	}
	return nil
}

type AccountUnmarshalerFactory struct {
	fc ModelFactory
}

func (f *AccountUnmarshalerFactory) CreateUnmarshaler() Unmarshaler {
	return f.fc.NewEmptyAccount()
}

func NewAccountUnmarshalerFactory(fc ModelFactory) UnmarshalerFactory {
	return &AccountUnmarshalerFactory{fc}
}

type PeerUnmarshalerFactory struct {
	fc ModelFactory
}

func (f *PeerUnmarshalerFactory) CreateUnmarshaler() Unmarshaler {
	return f.fc.NewEmptyPeer()
}

func NewPeerUnmarshalerFactory(fc ModelFactory) UnmarshalerFactory {
	return &PeerUnmarshalerFactory{fc}
}

type StorageUnmarshalerFactory struct {
	fc ModelFactory
}

func (f *StorageUnmarshalerFactory) CreateUnmarshaler() Unmarshaler {
	return f.fc.NewEmptyStorage()
}

func NewStorageUnmarshalerFactory(fc ModelFactory) UnmarshalerFactory {
	return &StorageUnmarshalerFactory{fc}
}