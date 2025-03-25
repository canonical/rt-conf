package model

import (
	"fmt"
	"reflect"
)

// Generic validation function for structs
func Validate(c any, fieldValidator func(string, string, string) error) error {
	// Validate fields based on struct tags
	v := reflect.ValueOf(c)
	t := reflect.TypeOf(c)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		tag := field.Tag.Get("validation")
		if tag == "" {
			continue // No validation tag, skip
		}

		// * * NOTE: For now it's okay to cast to string
		// * * since ther is only strings on KernelCmdline struct
		value, ok := v.Field(i).Interface().(string)
		if !ok {
			return fmt.Errorf("value for field %s is not a string", value)
		}

		if value == "" {
			continue
		}
		err := fieldValidator(field.Name, value, tag)
		if err != nil {
			return fmt.Errorf("validation failed for field %s: %v",
				field.Name, err)
		}
	}
	return nil
}
