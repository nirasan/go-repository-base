package test

import (
	"context"
	"google.golang.org/appengine/datastore"
)

type User struct {
	ID   int64 `repository:"id"`
	Name string
}

type UserRepository struct {
	ctx  context.Context
	kind string
}

func NewUserRepository(ctx context.Context) *UserRepository {
	return &UserRepository{ctx: ctx, kind: "User"}
}

func (r *UserRepository) Find(id int64) (*User, error) {
	e := &User{}
	if err := datastore.Get(r.ctx, r.NewKey(id), e); err != nil {
		return nil, err
	}
	return e, nil
}

func (r *UserRepository) FindAll() ([]*User, error) {
	list := []*User{}

	it := datastore.NewQuery(r.kind).Run(r.ctx)
	for {
		e := &User{}
		_, err := it.Next(e)
		if err != nil {
			break
		}
		list = append(list, e)
	}

	return list, nil
}

func (r *UserRepository) Create(e *User) error {
	key := datastore.NewIncompleteKey(r.ctx, r.kind, nil)
	newKey, err := datastore.Put(r.ctx, key, e)
	if err != nil {
		return err
	}
	e.ID = newKey.IntID()
	_, err = datastore.Put(r.ctx, newKey, e)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Delete(e *User) error {
	err := datastore.Delete(r.ctx, r.NewKey(e.ID))
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) NewKey(id int64) *datastore.Key {
	return datastore.NewKey(r.ctx, r.kind, "", id, nil)
}
