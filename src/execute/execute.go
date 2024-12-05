package execute

import "fmt"

func Exec() {
	// TODO: Add system detection functionality to print the message for each system
	// TODO: Move the message to a separate function
	fmt.Println("Successfully injected to file")
	fmt.Println("Please run:\nsudo update-grub\nto apply the changes")
}
