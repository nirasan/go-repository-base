package go_repository_base

import (
	"errors"
	"fmt"
	"reflect"
)

type InmemRepository struct {
	data     map[interface{}]interface{}
	id       int64
	entity   interface{}
	typeName string
	*IDManager
}

func NewInmemRepository(e interface{}) (*InmemRepository, error) {
	m, err := NewIDManager(e)
	if err != nil {
		return nil, err
	}
	return &InmemRepository{
		data:      make(map[interface{}]interface{}),
		id:        1,
		entity:    e,
		typeName:  reflect.TypeOf(e).String(),
		IDManager: m,
	}, nil
}

func (r *InmemRepository) Find(id interface{}) (interface{}, error) {
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
	if r.isIntID {
		r.CreateWithID(e, r.id)
	} else {
		r.CreateWithID(e, fmt.Sprintf("%d", r.id))
	}
	return nil
}

func (r *InmemRepository) CreateWithID(e interface{}, id interface{}) error {
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
