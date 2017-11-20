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
	m, err := NewIDManager(e)
	if err != nil {
		return nil, err
	}
	r := &DatastoreRepository{
		ctx:          ctx,
		entity:       e,
		typeName:     rt.String(),
		kind:         rt.Elem().String(),
		createEntity: createEntity,
		createList:   createList,
		IDManager:    m,
	}
	return r, nil
}

func (r *DatastoreRepository) Find(id interface{}) (interface{}, error) {
	var key *datastore.Key
	switch v := id.(type) {
	case int64:
		if !r.isIntID {
			return nil, errors.New("id must string")
		}
		key = datastore.NewKey(r.ctx, r.kind, "", v, nil)
	case string:
		if r.isIntID {
			return nil, errors.New("id must int64")
		}
		key = datastore.NewKey(r.ctx, r.kind, v, 0, nil)
	}
	e := r.createEntity()
	err := datastore.Get(r.ctx, key, e)
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
	if r.isIntID {
		return r.CreateWithID(e, newKey.IntID())
	} else {
		return r.CreateWithID(e, newKey.StringID())
	}
}

// Create with id
func (r *DatastoreRepository) CreateWithID(e interface{}, id interface{}) error {
	if err := r.ValidateEntity(e); err != nil {
		return err
	}
	if err := r.SetID(e, id); err != nil {
		return err
	}
	key, err := r.NewKey(e)
	if err != nil {
		return err
	}
	_, err = datastore.Put(r.ctx, key, e)
	return err
}

// Update
func (r *DatastoreRepository) Update(e interface{}) error {
	if err := r.ValidateEntity(e); err != nil {
		return err
	}
	key, err := r.NewKey(e)
	if err != nil {
		return err
	}
	_, err = datastore.Put(r.ctx, key, e)
	return err
}

// Delete
func (r *DatastoreRepository) Delete(e interface{}) error {
	if err := r.ValidateEntity(e); err != nil {
		return err
	}
	key, err := r.NewKey(e)
	if err != nil {
		return err
	}
	return datastore.Delete(r.ctx, key)
}

// New Datastore
func (r *DatastoreRepository) NewKey(e interface{}) (*datastore.Key, error) {
	id, err := r.GetID(e)
	if err != nil {
		return nil, err
	}
	if r.isIntID {
		return datastore.NewKey(r.ctx, r.kind, "", id.(int64), nil), nil
	} else {
		return datastore.NewKey(r.ctx, r.kind, id.(string), 0, nil), nil
	}
}

// Validation entity type
func (r *DatastoreRepository) ValidateEntity(e interface{}) error {
	rt := reflect.TypeOf(e)
	if r.typeName != rt.String() {
		return errors.New(fmt.Sprintf("Invalid entity type. expected:%s, actual: %s", r.typeName, rt.String()))
	}
	return nil
}
