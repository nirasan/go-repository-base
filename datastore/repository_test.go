package datastore

import (
	"google.golang.org/appengine/aetest"
	"testing"
)

type MyStruct struct {
	ID   int64
	Name string
}

func (s *MyStruct) GetID() int64 {
	return s.ID
}

func TestCommonRepository(t *testing.T) {
	ctx, f, err := aetest.NewContext()
	defer f()

	repo, err := NewDatastoreRepository(ctx, &MyStruct{})
	t.Logf("%+v, %+v\n", repo, err)

	e1 := &MyStruct{Name: "a", ID: 1}

	if err := repo.Create(e1); err != nil {
		t.Error(err)
	}

	e2 := &MyStruct{ID: 1}
	if err := repo.Find(e2); err != nil {
		t.Error(err)
	}
	if e2.Name != e1.Name {
		t.Errorf("Failed to Find: %+v, %+v", e2, e1)
	}

	list := []*MyStruct{}
	if err := repo.FindAll(&list); err != nil {
		t.Error(err)
	}
	if len(list) != 1 || list[0].ID != e1.ID {
		for _, e := range list {
			t.Logf("%#v", e)
			e.GetID()
		}
	}

	e2.Name = "b"
	if err = repo.Update(e2); err != nil {
		t.Error(err)
	}
	e3 := &MyStruct{ID: 1}
	if err := repo.Find(e3); err != nil || e3.Name != e2.Name {
		t.Errorf("Failed to Update: %+v, %+v, %+v", err, e3, e2)
	}

	if err = repo.Delete(e2); err != nil {
		t.Error(err)
	}
	e4 := &MyStruct{ID: 1}
	if err := repo.Find(e4); err == nil {
		t.Errorf("Failed to Delete: %+v, %+v", err, e4)
	}
}
