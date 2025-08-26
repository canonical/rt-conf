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
