package datastore

import (
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
	"testing"
)

type myStruct struct {
	ID   int64
	Name string
}

func (s *myStruct) GetID() int64 {
	return s.ID
}

func (s *myStruct) SetID(id int64) {
	s.ID = id
}

type myStruct2 struct {
	ID   int64
	Name string
}

func (s *myStruct2) GetID() int64 {
	return s.ID
}

func (s *myStruct2) SetID(id int64) {
	s.ID = id
}

func TestDatastoreRepository_ValidateEntity(t *testing.T) {
	ctx, f, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer f()

	r, err := NewDatastoreRepository(ctx, &myStruct{})
	if err != nil {
		t.Fatal()
	}

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

	r, err := NewDatastoreRepository(ctx, &myStruct{})
	if err != nil {
		t.Fatal()
	}

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

	r, _ := NewDatastoreRepository(ctx, &myStruct{})

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

	r, _ := NewDatastoreRepository(ctx, &myStruct{})

	e := &myStruct{Name: "bob"}
	r.Create(e)

	e2 := &myStruct{ID: e.ID}
	if err := r.Find(e2); err != nil {
		t.Error(err)
	}
	if e2.Name != e.Name {
		t.Errorf("invalid name: %+v", e2)
	}
	t.Logf("%+v", e2)
}

func TestDatastoreRepository_FindAll(t *testing.T) {
	ctx, f, _ := aetest.NewContext()
	defer f()

	r, _ := NewDatastoreRepository(ctx, &myStruct{})

	r.Create(&myStruct{Name: "alice"})
	r.Create(&myStruct{Name: "bob"})

	list := []*myStruct{}
	if err := r.FindAll(&list); err != nil {
		t.Error(err)
	}
	if len(list) != 2 {
		t.Errorf("invalid length: %+v", list)
	}
	for _, e := range list {
		t.Logf("%+v", e)
	}
}

func TestDatastoreRepository_Update(t *testing.T) {
	ctx, f, _ := aetest.NewContext()
	defer f()

	r, _ := NewDatastoreRepository(ctx, &myStruct{})

	e := &myStruct{Name: "bob"}
	r.Create(e)

	newName := "bob Jr."
	e.Name = newName
	if err := r.Update(e); err != nil {
		t.Error(err)
	}

	e2 := &myStruct{ID: e.ID}
	if err := r.Find(e2); err != nil {
		t.Error(err)
	}
	if e2.Name != newName {
		t.Errorf("invalid name: %+v", e2)
	}
	t.Logf("%+v", e2)
}

func TestDatastoreRepository_Delete(t *testing.T) {
	ctx, f, _ := aetest.NewContext()
	defer f()

	r, _ := NewDatastoreRepository(ctx, &myStruct{})

	e := &myStruct{Name: "bob"}
	r.Create(e)

	if err := r.Delete(e); err != nil {
		t.Error(err)
	}

	e2 := &myStruct{ID: e.ID}
	if err := r.Find(e2); err == nil {
		t.Errorf("deleted entity found: %+v", e2)
	}
	t.Logf("%+v", e2)
}
