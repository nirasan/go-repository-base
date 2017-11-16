# go-repository-base

* Generalize the CRUD processing of the repository with Google Datastore as an infrastructure to lower the burden of implementation.

# Usage

* Implement User entity's repository.

## Before

```go
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
```

## After

```go
type User struct {
	ID   int64 `repository:"id"`
	Name string
}

type UserDatastoreRepository struct {
	*go_repository_base.DatastoreRepository
}

func NewUserDatastoreRepository(ctx context.Context) (*UserDatastoreRepository, error) {
	r, err := go_repository_base.NewDatastoreRepository(ctx, func() interface{} { return &User{} }, func() interface{} { return []*User{} })
	if err != nil {
		return nil, err
	}
	return &UserDatastoreRepository{r}, nil
}
```

# Benchmark

* This package uses reflect, but "google.golang.org/appengine/datastore" does not have much performance change since it originally uses reflect.

```
BenchmarkBefore-4        1000          26454248 ns/op
BenchmarkAfter-4         1000          27537754 ns/op
```

