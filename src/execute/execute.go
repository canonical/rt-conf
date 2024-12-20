package execute

import (
	"fmt"
)

// TODO: Add a validation for the number of lines in each execute message
// NOTE: THESE MESSAGES MUST ALWAYS BE THE SAME NUMBER OF LINES SO IT DOESN'T
// BROKE THE TUI

// NOTE: These messages should be serialized into a struct/json for future use of TUI app
func GrubConclusion() []string {
	s := []string{
		"Successfully injected to file\n",
		"Please run:\n",
		"\n",
		"sudo update-grub\n",
		"\n",
		"to apply the changes\n",
	}
	return s
}

func RpiConclusion(cmdline []string) []string {
	s := []string{
		"Please, append the following to /boot/firmware/cmdline.txt:\n",
		"In case of old style boot partition,\n",
		"append to /boot/cmdline.txt\n\n",
	}
	for _, param := range cmdline {
		s = append(s, fmt.Sprintf("%s ", param))
	}
	return s
}
