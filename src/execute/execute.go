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

func ExecRaspberry(cmdline []string) {
	fmt.Println("Please, append the following to /boot/firmware/cmdline.txt:")
	fmt.Printf("In case of old style boot partition, \nappend to /boot/cmdline.txt\n\n")
	for _, param := range cmdline {
		fmt.Printf("%s ", param)
	}
	fmt.Printf("\n\n")
}
