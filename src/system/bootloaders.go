package system

type SystemType int

const (
	Unknown SystemType = iota
	Grub
	Rpi
	Uboot
	Core
)
