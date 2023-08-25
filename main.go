package main

import (
	"AWS_API/api"
	"fmt"
	"github.com/abiosoft/ishell/v2"
	"github.com/abiosoft/readline"
	"os"
	"reflect"
	"time"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	VERSION := "0.1.0"
	PROMPT := "aws_sdk_cli v" + VERSION + " $ "

	shell := ishell.NewWithConfig(&readline.Config{Prompt: PROMPT})
	shell.SetHomeHistoryPath(".aws_sdk_cli_shell_history")

	shell.Set("api", api.NewAPI())

	serviceList := []string{"ec2", "eks", "iam"}

	ec2List := []string{"DescribeInstances", "DescribeReservedInstances", "DescribeVolumes",
		"DescribeVolumeStatus", "DescribeSubnets", "DescribeSecurityGroups", "DescribeSecurityGroupRules",
		"DescribeNatGateways", "DescribeInternetGateways"}

	eksList := []string{"DescribeCluster", "DescribeNodegroup", "ListClusters", "ListNodegroups"}

	iamList := []string{"GetRole", "GetUser", "GetUserPolicy", "ListAccessKeys", "ListAttachedRolePolicies",
		"ListGroups", "ListPolicies", "ListRoles", "ListUsers", "ListUserTags"}

	shell.Set("service", serviceList)
	shell.Set("ec2", ec2List)
	shell.Set("eks", eksList)
	shell.Set("iam", iamList)

	shell.AddCmd(&ishell.Cmd{
		Name: "run",
		Help: "Usage: run <service> <action> <output_file> [params_file]",
		CompleterWithPrefix: func(prefix string, args []string) []string {
			if len(args) == 0 {
				return shell.Get("service").([]string)
			} else if len(args) == 1 {
				return shell.Get(args[0]).([]string)
			} else if len(args) == 2 || len(args) == 3 {
				return api.GetDirectory(prefix)
			}
			return []string{}
		},
		Func: run,
	})

	if len(os.Args) > 1 {
		err := shell.Process(os.Args[1:]...)
		if err != nil {
			api.PrintError(err)
		}
	} else {
		shell.Println("Interactive Shell:")
		shell.Run()
		shell.Close()
	}
}

func run(c *ishell.Context) {
	start := time.Now()

	a := c.Get("api").(*api.AWSAPI)

	inputFile := ""
	if len(c.Args) == 4 {
		inputFile = c.Args[3]
	} else if len(c.Args) != 3 {
		fmt.Println(c.Cmd.Help)
		return
	}

	service := c.Args[0]
	action := c.Args[1]
	outputFile := c.Args[2]

	var params []byte
	var err error
	if inputFile != "" {
		params, err = os.ReadFile(inputFile)
		if err != nil {
			api.PrintError(err)
			return
		}
	}

	result := a.Operate(service, action, params)

	if result == nil {
		fmt.Printf("Failed: Unable to process the input in %.2f seconds\n", time.Since(start).Seconds())
		return
	}

	entries := 1
	if reflect.ValueOf(result).Kind().String() == "slice" {
		entries = reflect.ValueOf(result).Len()
	}

	if api.DumpFile(result, outputFile, true) {
		fmt.Printf("Success: Processed %v entries in %.2f seconds\n", entries, time.Since(start).Seconds())
	}
}
