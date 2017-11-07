package inmem

import (
	. "github.com/nirasan/go-repository-base"
	"errors"
	"reflect"
	"fmt"
)

type InmemRepository struct {
	data map[int64]Entity
	id int64
	entity Entity
	typeName string
}

func NewInmemRepository(e Entity) *InmemRepository {
	return &InmemRepository{
		data: make(map[int64]Entity),
		id: 1,
		entity: e,
		typeName: reflect.TypeOf(e).String(),
	}
}

func (r *InmemRepository) Find(e Entity) error {
	if err := r.ValidateEntity(e); err != nil {
		return err
	}
	ee, ok := r.data[e.GetID()]
	if !ok {
		return errors.New("not found")
	}
	rv := reflect.ValueOf(e)
	rt := rv.Type()
	if rt.Kind() != reflect.Ptr || rt.Elem().Kind() != reflect.Struct {
		return errors.New("invalid type")
	}
	rv.Elem().Set(reflect.ValueOf(ee).Elem())
	return nil
}

func (r *InmemRepository) FindAll(list interface{}) error {
	rt := reflect.TypeOf(list)
	if rt.Kind() != reflect.Ptr || rt.Elem().Kind() != reflect.Slice || !rt.Elem().Elem().Implements(EntityType) {
		return errors.New(fmt.Sprintf("Invalid entity type. expected:[]%s, actual: %s", reflect.TypeOf(r.entity), rt.String()))
	}
	rlist := reflect.ValueOf(list).Elem()
	for _, e := range r.data {
		rlist.Set(reflect.Append(rlist, reflect.ValueOf(e)))
	}
	return nil
}

func (r *InmemRepository) Create(e Entity) error {
	if err := r.ValidateEntity(e); err != nil {
		return err
	}
	r.id = r.id + 1
	r.CreateWithID(e, r.id)
	return nil
}

func (r *InmemRepository) CreateWithID(e Entity, id int64) error {
	if err := r.ValidateEntity(e); err != nil {
		return err
	}
	r.data[id] = e
	e.SetID(r.id)
	return nil
}

func (r *InmemRepository) Update(e Entity) error {
	if err := r.ValidateEntity(e); err != nil {
		return err
	}
	r.data[e.GetID()] = e
	return nil
}

func (r *InmemRepository) Delete(e Entity) error {
	if err := r.ValidateEntity(e); err != nil {
		return err
	}
	delete(r.data, e.GetID())
	return nil
}

// Validation entity type
func (r *InmemRepository) ValidateEntity(e Entity) error {
	rt := reflect.TypeOf(e)
	if r.typeName != rt.String() {
		return errors.New(fmt.Sprintf("Invalid entity type. expected:%s, actual: %s", r.typeName, rt.String()))
	}
	return nil
}
