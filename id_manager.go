package go_repository_base

import (
	"errors"
	"fmt"
	"reflect"
)

const idManagerTag = "repository"

type IDManager struct {
}

func (m *IDManager) GetIDFieldName(e interface{}) (string, error) {
	rv := reflect.ValueOf(e)
	rt := rv.Type()
	if rt.Kind() != reflect.Ptr || rt.Elem().Kind() != reflect.Struct {
		return "", errors.New(fmt.Sprintf("Entity must struct or pointer of struct: %+v", e))
	}
	rt = rt.Elem()
	for i := rt.NumField() - 1; i >= 0; i = i - 1 {
		f := rt.Field(i)
		if f.Type.Kind() == reflect.Int64 && f.Tag.Get(idManagerTag) == "id" {
			return f.Name, nil
		}
	}
	return "", errors.New("not found")
}

func (m *IDManager) SetID(e interface{}, id int64) error {
	name, err := m.GetIDFieldName(e)
	if err != nil {
		return err
	}
	rv := reflect.Indirect(reflect.ValueOf(e))
	rv.FieldByName(name).SetInt(id)
	return nil
}

func (m *IDManager) GetID(e interface{}) (int64, error) {
	name, err := m.GetIDFieldName(e)
	if err != nil {
		return 0, err
	}
	rv := reflect.Indirect(reflect.ValueOf(e))
	return rv.FieldByName(name).Int(), nil
}
