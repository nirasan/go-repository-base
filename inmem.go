package go_repository_base

import (
	"errors"
	"fmt"
	"reflect"
)

type InmemRepository struct {
	Data     map[interface{}]interface{}
	Id       int64
	Entity   interface{}
	TypeName string
	*IDManager
}

func NewInmemRepository(e interface{}) (*InmemRepository, error) {
	m, err := NewIDManager(e)
	if err != nil {
		return nil, err
	}
	return &InmemRepository{
		Data:      make(map[interface{}]interface{}),
		Id:        1,
		Entity:    e,
		TypeName:  reflect.TypeOf(e).String(),
		IDManager: m,
	}, nil
}

func (r *InmemRepository) Find(id interface{}) (interface{}, error) {
	e, ok := r.Data[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return e, nil
}

func (r *InmemRepository) FindAll() (interface{}, error) {
	rt := reflect.SliceOf(reflect.TypeOf(r.Entity))
	rlist := reflect.MakeSlice(rt, 0, 0)
	for _, e := range r.Data {
		rlist = reflect.Append(rlist, reflect.ValueOf(e))
	}
	return rlist.Interface(), nil
}

func (r *InmemRepository) Create(e interface{}) error {
	if err := r.ValidateEntity(e); err != nil {
		return err
	}
	r.Id = r.Id + 1
	if r.isIntID {
		r.CreateWithID(e, r.Id)
	} else {
		r.CreateWithID(e, fmt.Sprintf("%d", r.Id))
	}
	return nil
}

func (r *InmemRepository) CreateWithID(e interface{}, id interface{}) error {
	if err := r.ValidateEntity(e); err != nil {
		return err
	}
	r.Data[id] = e
	r.SetID(e, r.Id)
	return nil
}

func (r *InmemRepository) Update(e interface{}) error {
	if err := r.ValidateEntity(e); err != nil {
		return err
	}
	id, _ := r.GetID(e)
	r.Data[id] = e
	return nil
}

func (r *InmemRepository) Delete(e interface{}) error {
	if err := r.ValidateEntity(e); err != nil {
		return err
	}
	id, _ := r.GetID(e)
	delete(r.Data, id)
	return nil
}

// Validation entity type
func (r *InmemRepository) ValidateEntity(e interface{}) error {
	rt := reflect.TypeOf(e)
	if r.TypeName != rt.String() {
		return errors.New(fmt.Sprintf("Invalid entity type. expected:%s, actual: %s", r.TypeName, rt.String()))
	}
	return nil
}
