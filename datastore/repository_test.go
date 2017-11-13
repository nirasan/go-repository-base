package datastore

import (
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
	"testing"
	"context"
)

type myStruct struct {
	ID   int64 `repository:"id"`
	Name string
}

func createMyStructRepository(ctx context.Context) *DatastoreRepository {
	r, _ := NewDatastoreRepository(ctx, func() interface{} { return &myStruct{} }, func() interface{} { return []*myStruct{} })
	return r
}

type myStruct2 struct {
	ID   int64
	Name string
}

func TestDatastoreRepository_SetID(t *testing.T) {
	ctx, f, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer f()
	r := createMyStructRepository(ctx)
	e := &myStruct{}
	if err := r.SetID(e, 100); err != nil {
		t.Error(err)
	}
	if e.ID != 100 {
		t.Error("invalid id: ", e.ID)
	}
}

func TestDatastoreRepository_GetID(t *testing.T) {
	ctx, f, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer f()
	r := createMyStructRepository(ctx)
	e := &myStruct{ID: 999}
	id, err := r.GetID(e)
	if err != nil {
		t.Error(err)
	}
	if id != 999 {
		t.Error("invalid id: ", id)
	}
}

func TestDatastoreRepository_ValidateEntity(t *testing.T) {
	ctx, f, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer f()

	r := createMyStructRepository(ctx)

	e1 := &myStruct2{Name: "hello"}
	err = r.ValidateEntity(e1)
	if err == nil {
		t.Error("entity must myStruct")
	}
	t.Log(err)
}

func TestDatastoreRepository_CreateWithID(t *testing.T) {
	ctx, f, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer f()

	r := createMyStructRepository(ctx)

	e := &myStruct{Name: "bob"}
	var id int64 = 100
	if err := r.CreateWithID(e, id); err != nil {
		t.Error(err)
	}

	e2 := new(myStruct)
	if err := datastore.Get(ctx, datastore.NewKey(ctx, r.kind, "", id, nil), e2); err != nil {
		t.Error(err)
	}
	if e2.ID != id {
		t.Errorf("invalid id: %+v", e2)
	}

	t.Logf("%+v", e2)
}

func TestDatastoreRepository_Create(t *testing.T) {
	ctx, f, _ := aetest.NewContext()
	defer f()

	r := createMyStructRepository(ctx)

	e := &myStruct{Name: "bob"}
	if err := r.Create(e); err != nil {
		t.Error(err)
	}

	e2 := new(myStruct)
	if err := datastore.Get(ctx, datastore.NewKey(ctx, r.kind, "", e.ID, nil), e2); err != nil {
		t.Error(err)
	}
	if e2.Name != e.Name {
		t.Errorf("invalid name: %+v", e2)
	}
	t.Logf("%+v", e2)
}

func TestDatastoreRepository_Find(t *testing.T) {
	ctx, f, _ := aetest.NewContext()
	defer f()

	r := createMyStructRepository(ctx)

	e := &myStruct{Name: "bob"}
	r.Create(e)

	ie, err := r.Find(e.ID)
	if err != nil {
		t.Error(err)
	}
	e2 := ie.(*myStruct)
	if e2.Name != e.Name {
		t.Errorf("invalid name: %+v", e2)
	}
	t.Logf("%+v", e2)
}

func TestDatastoreRepository_FindAll(t *testing.T) {
	ctx, f, _ := aetest.NewContext()
	defer f()

	r := createMyStructRepository(ctx)

	r.Create(&myStruct{Name: "alice"})
	r.Create(&myStruct{Name: "bob"})

	ret, err := r.FindAll()
	if err != nil {
		t.Error(err)
	}
	list := ret.([]*myStruct)
	if len(list) != 2 {
		t.Errorf("invalid length: %+v", list)
	}
	for _, e := range list {
		t.Logf("%+v", e)
	}
}

func TestDatastoreRepository_FindByQuery(t *testing.T) {
	ctx, f, _ := aetest.NewContext()
	defer f()

	r := createMyStructRepository(ctx)

	ee1 := &myStruct{Name: "alice"}
	ee2 := &myStruct{Name: "bob"}
	ee3 := &myStruct{Name: "carol"}
	r.Create(ee1)
	r.Create(ee2)
	r.Create(ee3)

	ret, err := r.FindByQuery(datastore.NewQuery(r.kind))
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", ret)
	list := ret.([]*myStruct)
	t.Logf("%+v", list)
	if len(list) != 3 {
		t.Errorf("invalid length: %+v", list)
	}
	for _, e := range list {
		t.Logf("%+v", e)
	}
}

func TestDatastoreRepository_Update(t *testing.T) {
	ctx, f, _ := aetest.NewContext()
	defer f()

	r := createMyStructRepository(ctx)

	e := &myStruct{Name: "bob"}
	r.Create(e)

	newName := "bob Jr."
	e.Name = newName
	if err := r.Update(e); err != nil {
		t.Error(err)
	}

	ie, err := r.Find(e.ID)
	if err != nil {
		t.Error(err)
	}
	e2 := ie.(*myStruct)
	if e2.Name != newName {
		t.Errorf("invalid name: %+v", e2)
	}
	t.Logf("%+v", e2)
}

func TestDatastoreRepository_Delete(t *testing.T) {
	ctx, f, _ := aetest.NewContext()
	defer f()

	r := createMyStructRepository(ctx)

	e := &myStruct{Name: "bob"}
	r.Create(e)

	if err := r.Delete(e); err != nil {
		t.Error(err)
	}

	ie, err := r.Find(e.ID)
	if err == nil {
		t.Errorf("deleted entity found: %+v", ie)
	}
	e2 := ie.(*myStruct)
	t.Logf("%+v", e2)
}
