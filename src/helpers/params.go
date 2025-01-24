package helpers

import (
	"fmt"
	"reflect"
)

func ReconstructKeyValuePairs(v interface{}) ([]string, error) {
	var keyValuePairs []string

	val := reflect.TypeOf(v)
	valValue := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		valValue = valValue.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		key := field.Tag.Get("yaml")
		value := valValue.Field(i).String()
		if key == "" || value == "" {
			continue
		}

		keyValuePairs = append(keyValuePairs, fmt.Sprintf("%s=%s", key, value))
	}

	return keyValuePairs, nil
}
