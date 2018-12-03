package model

type Storage interface {
	GetObject() map[string]Object
}
