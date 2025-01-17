package execute

import (
	"fmt"
)

func GrubConclusion() []string {
	s := []string{
		"\n",                   // 1
		"Please run:\n",        // 2
		"\n",                   // 3
		"\tsudo update-grub\n", // 4
		"\n",                   // 5
		"to apply the changes to your bootloader.\n", // 6
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

	return s
}
