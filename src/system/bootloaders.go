package system

type Bootloader int

const (
	Unknown Bootloader = iota
	Grub
	Rpi
	Uboot
)
