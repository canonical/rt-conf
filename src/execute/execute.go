package execute

import (
	"fmt"
)

// NOTE: THESE MESSAGES MUST ALWAYS BE THE SAME NUMBER OF LINES SO IT DOESN'T
// BROKE THE TUI

const lineNums = 6

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
	if len(s) != lineNums {
		panic("GrubConclusion: invalid number of lines")
	}
	return s
}

func RpiConclusion(cmdline []string) []string {
	s := []string{
		"Please, append the following to /boot/firmware/cmdline.txt:\n", // 1
		"In case of old style boot partition,\n",                        // 2
		"append to /boot/cmdline.txt\n",                                 // 3
		"\n",                                                            // 4
	}
	kcmdline := ""
	for _, param := range cmdline {
		kcmdline += fmt.Sprintf("%s ", param)
	}
	s = append(s, fmt.Sprintf("%s\n", kcmdline)) // 5
	s = append(s, "\n")                          // 6
	if len(s) != lineNums {
		panic("RpiConclusion: invalid number of lines")
	}
	return s
}
