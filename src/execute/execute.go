package execute

import "fmt"

func GrubConclusion() {
	// TODO: Add system detection functionality to print the message for each system
	fmt.Println("Successfully injected to file")
	fmt.Println("Please run:")
	fmt.Println("")
	fmt.Println("sudo update-grub")
	fmt.Println("")
	fmt.Println("to apply the changes")
}
