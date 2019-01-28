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
