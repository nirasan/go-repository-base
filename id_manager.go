package go_repository_base

import (
	"errors"
	"fmt"
	"reflect"
)

const idManagerTag = "repository"

type IDManager struct {
	entityName  string
	idFieldName string
	isIntID     bool
}

func NewIDManager(e interface{}) (*IDManager, error) {
	m := &IDManager{}
	rv := reflect.ValueOf(e)
	rt := rv.Type()
	if rt.Kind() != reflect.Ptr || rt.Elem().Kind() != reflect.Struct {
		return nil, errors.New(fmt.Sprintf("Entity must struct or pointer of struct: %+v, %+v", e, rt.Kind()))
	}
	m.entityName = rt.Name()
	rt = rt.Elem()
	for i := rt.NumField() - 1; i >= 0; i = i - 1 {
		f := rt.Field(i)
		if f.Tag.Get(idManagerTag) == "id" {
			m.idFieldName = f.Name
			if f.Type.Kind() == reflect.Int64 {
				m.isIntID = true
				return m, nil
			} else if f.Type.Kind() == reflect.String {
				m.isIntID = false
				return m, nil
			}
		}
	}
	return nil, errors.New("ID field not found")
}

func (m *IDManager) SetID(e interface{}, id interface{}) error {
	if err := m.validate(e); err != nil {
		return err
	}
	rv := reflect.Indirect(reflect.ValueOf(e))
	f := rv.FieldByName(m.idFieldName)
	if m.isIntID {
		i, ok := id.(int64)
		if !ok {
			return errors.New(fmt.Sprintf("id must int64: %+v", id))
		}
		f.SetInt(i)
		return nil
	} else {
		s, ok := id.(string)
		if !ok {
			return errors.New(fmt.Sprintf("id must string: &+v", id))
		}
		f.SetString(s)
		return nil
	}
}

func (m *IDManager) GetID(e interface{}) (interface{}, error) {
	if err := m.validate(e); err != nil {
		return nil, err
	}
	rv := reflect.Indirect(reflect.ValueOf(e))
	f := rv.FieldByName(m.idFieldName)
	if m.isIntID {
		return f.Int(), nil
	} else {
		return f.String(), nil
	}
}

func (m *IDManager) validate(e interface{}) error {
	rt := reflect.TypeOf(e)
	if rt.Name() != m.entityName {
		return fmt.Errorf("invalid entity: %+v", e)
	}
	return nil
}
