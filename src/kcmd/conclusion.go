package kcmd

func GrubConclusion(grubFile string) []string {
	s := []string{
		"Detected bootloader: GRUB\n",
		"Updated default grub file: " + grubFile + "\n",
		"\n",
		"Please run:\n",
		"\n",
		"\tsudo update-grub\n",
		"\n",
		"to apply the changes to your bootloader.\n",
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

func UbuntuCoreConclusion() []string {
	s := []string{
		"Detected bootloader: Ubuntu Core managed\n",
		"\n",
		"Sucessfully applied the changes.\n",
		"Please reboot your system to apply the changes.\n",
	}
	return s
}
