package utils

import (
	"fmt"
	"log"
	"strings"
)

const (
	// Bold text
	initColor = "\033[1m"

	// Reset color
	endColor = "\033[0m"
)

// Print title in bold inside box with an orange background color
func PrintTitle(title string) {
	log.Println()
	tittleLine := strings.Repeat("─", len(title)+2)
	printlnBoldBgText("  ┌" + tittleLine + "┐")
	printlnBoldBgText("  │ %s │", title)
	printBoldBgText("  └" + tittleLine + "┘")
	log.Println()
}

func printBoldBgText(format string, args ...any) {
	text := fmt.Sprintf(format, args...)
	log.Printf("%s%s%s", initColor, text, endColor)
}

func printlnBoldBgText(format string, args ...any) {
	text := fmt.Sprintf(format, args...)
	log.Printf("%s%s%s\n", initColor, text, endColor)
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
