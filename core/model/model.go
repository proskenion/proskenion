package model

type Marshaler interface {
	Marshal() ([]byte, error)
}

type Unmarshaler interface {
	Unmarshal([]byte) error
}

type Hasher interface {
	Hash() (Hash, error)
}

type Modelor interface {
	Marshaler
	Unmarshaler
	Hasher
}