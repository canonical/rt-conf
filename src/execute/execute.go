package execute

func GrubConclusion(grubFile string) []string {
	s := []string{
		"Detected bootloader: GRUB\n",
		"Updated default grub file: " + grubFile + "\n",
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
		"Detected bootloader: Raspberry Pi\n",
		"\n",
		"Please, append the following to /boot/firmware/cmdline.txt:\n",
		"In case of old style boot partition,\n",
		"append to /boot/cmdline.txt\n",
		cmdline,
		"\n",
	}
	return s
}

func UbuntuCoreConclusion(change string) []string {
	s := []string{
		"Detected bootloader: Ubuntu Core managed\n",
		"\n",
		"Sucessfully applied the changes.\n",
		"Snapd change: " + change + "\n",
		"Please reboot your system to apply the changes.\n",
	}
	return s
}
