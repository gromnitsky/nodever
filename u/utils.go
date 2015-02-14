// Do `u.Conf = mymap` before using any routing from the pkg.
package u

import (
	"fmt"
	"os"
)

var Conf map[string]interface{}

func Veputs(level int, format string, args ...interface{}) {
	if *Conf["verbose"].(*int) >= level {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}

func Warnx(format string, args ...interface{}) {
	if *Conf["verbose"].(*int) >= 1 {
		fmt.Fprintf(os.Stderr, Conf["name"].(string) + " warning: " + format + "\n", args...)
	}
}

func Errx(exit_code int, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, Conf["name"].(string) + " error: " + format + "\n", args...)
	if exit_code > 0 {
		os.Exit(exit_code)
	}
}
