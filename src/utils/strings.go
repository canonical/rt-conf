// Package utils provides utility functions used across the application
package utils

import (
	"strings"
)

func TrimSurroundingQuotes(value string) string {
	if len(value) >= 2 {
		if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) ||
			strings.HasPrefix(value, `'`) && strings.HasSuffix(value, `'`) {
			return value[1 : len(value)-1]
		}
	}
	return value
}

func TrimSurroundingDoubleQuotes(value string) string {
	if len(value) >= 2 {
		// on https://docs.kernel.org/admin-guide/kernel-parameters.html#special-handling
		// kernel docs mentions that double-quotes can be used to protect spaces on the value
		if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
			return value[1 : len(value)-1]
		}
	}
	return value
}
