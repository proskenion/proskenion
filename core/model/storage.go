package model

type Storage interface {
	GetObject() map[string]Object
	GetFromKey(key string) Object
	Modelor
}
