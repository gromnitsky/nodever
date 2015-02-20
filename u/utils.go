// Do `u.Conf = mymap` before using any routing from the pkg.
package u

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
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

func Run(program string, args []string) {
	cmd := exec.Command(program, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	Veputs(1, "RUN: %s %s\n", program, args)
	if err := cmd.Run(); err != nil && cmd.ProcessState == nil {
		// fork error
		Errx(65, "%s", err)
	}

	os.Exit(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus())
}
