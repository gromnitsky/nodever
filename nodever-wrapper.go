package main

import (
	"os"
	"os/exec"
	"syscall"
	"path"
	"flag"
	"fmt"

	"github.com/mattn/go-shellwords"

	"./u"
	"./nodeinfo"
	"./meta"
)

var conf = map[string]interface{} {
	"name": "nodever-wrapper",
	"wrapper": path.Base(os.Args[0]),
	"wrapper_env": "NODEVER_WRAPPER",
	"config_var": "NODEVER",

	// in NODEVER_WRAPPER env var, not the command line
	"verbose": flag.Int("v", 0, "verbose level"),
	"config": flag.String("config", ".nodever.json", "debug"),
	"version": flag.Bool("V", false, "version info"),
}

func parse_debug_env() (err error) {
	args, err := shellwords.Parse(os.Getenv(conf["wrapper_env"].(string)))
	if err != nil {
		u.Errx(1, "%s: invalid shell words", conf["wrapper_env"].(string))
	}
	flag.CommandLine.Parse(args)
	return
}

func run(program string, args []string) {
	cmd := exec.Command(program, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	u.Veputs(1, "RUN: %s %s\n", program, args)
	if err := cmd.Run(); err != nil {
		if cmd.ProcessState != nil {
			// exit w/ captured exit status
			os.Exit(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus())
		} else {
			// fork error
			u.Errx(65, "%s", err)
		}
	}

	os.Exit(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus())
}

func main() {
	u.Conf = conf
	parse_debug_env()
	if *conf["version"].(*bool) {
		fmt.Println(meta.Version)
		os.Exit(0)
	}

	var dir string
	var err error
	variants := []nodeinfo.Finder {
		&nodeinfo.DataVar{*&nodeinfo.DataFile{conf["config_var"].(string)}},
		&nodeinfo.DataFile{*conf["config"].(*string)},
	}
	for _,data := range variants {
		if dir, err = data.Dirname(); err == nil {
			u.Veputs(1, "FOUND: %s\n", dir)
			run(path.Join(dir, "bin", conf["wrapper"].(string)), os.Args[1:])
			break
		} else {
			u.Warnx("%s", err)
		}
	}

	u.Errx(66, "cannot find node; run\n\n  $ %s='-v 1' %s\n\nfor more info; " +
		"for help see %s",
		conf["wrapper_env"], conf["wrapper"], meta.Website)
}
