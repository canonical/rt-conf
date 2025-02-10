package config

type Views int

const (
	INIT_VIEW_ID Views = iota
	KCMD_VIEW_ID
	KCMD_CONCLUSION_VIEW_ID
	IRQ_VIEW_ID
	IRQ_ADD_EDIT_VIEW_ID
	IRQ_CONCLUSION_VIEW_ID
	PWMGMT_VIEW_ID
)

// Map of views to their names (for logging)
var Menu map[Views]string = map[Views]string{
	INIT_VIEW_ID:            "INIT VIEW: Main Menu",
	KCMD_VIEW_ID:            "Kcmd Menu",
	KCMD_CONCLUSION_VIEW_ID: "Kcmd Conclusion",
	IRQ_VIEW_ID:             "IRQ Affinity",
	IRQ_ADD_EDIT_VIEW_ID:    "Inner IRQ Affinity (add/edit)",
	PWMGMT_VIEW_ID:          "Power Management",
}

// Init menu view names
const (
	MENU_KCMDLINE    = "Kernel cmdline"
	MENU_IRQAFFINITY = "IRQ Affinity"
	MENU_PWRMGMT     = "Power Management"
)

const (
	DESC_KCMDLINE    = "Configure Kernel cmdline parameters"
	DESC_IRQAFFINITY = "Isolate CPUs from serving IRQs"
	DESC_PWRMGMT     = "Configure CPU power management settings"
)

const (
	NUMBER_OF_MENUS = 3
)

const CpuListPlaceholder = "Enter a CPU list like: 4-n or 3-5 or 2,4,5 "
const IrqFilterPlaceholder = "Insert filter parameters for IRQs"

const (
	PrefixIRQFilter = "Filter > "
	PrefixCpuList   = "CPU Range > "
)
