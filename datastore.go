package go_repository_base

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/appengine/datastore"
	"reflect"
)

type DatastoreRepository struct {
	ctx          context.Context
	kind         string
	entity       interface{}
	typeName     string
	createEntity func() interface{}
	createList   func() interface{}
	*IDManager
}

// Create Repository
func NewDatastoreRepository(ctx context.Context, createEntity func() interface{}, createList func() interface{}) (*DatastoreRepository, error) {
	e := createEntity()
	rt := reflect.TypeOf(e)
	if rt.Kind() != reflect.Ptr || rt.Elem().Kind() != reflect.Struct {
		return nil, errors.New(fmt.Sprintf("Invalid entity type must be Ptr of Struct. actual: %s", rt.String()))
	}
	r := &DatastoreRepository{
		ctx:          ctx,
		entity:       e,
		typeName:     rt.String(),
		kind:         rt.Elem().String(),
		createEntity: createEntity,
		createList:   createList,
		IDManager:    &IDManager{},
	}
	if _, err := r.GetIDFieldName(e); err != nil {
		return nil, err
	}
	return r, nil
}

// Find one Entity
func (r *DatastoreRepository) Find(id int64) (interface{}, error) {
	e := r.createEntity()
	err := datastore.Get(r.ctx, r.NewKey(id), e)
	return e, err
}

// Find all entity
func (r *DatastoreRepository) FindAll() (interface{}, error) {
	return r.FindByQuery(datastore.NewQuery(r.kind))
}

// Find by query
func (r *DatastoreRepository) FindByQuery(query *datastore.Query) (interface{}, error) {
	list := reflect.ValueOf(r.createList())

	it := query.Run(r.ctx)
	for {
		e := r.createEntity()
		_, err := it.Next(e)
		if err != nil {
			break
		}
		list = reflect.Append(list, reflect.ValueOf(e))
	}

	return list.Interface(), nil
}

// Create
func (r *DatastoreRepository) Create(e interface{}) error {
	if err := r.ValidateEntity(e); err != nil {
		return err
	}
	key := datastore.NewIncompleteKey(r.ctx, r.kind, nil)
	newKey, err := datastore.Put(r.ctx, key, e)
	if err != nil {
		return err
	}
	return r.CreateWithID(e, newKey.IntID())
}

// Create with id
func (r *DatastoreRepository) CreateWithID(e interface{}, id int64) error {
	if err := r.ValidateEntity(e); err != nil {
		return err
	}
	if err := r.SetID(e, id); err != nil {
		return err
	}
	_, err := datastore.Put(r.ctx, r.NewKey(id), e)
	return err
}

// Update
func (r *DatastoreRepository) Update(e interface{}) error {
	if err := r.ValidateEntity(e); err != nil {
		return err
	}
	id, err := r.GetID(e)
	if err != nil {
		return err
	}
	_, err = datastore.Put(r.ctx, r.NewKey(id), e)
	return err
}

// Delete
func (r *DatastoreRepository) Delete(e interface{}) error {
	if err := r.ValidateEntity(e); err != nil {
		return err
	}
	id, err := r.GetID(e)
	if err != nil {
		return err
	}
	return datastore.Delete(r.ctx, r.NewKey(id))
}

// New Datastore
func (r *DatastoreRepository) NewKey(id int64) *datastore.Key {
	return datastore.NewKey(r.ctx, r.kind, "", id, nil)
}

// Validation entity type
func (r *DatastoreRepository) ValidateEntity(e interface{}) error {
	rt := reflect.TypeOf(e)
	if r.typeName != rt.String() {
		return errors.New(fmt.Sprintf("Invalid entity type. expected:%s, actual: %s", r.typeName, rt.String()))
	}
	return nil
}
