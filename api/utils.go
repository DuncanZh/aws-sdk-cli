package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetDirectory(prefix string) []string {
	path := "./" + prefix
	if f, err := os.Stat(path); err == nil {
		if !f.IsDir() {
			return []string{prefix}
		}
	} else {
		path = path[:strings.LastIndex(path, "/")] + "/"
	}

	entries, _ := os.ReadDir(path)

	var es []string
	for _, e := range entries {
		if path == "./" {
			es = append(es, e.Name())
		} else {
			es = append(es, filepath.Dir(path[2:]+"/")+"/"+e.Name())
		}
	}
	return es
}

func DumpFile(result interface{}, file string, pretty bool) bool {
	j, err := json.Marshal(result)
	if pretty {
		j, err = json.MarshalIndent(result, "", " ")
	}

	if err != nil {
		PrintError(err)
		return false
	}

	err = os.MkdirAll(filepath.Dir(file), 0777)
	if err != nil {
		PrintError(err)
		return false
	}

	f, err := os.Create(file)
	if err != nil {
		PrintError(err)
		return false
	}
	defer f.Close()

	_, err = f.Write(j)

	if err != nil {
		PrintError(err)
		return false
	}
	return true
}

func PrintError(err error) {
	fmt.Printf("\nError: %s\n", err)
}

func PrintErrorString(str string) {
	fmt.Printf("\nError: %s\n", errors.New(str))
}
