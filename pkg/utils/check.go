package utils

import (
	"errors"
	"reflect"
	"strings"
)

func ValidateStruct(model interface{}) error {
	val := reflect.ValueOf(model)

	// Если указатель → разыменовываем
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Проверка только struct
	if val.Kind() != reflect.Struct {
		return errors.New("ValidateStruct: expected struct")
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		// Проверяем только string-поля
		if field.Kind() == reflect.String {
			value := strings.TrimSpace(field.String())
			if value == "" {
				return errors.New("field " + typ.Field(i).Name + " is required")
			}
		}
	}

	return nil
}
