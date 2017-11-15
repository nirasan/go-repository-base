package inmem

import (
	"errors"
	"fmt"
	"reflect"
)

const tagKey = "repository"

type InmemRepository struct {
	data     map[int64]interface{}
	id       int64
	entity   interface{}
	typeName string
}

func NewInmemRepository(e interface{}) *InmemRepository {
	return &InmemRepository{
		data:     make(map[int64]interface{}),
		id:       1,
		entity:   e,
		typeName: reflect.TypeOf(e).String(),
	}
}

func (r *InmemRepository) GetIDFieldName(e interface{}) (string, error) {
	if err := r.ValidateEntity(e); err != nil {
		return "", err
	}
	rv := reflect.Indirect(reflect.ValueOf(e))
	rt := rv.Type()
	if rt.Kind() != reflect.Struct {
		return "", errors.New(fmt.Sprintf("Entity must struct or pointer of struct: %+v", e))
	}
	for i := rt.NumField() - 1; i >= 0; i = i - 1 {
		f := rt.Field(i)
		if f.Type.Kind() == reflect.Int64 && f.Tag.Get(tagKey) == "id" {
			return f.Name, nil
		}
	}
	return "", errors.New("not found")
}

func (r *InmemRepository) SetID(e interface{}, id int64) error {
	name, err := r.GetIDFieldName(e)
	if err != nil {
		return err
	}
	rv := reflect.Indirect(reflect.ValueOf(e))
	rv.FieldByName(name).SetInt(id)
	return nil
}

func (r *InmemRepository) GetID(e interface{}) (int64, error) {
	name, err := r.GetIDFieldName(e)
	if err != nil {
		return 0, err
	}
	rv := reflect.Indirect(reflect.ValueOf(e))
	return rv.FieldByName(name).Int(), nil
}

func (r *InmemRepository) Find(id int64) (interface{}, error) {
	e, ok := r.data[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return e, nil
}

func (r *InmemRepository) FindAll() (interface{}, error) {
	rt := reflect.SliceOf(reflect.TypeOf(r.entity))
	rlist := reflect.MakeSlice(rt, 0, 0)
	for _, e := range r.data {
		rlist = reflect.Append(rlist, reflect.ValueOf(e))
	}
	return rlist.Interface(), nil
}

func (r *InmemRepository) Create(e interface{}) error {
	if err := r.ValidateEntity(e); err != nil {
		return err
	}
	r.id = r.id + 1
	r.CreateWithID(e, r.id)
	return nil
}

func (r *InmemRepository) CreateWithID(e interface{}, id int64) error {
	if err := r.ValidateEntity(e); err != nil {
		return err
	}
	r.data[id] = e
	r.SetID(e, r.id)
	return nil
}

func (r *InmemRepository) Update(e interface{}) error {
	if err := r.ValidateEntity(e); err != nil {
		return err
	}
	id, _ := r.GetID(e)
	r.data[id] = e
	return nil
}

func (r *InmemRepository) Delete(e interface{}) error {
	if err := r.ValidateEntity(e); err != nil {
		return err
	}
	id, _ := r.GetID(e)
	delete(r.data, id)
	return nil
}

// Validation entity type
func (r *InmemRepository) ValidateEntity(e interface{}) error {
	rt := reflect.TypeOf(e)
	if r.typeName != rt.String() {
		return errors.New(fmt.Sprintf("Invalid entity type. expected:%s, actual: %s", r.typeName, rt.String()))
	}
	return nil
}
