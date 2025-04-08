package debug

import (
	"log"
)

var debug bool

func Enable() {
	debug = true
}

func Printf(format string, v ...any) {
	if debug {
		log.Printf(format, v...)
	}
}

func Println(v ...any) {
	if debug {
		log.Println(v...)
	}
}
