package cpu

import "fmt"

func ValidateList(s string) error {

	max, err := TotalAvailable()
	if err != nil {
		return fmt.Errorf("failed to get total available CPUs: %v", err)
	}
	_, err = ParseCPUs(s, max)
	return err
}
