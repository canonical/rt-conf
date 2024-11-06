package validator

import (
	"fmt"
	"os"
	"reflect"
)

var TypeEnum = map[string]reflect.Type{
	"string": reflect.TypeOf(""),
	"bool":   reflect.TypeOf(true),
	"int":    reflect.TypeOf(0),
}

// validateType checks if the parameter matches the expected type.
func ValidateType(expected reflect.Type, paramName string, param interface{}) error {
	// Check if the value's type matches the expected type
	actualType := reflect.TypeOf(param)
	if actualType != expected {
		os.Stderr.WriteString(fmt.Sprintf("bad YAML format: parameter '%s' is of type '%s', expected '%s'\n", paramName, actualType, expected))
		os.Exit(1)
	}
	return nil
}
