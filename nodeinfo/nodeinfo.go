package nodeinfo

import (
	"fmt"
	"os"
	"strings"
	"encoding/json"
	"regexp"
	"io"

	"github.com/gromnitsky/nodever/u"
)

type Finder interface {
	Dirname() (string, *NodeInfo, error)
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

func (df *DataFile) Dirname() (source string, ni *NodeInfo, err error) {
	const sep = string(os.PathSeparator)

	pwd, _ := os.Getwd()
	arr := strings.Split(pwd, sep)

	for idx := len(arr)-1; idx >= 0; idx-- {
		cur := strings.Join(append(arr[0:idx+1], df.Name), sep)

		fd, err := os.Open(cur)
		if err != nil {
			u.Veputs(1, "%s\n", err.Error())
			continue
		}
		defer fd.Close()

		if ni = json_parse(fd, &cur); ni != nil {
			source = cur
			break
		}
	}

	if ni == nil { err = df.dirname_failure() }
	return
}

func (dv *DataVar) Dirname() (source string, ni *NodeInfo, err error) {
	source = dv.Name + " env var"
	env := strings.NewReader(os.Getenv(dv.Name))
	if ni = json_parse(env, &dv.Name); ni == nil {
		err = dv.dirname_failure()
	}
	return
}

type NodeInfo struct {
	Dir string
	Def string
}

func json_parse(reader io.Reader, src *string) *NodeInfo {
	var nodeinfo NodeInfo

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

func json_validate(json *NodeInfo) bool {
	arr := []string{json.Dir, json.Def}
	for _,str := range arr {
		if m, _ := regexp.MatchString("^\\s*$", str); m { return false }
	}

	return true
}
