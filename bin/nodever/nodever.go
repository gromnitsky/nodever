package main

import (
	"os"
	"os/user"
	"flag"
	"fmt"
	"path"
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/gromnitsky/nodever/u"
	"github.com/gromnitsky/nodever/nodeinfo"
	"github.com/gromnitsky/nodever/meta"
)

var conf = map[string]interface{} {
	"name": path.Base(os.Args[0]),
	"config_var": "NODEVER",
	"verbose": flag.Int("v", 0, "verbose level"),
	"config": flag.String("config", ".nodever.json", "debug"),
	"useronly": flag.Bool("u", false, "start from $HOME"),
	"version": flag.Bool("V", false, "version info"),
}

func mode_info() {
	var ni *nodeinfo.NodeInfo
	var err error
	var source string
	variants := []nodeinfo.Finder {
		&nodeinfo.DataVar{*&nodeinfo.DataFile{conf["config_var"].(string)}},
		&nodeinfo.DataFile{*conf["config"].(*string)},
	}
	for _,data := range variants {
		if source, ni, err = data.Dirname(); err == nil {
			fmt.Printf("%s (%s)\n", ni.Def, source)
			break
		} else {
			u.Warnx("%s", err)
		}
	}

	if ni == nil {
		u.Errx(66, "cannot find node; rerun w/ '-v 1' argument; for help see %s",
			meta.Website)
	}
}

func mode_init() {
	var config = *conf["config"].(*string)
	ni := &nodeinfo.NodeInfo{Dir: "/opt/s", Def: "SET ME"}
	config_write(config, ni)
}

func config_write(filename string, ni *nodeinfo.NodeInfo) {
	u.Veputs(1, "writing %s\n", filename)
	fd, err := os.Create(filename)
	if err != nil {
		u.Errx(1, err.Error())
	}
	defer fd.Close()

	if err := json.NewEncoder(fd).Encode(ni); err != nil {
		u.Errx(1, err.Error())
	}
}

func mode_list() {
	var list []NodeVersion
	var err error
	var source string
	if source, _, list, err = node_versions(); err != nil {
		u.Errx(1, "cannot read config file, run `%s init`: %s", conf["name"], err)
	}

	fmt.Printf("(%s)\n", source)
	for _, val := range list {
		fmt.Println(&val)
	}
}

type NodeVersion struct {
	name string
	is_cur bool
}

func (nv *NodeVersion) String() string {
	// where is ternary operator? why golang, why?
	mark := " "
	if nv.is_cur { mark = "*"}
	return mark + " " + nv.name
}

// get dir from config, return all node subdirs from it
func node_versions() (source string, ni *nodeinfo.NodeInfo, list []NodeVersion, err error) {
	df := &nodeinfo.DataFile{*conf["config"].(*string)}
	if source, ni, err = df.Dirname(); err != nil { return }

	files, err := ioutil.ReadDir(ni.Dir)
	for _, file := range files {
		if is_node_dir(file.Name()) {
			is_cur := false
			if ni.Def == file.Name() { is_cur = true }
			list = append(list, NodeVersion{file.Name(), is_cur})
		}
	}

	return
}

// well, in ruby it would have been much prettier
func node_versions_filter(list []NodeVersion, str string) (r []NodeVersion) {
	if str == "" { return list }
	for _, val := range list {
		if strings.Contains(val.name, str) {
			r = append(r, val)
		}
	}
	return
}

func is_node_dir(name string) bool {
	m, _ := regexp.MatchString("^(node|iojs)-v?\\d+\\.\\d+\\.\\d+", name)
	return m
}

func mode_use(filter string) {
	source, ni := select_version(filter)
	config_write(source, ni)
}

func select_version(filter string) (source string, ni *nodeinfo.NodeInfo) {
	var list []NodeVersion
	var err error
	if source, ni, list, err = node_versions(); err != nil {
		u.Errx(1, "cannot read config file, run `%s init`: %s", conf["name"], err)
	}
	ver := node_versions_filter(list, filter)
	if len(ver) > 1 {
		u.Errx(0, "the query must resolve in 1 entry, you got:\n")
		for _,val := range ver {
			fmt.Fprintln(os.Stderr, &val)
		}
		os.Exit(1)
	}
	if len(ver) == 0 { u.Errx(1, "`%s` doesn't match any node installations in %s", filter, ni.Dir) }

	ni.Def = ver[0].name
	return
}

func main() {
	u.Conf = conf
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [mode [arg]]:\n", conf["name"])
		fmt.Fprintf(os.Stderr, "Modes:\n" +
			"  init       create/overwrite a config in the current dir or $HOME (see -u)\n" +
			"  list       list node/iojs installations\n" +
			"  use VER    select node version (write it to the config)\n" +
			"  (no mode)  show currrent node version & its source\n" +
			"\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	if *conf["version"].(*bool) {
		fmt.Println(meta.Version)
		os.Exit(0)
	}

	if *conf["useronly"].(*bool) {
		u.Veputs(1, "cd $HOME\n")
		// TODO: check errors
		usr, _ := user.Current()
		os.Chdir(usr.HomeDir)
	}

	switch flag.Arg(0) {
	case "init":
		mode_init()
	case "list":
		mode_list()
	case "use":
		mode_use(flag.Arg(1))
	case "":					// yep
		mode_info()
	default:
		u.Errx(1, "unknown mode")
	}
}
