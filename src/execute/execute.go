package execute

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

func RpiConclusion(cmdline string) []string {
	s := []string{
		"\n", //1
		"Please, append the following to /boot/firmware/cmdline.txt:\n", // 2
		"In case of old style boot partition,\n",                        // 3
		"append to /boot/cmdline.txt\n",                                 // 4
		cmdline,                                                         // 5
		"\n",                                                            // 6
	}
	return s
}

func UbuntuCoreConclusion(change string) []string {
	s := []string{
		"\n",                                 // 1
		"Sucessfully applied the changes.\n", // 2
		"Snapd change: " + change + "\n",     // 3
		"Please reboot your system to apply the changes.\n", // 4
	}
	return s
}
