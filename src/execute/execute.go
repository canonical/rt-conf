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
		"\n",                   // 1
		"Please run:\n",        // 2
		"\n",                   // 3
		"\tsudo update-grub\n", // 4
		"\n",                   // 5
		"to apply the changes to your bootloader.\n", // 6
	}
	if len(s) != lineNums {
		panic("GrubConclusion: invalid number of lines")
	}
	return s
}

func RpiConclusion(cmdline []string) []string {
	s := []string{
		"\n", //1
		"Please, append the following to /boot/firmware/cmdline.txt:\n", // 2
		"In case of old style boot partition,\n",                        // 3
		"append to /boot/cmdline.txt\n",                                 // 4
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
