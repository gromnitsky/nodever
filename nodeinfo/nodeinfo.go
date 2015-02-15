package nodeinfo

import (
	"fmt"
	"os"
	"path"
	"strings"
	"encoding/json"
	"regexp"
	"io"

	"../u"
)

type Finder interface {
	Dirname() (string, error)
}


type DataFile struct {
	Name string
}

func (df *DataFile) dirname_failure() error {
	return fmt.Errorf("cannot get node path from `%s`", df.Name)
}

type DataVar struct {
	DataFile
}

func (df *DataFile) Dirname() (dir string, err error) {
	const sep = string(os.PathSeparator)

	pwd, _ := os.Getwd()
	arr := strings.Split(pwd, sep)

	for idx := len(arr)-1; idx >= 0; idx-- {
		cur := strings.Join(append(arr[0:idx+1], df.Name), sep)
//		fmt.Printf("%v\n", cur)

		fd, err := os.Open(cur)
		if err != nil {
			u.Veputs(1, "%s\n", err.Error())
			continue
		}
		defer fd.Close()

		json := json_parse(fd, &cur)
		if dir = node_path(json); dir != "" { break }
	}

	if dir == "" { err = df.dirname_failure() }
	return
}

func (dv *DataVar) Dirname() (dir string, err error) {
	env := strings.NewReader(os.Getenv(dv.Name))
	json := json_parse(env, &dv.Name)
	if dir = node_path(json); dir == "" {
		err = dv.dirname_failure()
	}
	return
}

type nodeInfo struct {
	Dir string
	Def string
}

func json_parse(reader io.Reader, src *string) *nodeInfo {
	var nodeinfo nodeInfo

	if err := json.NewDecoder(reader).Decode(&nodeinfo); err != nil {
		u.Warnx(*src + ": " + err.Error())
		return nil
	}
	if !json_validate(&nodeinfo) {
		u.Warnx(*src + ": invalid values in the config")
		return nil
	}

	return &nodeinfo
}

func json_validate(json *nodeInfo) bool {
	arr := []string{json.Dir, json.Def}
	for _,str := range arr {
		if m, _ := regexp.MatchString("^\\s*$", str); m { return false }
	}

	return true
}

func node_path(json *nodeInfo) string {
	if json == nil { return "" }
	dir := path.Join(json.Dir, json.Def)
	if _, err := os.Stat(dir); os.IsNotExist(err) {	return "" }
	return dir
}