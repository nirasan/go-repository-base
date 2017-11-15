package inmem

import (
	"testing"
)

type myStruct struct {
	ID   int64 `repository:"id"`
	Name string
}

func TestInmemRepository_Find(t *testing.T) {
	r := NewInmemRepository(&myStruct{})
	e := &myStruct{ID: 1, Name: "bob"}
	r.Create(e)

	if e2, err := r.Find(e.ID); err == nil {
		t.Logf("%+v", e2)
	} else {
		t.Error(err)
	}
}

func TestInmemRepository_FindAll(t *testing.T) {
	r := NewInmemRepository(&myStruct{})

	r.Create(&myStruct{Name: "alice"})
	r.Create(&myStruct{Name: "bob"})

	var list []*myStruct
	if l, err := r.FindAll(); err == nil {
		list = l.([]*myStruct)
	} else {
		t.Error(err)
	}
	if len(list) != 2 {
		t.Errorf("invalid response: %+v", list)
	}
	for _, e := range list {
		t.Logf("%+v", e)
	}
}

func TestInmemRepository_CreateWithID(t *testing.T) {
	r := NewInmemRepository(&myStruct{})

	var id int64 = 100
	e := &myStruct{Name: "bob"}
	if err := r.CreateWithID(e, id); err != nil {
		t.Error(err)
	}
	if len(r.data) != 1 {
		t.Error("failed to create with id")
	}

	if _, ok := r.data[id]; !ok {
		t.Error("not found")
	}

	t.Logf("%+v", r.data)
}

func TestInmemRepository_Create(t *testing.T) {
	r := NewInmemRepository(&myStruct{})

	e1 := &myStruct{Name: "alice"}
	e2 := &myStruct{Name: "bob"}

	if err := r.Create(e1); err != nil {
		t.Error(err)
	}
	if err := r.Create(e2); err != nil {
		t.Error(err)
	}

	if e1.ID == e2.ID {
		t.Errorf("failed to specify id: %+v, %+v", e1, e2)
	}
	t.Logf("%+v, %+v", e1, e2)
}

func TestInmemRepository_Update(t *testing.T) {
	r := NewInmemRepository(&myStruct{})

	e1 := &myStruct{Name: "alice"}
	r.Create(e1)

	newName := "alice Jr."
	e1.Name = newName
	if err := r.Update(e1); err != nil {
		t.Error(err)
	}

	e2, _ := r.Find(e1.ID)
	if e2.(*myStruct).Name != newName {
		t.Errorf("failed to update: %+v", e2)
	}
	t.Logf("%+v", e2)
}

func TestInmemRepository_Delete(t *testing.T) {
	r := NewInmemRepository(&myStruct{})

	e1 := &myStruct{Name: "alice"}
	r.Create(e1)

	if err := r.Delete(e1); err != nil {
		t.Error(err)
	}

	if e2, err := r.Find(e1.ID); err == nil {
		t.Error("failed to delete: %+v, %+v", e2, err)
	}
	t.Logf("%+v", r)
}
