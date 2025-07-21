package utils

import (
	"fmt"
	"log"
)

// Print title in bold inside box with an orange background color
func PrintTitle(title string) {

	printBoldBgText("┌")
	for i := 0; i <= len(title)+1; i++ {
		printBoldBgText("─")
	}
	printlnBoldBgText("┐")

	printlnBoldBgText("│ %s │", title)

	printBoldBgText("└")
	for i := 0; i <= len(title)+1; i++ {
		printBoldBgText("─")
	}
	printlnBoldBgText("┘")

	log.Println()
}

func printBoldBgText(format string, args ...any) {
	// Ubuntu orange color
	r, g, b := 0xE9, 0x54, 0x20
	text := fmt.Sprintf(format, args...)
	fmt.Printf("\033[1;48;2;%d;%d;%dm%s\033[0m", r, g, b, text)
}

func printlnBoldBgText(format string, args ...any) {
	// Ubuntu orange color
	r, g, b := 0xE9, 0x54, 0x20
	text := fmt.Sprintf(format, args...)
	fmt.Printf("\033[1;48;2;%d;%d;%dm%s\033[0m\n", r, g, b, text)
}

func LogTreeStyle(entries []string) {
	for i, entry := range entries {
		prefix := "├── "
		if i == len(entries)-1 {
			prefix = "└── "
		}
		log.Printf("%s%s\n", prefix, entry)
	}
	log.Println()
}
