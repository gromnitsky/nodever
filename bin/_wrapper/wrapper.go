package main

import (
	"os"
	"path"
	"flag"
	"fmt"
	"encoding/json"

	"github.com/mattn/go-shellwords"

	"github.com/gromnitsky/nodever/u"
	"github.com/gromnitsky/nodever/nodeinfo"
	"github.com/gromnitsky/nodever/meta"
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

func set_nodever(ni *nodeinfo.NodeInfo) {
	if os.Getenv(conf["config_var"].(string)) != "" { return }

	b, _ := json.Marshal(ni)
	u.Veputs(1, "set %s=%s for possible subshells\n",
		conf["config_var"].(string), string(b))
	os.Setenv(conf["config_var"].(string), string(b))
}

func main() {
	u.Conf = conf
	parse_debug_env()
	if *conf["version"].(*bool) {
		fmt.Println(meta.Version)
		os.Exit(0)
	}

	var ni *nodeinfo.NodeInfo
	var err error
	variants := []nodeinfo.Finder {
		&nodeinfo.DataVar{*&nodeinfo.DataFile{conf["config_var"].(string)}},
		&nodeinfo.DataFile{*conf["config"].(*string)},
	}
	for idx,data := range variants {
		if _, ni, err = data.Dirname(); err == nil {
			u.Veputs(1, "FOUND/%d: %s\n", idx, ni.Def)
			set_nodever(ni)
			u.Run(path.Join(ni.Dir, ni.Def, "bin", conf["wrapper"].(string)), os.Args[1:])
			break
		} else {
			u.Warnx("%s", err)
		}
	}

	u.Errx(66, "cannot find node; run\n\n  $ %s='-v 1' %s\n\nfor more info; " +
		"for help see %s",
		conf["wrapper_env"], conf["wrapper"], meta.Website)
}
