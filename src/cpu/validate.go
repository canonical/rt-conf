package cpu

func ValidateList(s string, max int) error {
	_, err := ParseCPUs(s, max)
	return err
}
