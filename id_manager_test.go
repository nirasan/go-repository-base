package go_repository_base

import "testing"

func TestIDManager_GetIDFieldName(t *testing.T) {
	m := &IDManager{}

	if _, err := m.GetIDFieldName(1); err == nil {
		t.Error("invalid type")
	}

	if _, err := m.GetIDFieldName(struct{ ID int64 }{}); err == nil {
		t.Error("tag not found")
	}

	if _, err := m.GetIDFieldName(struct {
		ID int64 `repository:"pk"`
	}{}); err == nil {
		t.Error("invalid tag string")
	}

	if _, err := m.GetIDFieldName(struct {
		ID string `repository:"id"`
	}{}); err == nil {
		t.Error("invalid id type")
	}

	if _, err := m.GetIDFieldName(struct {
		ID int64 `repository:"id"`
	}{}); err == nil {
		t.Error("entity must pointer of struct")
	}

	if id, err := m.GetIDFieldName(&struct {
		ID int64 `repository:"id"`
	}{}); err != nil || id != "ID" {
		t.Error("invalid result")
	}
}

func TestIDManager_SetID(t *testing.T) {
	m := &IDManager{}
	s := &struct {
		ID int64 `repository:"id"`
	}{}

	err := m.SetID(s, 1000)
	if err != nil {
		t.Error(err)
	}

	if s.ID != 1000 {
		t.Errorf("invalid result: %+v", s)
	}
}

func TestIDManager_GetID(t *testing.T) {
	m := &IDManager{}
	s := &struct {
		ID int64 `repository:"id"`
	}{ID: 100}

	id, err := m.GetID(s)
	if err != nil {
		t.Error(err)
	}

	if id != 100 {
		t.Errorf("invalid result: %+v", s)
	}
}
