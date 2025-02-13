package debug

import (
	"log"
	"os"
)

var d bool

func init() {
	d = os.Getenv("DEBUG") == "1"
}

func Printf(format string, v ...any) {
	if d {
		log.Printf(format, v...)
	}
}

func Println(v ...any) {
	if d {
		log.Println(v...)
	}
}
