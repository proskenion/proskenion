package model

type Storage interface {
	GetObject() map[string]Object
	GetId() string
	GetFromKey(key string) Object
	Modelor
}
