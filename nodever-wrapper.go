package main

import (
	"fmt"
	"os"
	"path"
	"flag"
	"strings"
	"encoding/json"
	"regexp"
	"io"

	"./u"
)

var conf = map[string]interface{} {
	"name": "nodever-wrapper",
	"wrapper": path.Base(os.Args[0]),
	"verbose": flag.Int("v", 0, "verbose level"),
	"config": flag.String("config", ".nodever.json", "debug"),
}

type NodeInfo struct {
	Dir string
	Def string
}

func slocal() (cmd string, err error) {
	file := *conf["config"].(*string)
	const sep = string(os.PathSeparator)

	pwd, _ := os.Getwd()
	arr := strings.Split(pwd, sep)
	dir := ""

	for idx := len(arr)-1; idx >= 0; idx-- {
		cur := strings.Join(append(arr[0:idx+1], file), sep)
//		fmt.Printf("%v\n", cur)

		fd, err := os.Open(cur)
		if err != nil {
			u.Veputs(2, "%s\n", err.Error())
			continue
		}
		defer fd.Close()

		json := json_parse(fd, cur)
		if dir = node_path(json); dir != "" { break }
	}

	if dir == "" { err = fmt.Errorf("cannot get node path from `%s`", file) }
	return cmd_get(dir), err
}

func json_parse(reader io.Reader, cur string) *NodeInfo {
	var nodeinfo NodeInfo

	if err := json.NewDecoder(reader).Decode(&nodeinfo); err != nil {
		u.Warnx(cur + ": " + err.Error())
		return nil
	}
	if !json_validate(&nodeinfo) {
		u.Warnx(cur + ": invalid values in the config")
		return nil
	}

	return &nodeinfo
}

func json_validate(json *NodeInfo) bool {
	arr := []string{json.Dir, json.Def}
	for _,str := range arr {
		if m, _ := regexp.MatchString("^\\s*$", str); m { return false }
	}

	return true
}

func node_path(json *NodeInfo) string {
	if json == nil { return "" }
	dir := path.Join(json.Dir, json.Def)
	if _, err := os.Stat(dir); os.IsNotExist(err) {	return "" }
	return dir
}

func cmd_get(dir string) string {
	if dir == "" { return "" }
	return path.Join(dir, "bin", conf["wrapper"].(string))
}

func senv() (cmd string, err error) {
	err = fmt.Errorf("cannot get node path from NODEVER env var")
	return
}

func main() {
	u.Conf = conf
	flag.Parse()

	variants := []func() (string, error) { senv, slocal }
	var cmd string
	var err error

	for _,fnc := range variants {
		if cmd, err = fnc(); err == nil {
			u.Veputs(2, "FOUND: %s\n", cmd)
			break
		} else {
			u.Warnx("%s", err)
		}
	}

	if cmd == "" {
		u.Errx(1, "noooooooo!")
	}
}
