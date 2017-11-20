package go_repository_base

import "testing"

func TestNewIDManager(t *testing.T) {
	if _, err := NewIDManager(1); err == nil {
		t.Error("invalid type")
	}

	if _, err := NewIDManager(struct{ ID int64 }{}); err == nil {
		t.Error("tag not found")
	}

	if _, err := NewIDManager(struct {
		ID int64 `repository:"pk"`
	}{}); err == nil {
		t.Error("invalid tag string")
	}

	if _, err := NewIDManager(struct {
		ID string `repository:"id"`
	}{}); err == nil {
		t.Error("invalid id type")
	}

	if _, err := NewIDManager(struct {
		ID int64 `repository:"id"`
	}{}); err == nil {
		t.Error("entity must pointer of struct")
	}

	s1 := struct {
		ID int64 `repository:"id"`
	}{}
	if m, err := NewIDManager(&s1); err != nil || m.idFieldName != "ID" || !m.isIntID {
		t.Error("invalid result")
	}

	s2 := struct {
		ID string `repository:"id"`
	}{}
	if m, err := NewIDManager(&s2); err != nil || m.idFieldName != "ID" || m.isIntID {
		t.Error("invalid result")
	}
}

func TestIDManager_SetID(t *testing.T) {
	s := &struct {
		ID int64 `repository:"id"`
	}{}
	m, err := NewIDManager(s)
	if err != nil {
		t.Error(err)
	}

	var id int64 = 1000
	err = m.SetID(s, id)
	if err != nil {
		t.Error(err)
	}

	if s.ID != id {
		t.Errorf("invalid result: %+v", s)
	}

	strID := "2000"
	if err = m.SetID(s, strID); err == nil {
		t.Errorf("invalid result: %+v", s)
	}
}

func TestIDManager_GetID(t *testing.T) {
	var id int64 = 100
	s := &struct {
		ID int64 `repository:"id"`
	}{ID: id}
	m, err := NewIDManager(s)
	if err != nil {
		t.Error(err)
	}

	id2, err := m.GetID(s)
	if err != nil {
		t.Error(err)
	}

	if id2 != id {
		t.Errorf("invalid result: %+v", s)
	}
}
