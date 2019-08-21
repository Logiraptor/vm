package debug

import (
	"fmt"
	"os"
)

const Debug = false

func Debugln(args ...interface{}) {
	if Debug {
		fmt.Fprintln(os.Stdout, args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if Debug {
		fmt.Fprintf(os.Stdout, format, args...)
	}
}
