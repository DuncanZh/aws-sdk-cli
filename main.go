package main

import (
	"AWS_API/api"
	"encoding/json"
	"fmt"
	"github.com/abiosoft/ishell/v2"
	"github.com/abiosoft/readline"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	VERSION := "0.1.0"
	PROMPT := "msgraph_cli v" + VERSION + " $ "

	shell := ishell.NewWithConfig(&readline.Config{Prompt: PROMPT})
	shell.SetHomeHistoryPath(".msgraph_shell_history")

	shell.Set("api", api.NewAPI())

	describeResources := []string{"instances", "reservedInstances", "volumes", "volumeStatus", "subnets",
		"securityGroups", "securityGroupRules", "natGateways", "internetGateways"}

	shell.Set("describe", describeResources)

	shell.AddCmd(&ishell.Cmd{
		Name: "describe",
		Help: "Usage: describe <resource> <output_file> [ids...]",
		CompleterWithPrefix: func(prefix string, args []string) []string {
			if len(args) == 0 {
				return shell.Get("describe").([]string)
			} else if len(args) == 1 {
				return getDirectory(prefix)
			}
			return []string{}
		},
		Func: describe,
	})

	if len(os.Args) > 1 {
		err := shell.Process(os.Args[1:]...)
		if err != nil {
			fmt.Println("Error: " + err.Error())
		}
	} else {
		shell.Println("Interactive Shell:")
		shell.Run()
		shell.Close()
	}
}

func describe(c *ishell.Context) {
	start := time.Now()

	a := c.Get("api").(*api.AwsAPI)

	var ids []string
	if len(c.Args) > 2 {
		ids = c.Args[2:]
	} else if len(c.Args) < 2 {
		fmt.Println(c.Cmd.Help)
		return
	}

	resource := c.Args[0]
	outputFile := c.Args[1]

	result := a.Describe(resource, ids)

	if result == nil {
		fmt.Printf("Failed: Unable to process the input in %.2f seconds\n", time.Since(start).Seconds())
		return
	}

	if dumpFile(result, outputFile, true) {
		fmt.Printf("Success: Processed %v entries in %.2f seconds\n", reflect.ValueOf(result).Len(), time.Since(start).Seconds())
	}
}

func getDirectory(prefix string) []string {
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

func dumpFile(result interface{}, file string, pretty bool) bool {
	j, err := json.Marshal(result)
	if pretty {
		j, err = json.MarshalIndent(result, "", " ")
	}

	if err != nil {
		api.PrintError(err.Error())
		return false
	}

	err = os.MkdirAll(filepath.Dir(file), 0777)
	if err != nil {
		api.PrintError(err.Error())
		return false
	}

	f, err := os.Create(file)
	if err != nil {
		api.PrintError(err.Error())
		return false
	}
	defer f.Close()

	_, err = f.Write(j)

	if err != nil {
		api.PrintError(err.Error())
		return false
	}

	return true
}
