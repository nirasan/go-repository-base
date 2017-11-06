package datastore

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/appengine/datastore"
	"reflect"
)

type Entity interface {
	GetID() int64
}

var EntityType = reflect.TypeOf((*Entity)(nil)).Elem()

type DatastoreRepository struct {
	ctx      context.Context
	kind     string
	entity   Entity
	typeName string
}

// Create Repository
func NewDatastoreRepository(ctx context.Context, e Entity) (*DatastoreRepository, error) {
	rt := reflect.TypeOf(e)
	if rt.Kind() != reflect.Ptr || rt.Elem().Kind() != reflect.Struct {
		return nil, errors.New(fmt.Sprintf("Invalid entity type must be Ptr of Struct. actual: %s", rt.String()))
	}
	r := &DatastoreRepository{ctx: ctx, entity: e}
	r.typeName = rt.String()
	r.kind = rt.Elem().String()
	return r, nil
}

// Find one entity
func (r *DatastoreRepository) Find(e Entity) error {
	if err := r.ValidateEntity(e); err != nil {
		return err
	}
	err := datastore.Get(r.ctx, r.NewKey(e.GetID()), e)
	return err
}

// Find all entity
func (r *DatastoreRepository) FindAll(list interface{}) error {
	rt := reflect.TypeOf(list)
	if rt.Kind() != reflect.Ptr || rt.Elem().Kind() != reflect.Slice || !rt.Elem().Elem().Implements(EntityType) {
		return errors.New(fmt.Sprintf("Invalid entity type. expected:[]%s, actual: %s", r.typeName, rt.String()))
	}
	_, err := datastore.NewQuery(r.kind).GetAll(r.ctx, list)
	return err
}

// Create
func (r *DatastoreRepository) Create(e Entity) error {
	if err := r.ValidateEntity(e); err != nil {
		return err
	}
	_, err := datastore.Put(r.ctx, r.NewKey(e.GetID()), e)
	return err
}

// Update
func (r *DatastoreRepository) Update(e Entity) error {
	if err := r.ValidateEntity(e); err != nil {
		return err
	}
	_, err := datastore.Put(r.ctx, r.NewKey(e.GetID()), e)
	return err
}

// Delete
func (r *DatastoreRepository) Delete(e Entity) error {
	if err := r.ValidateEntity(e); err != nil {
		return err
	}
	return datastore.Delete(r.ctx, r.NewKey(e.GetID()))
}

// New Datastore
func (r *DatastoreRepository) NewKey(id int64) *datastore.Key {
	return datastore.NewKey(r.ctx, r.kind, "", id, nil)
}

// Validation entity type
func (r *DatastoreRepository) ValidateEntity(e Entity) error {
	rt := reflect.TypeOf(e)
	if r.typeName != rt.String() {
		return errors.New(fmt.Sprintf("Invalid entity type. expected:%s, actual: %s", r.typeName, rt.String()))
	}
	return nil
}
