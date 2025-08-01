package db

import (
	"errors"
	"reflect"
	"strings"
)

// Use 'db' tag to to match the name of a structure field with the name of a column from the database
// mapToStruct конвертирует map в struct через reflect
func MapToStruct(m map[string]interface{}, out interface{}) error {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return errors.New("out must be a non-nil pointer to struct")
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("out must be pointer to struct")
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}

		colName := field.Tag.Get("db")
		if colName == "" {
			colName = toSnakeCase(field.Name)
		}

		if val, ok := m[colName]; ok {
			fv := v.Field(i)
			if !fv.CanSet() {
				continue
			}

			valRef := reflect.ValueOf(val)
			if valRef.Type().AssignableTo(fv.Type()) {
				fv.Set(valRef)
			} else if valRef.Type().ConvertibleTo(fv.Type()) {
				fv.Set(valRef.Convert(fv.Type()))
			} else {
				// Попытка конвертировать []byte в string если поле string
				if b, ok := val.([]byte); ok && fv.Kind() == reflect.String {
					fv.SetString(string(b))
				}
			}
		}
	}

	return nil
}

// structToMap конвертирует struct в map[string]interface{} через reflect
func StructToMap(entity interface{}) (map[string]interface{}, error) {
	v := reflect.ValueOf(entity)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, errors.New("entity must be a struct or pointer to struct")
	}

	m := make(map[string]interface{})
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}

		// Берём имя колонки из тега db, если есть, иначе имя поля в snake_case
		colName := field.Tag.Get("db")
		if colName == "" {
			colName = toSnakeCase(field.Name)
		}
		m[colName] = v.Field(i).Interface()
	}

	return m, nil
}

// toSnakeCase простой конвертер CamelCase в snake_case
func toSnakeCase(s string) string {
	var sb strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			sb.WriteByte('_')
		}
		sb.WriteRune(r)
	}
	return strings.ToLower(sb.String())
}
