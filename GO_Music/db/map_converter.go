package db

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// [RU] MapToStruct конвертирует map в struct через reflect
// [ENG]MapToStruct convert map[string]interface{} into struct with reflect;
// Use 'db' tag to to match the name of a structure field with the name of a column from the database;
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
		if field.PkgPath != "" { // unexported
			continue
		}

		colName := field.Tag.Get("db")
		if colName == "" {
			colName = ToSnakeCase(field.Name)
		}

		if val, ok := m[colName]; ok {
			fv := v.Field(i)
			if !fv.CanSet() {
				continue
			}

			// Обработка nil значений
			if val == nil {
				if fv.Kind() == reflect.Ptr {
					fv.Set(reflect.Zero(fv.Type()))
				}
				continue
			}

			// Обработка указателей
			if fv.Kind() == reflect.Ptr {
				// Создаем новый указатель нужного типа
				ptrVal := reflect.New(fv.Type().Elem())

				// Устанавливаем значение
				if err := setFieldValue(ptrVal.Elem(), val); err != nil {
					return fmt.Errorf("field %s: %v", field.Name, err)
				}

				fv.Set(ptrVal)
			} else {
				// Обработка обычных полей
				if err := setFieldValue(fv, val); err != nil {
					return fmt.Errorf("field %s: %v", field.Name, err)
				}
			}
		}
	}

	return nil
}

// [RU]StructToMap конвертирует struct в map[string]interface{} через reflect;
// [ENG]Struct To Map convert struct into map[string]interface{} with reflect
func StructToMap(entity interface{}) (map[string]interface{}, error) {
	v := reflect.ValueOf(entity)
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil, errors.New("nil pointer passed")
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, errors.New("entity must be a struct or pointer to struct")
	}

	m := make(map[string]interface{})
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		// Пропускаем непубличные поля
		if field.PkgPath != "" {
			continue
		}

		colName := field.Tag.Get("db")
		if colName == "" {
			colName = ToSnakeCase(field.Name)
		}

		fieldVal := v.Field(i)

		if fieldVal.Kind() == reflect.Pointer {
			if fieldVal.IsNil() {
				m[colName] = nil
				continue
			}
			// Для *string кладём сам указатель (например, *string с валидным значением)
			if fieldVal.Type() == reflect.TypeOf((*string)(nil)) {
				m[colName] = fieldVal.Interface()
				continue
			}
			// Для других указателей разыменовываем
			fieldVal = fieldVal.Elem()
		}

		m[colName] = fieldVal.Interface()
	}

	return m, nil
}

// [RU] ToSnakeCase простой конвертер CamelCase в snake_case
// [ENG] ToSnakeCase is a simple converter from 'CamelCase' into 'snake_case'
func ToSnakeCase(s string) string {
	var sb strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]

		if c >= 'A' && c <= 'Z' {
			if i > 0 {
				prev := s[i-1]
				nextLower := false
				if i < len(s)-1 {
					next := s[i+1]
					nextLower = next >= 'a' && next <= 'z'
				}

				if (prev >= 'a' && prev <= 'z') ||
					(prev >= '0' && prev <= '9') ||
					nextLower {
					sb.WriteByte('_')
				}
			}
			sb.WriteByte(c + 32)
		} else {
			sb.WriteByte(c)
		}
	}
	return sb.String()
}

// setFieldValue устанавливает значение в reflect.Value с учетом типов
func setFieldValue(fv reflect.Value, val interface{}) error {
	valRef := reflect.ValueOf(val)

	// Если типы совместимы напрямую
	if valRef.Type().AssignableTo(fv.Type()) {
		fv.Set(valRef)
		return nil
	}

	// Если можно конвертировать
	if valRef.Type().ConvertibleTo(fv.Type()) {
		fv.Set(valRef.Convert(fv.Type()))
		return nil
	}

	// Обработка специальных случаев
	switch fv.Kind() {
	case reflect.String:
		switch v := val.(type) {
		case []byte:
			fv.SetString(string(v))
		case fmt.Stringer:
			fv.SetString(v.String())
		default:
			fv.SetString(fmt.Sprintf("%v", val))
		}
		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v := val.(type) {
		case int64:
			fv.SetInt(v)
		case float64:
			fv.SetInt(int64(v))
		case float32:
			fv.SetInt(int64(v))
		case int:
			fv.SetInt(int64(v))
		default:
			return fmt.Errorf("cannot convert %T to int", val)
		}
		return nil

	case reflect.Float32, reflect.Float64:
		switch v := val.(type) {
		case float64:
			fv.SetFloat(v)
		case int64:
			fv.SetFloat(float64(v))
		case int:
			fv.SetFloat(float64(v))
		default:
			return fmt.Errorf("cannot convert %T to float", val)
		}
		return nil

	case reflect.Bool:
		switch v := val.(type) {
		case bool:
			fv.SetBool(v)
		case int64:
			fv.SetBool(v != 0)
		case int:
			fv.SetBool(v != 0)
		default:
			return fmt.Errorf("cannot convert %T to bool", val)
		}
		return nil

	case reflect.Struct:
		if fv.Type() == reflect.TypeOf(time.Time{}) {
			switch v := val.(type) {
			case time.Time:
				fv.Set(reflect.ValueOf(v))
			case string:
				t, err := time.Parse(time.RFC3339, v)
				if err != nil {
					return fmt.Errorf("cannot parse time: %v", err)
				}
				fv.Set(reflect.ValueOf(t))
			default:
				return fmt.Errorf("cannot convert %T to time.Time", val)
			}
			return nil
		}
	}

	return fmt.Errorf("cannot set field of type %s with value %v (%T)",
		fv.Type(), val, val)
}
