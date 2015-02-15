package main

import (
	"os"
	"path"
	"flag"

	"github.com/mattn/go-shellwords"

	"./u"
	"./nodeinfo"
)

var conf = map[string]interface{} {
	"name": "nodever-wrapper",
	"wrapper": path.Base(os.Args[0]),
	"wrapper_env": "NODEVER_WRAPPER",
	"config_var": "NODEVER",

	// in NODEVER_WRAPPER env var, not the command line
	"verbose": flag.Int("v", 0, "verbose level"),
	"config": flag.String("config", ".nodever.json", "debug"),
}

func parse_debug_env() (err error) {
	args, err := shellwords.Parse(os.Getenv(conf["wrapper_env"].(string)))
	if err != nil {
		u.Errx(1, "%s: invalid shell words", conf["wrapper_env"].(string))
	}
	flag.CommandLine.Parse(args)
	return
}


func main() {
	u.Conf = conf
	parse_debug_env()

	var dir string
	var err error
	variants := []nodeinfo.Finder {
		&nodeinfo.DataVar{*&nodeinfo.DataFile{conf["config_var"].(string)}},
		&nodeinfo.DataFile{*conf["config"].(*string)},
	}
	for _,data := range variants {
		if dir, err = data.Dirname(); err == nil {
			u.Veputs(1, "FOUND: %s\n",
				path.Join(dir, "bin", conf["wrapper"].(string)))
			break
		} else {
			u.Warnx("%s", err)
		}
	}

	if dir == "" {
		u.Errx(1, "NOOOOOOOO!")
	}
}
