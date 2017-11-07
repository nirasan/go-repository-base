package go_repository_base

import "reflect"

type Entity interface {
	GetID() int64
	SetID(int64)
}

var EntityType = reflect.TypeOf((*Entity)(nil)).Elem()

type Repository interface {
	Find(e Entity) error
	FindAll(list interface{}) error
	Create(e Entity) error
	CreateWithID(e Entity, id int64) error
	Update(e Entity) error
	Delete(e Entity) error
}
