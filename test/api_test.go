package test

import (
	"AWS_API/api"
	"fmt"
	"os"
	"testing"
)

func TestEC2(t *testing.T) {
	service := "ec2"

	ec2List := []string{"DescribeInstances", "DescribeReservedInstances", "DescribeVolumes",
		"DescribeVolumeStatus", "DescribeSubnets", "DescribeSecurityGroups", "DescribeSecurityGroupRules",
		"DescribeNatGateways", "DescribeInternetGateways"}

	run(service, ec2List)
}

func TestEKS(t *testing.T) {
	service := "eks"

	eksList := []string{"DescribeCluster", "DescribeNodegroup", "ListClusters", "ListNodegroups"}

	run(service, eksList)
}

func TestIAM(t *testing.T) {
	service := "iam"

	iamList := []string{"GetRole", "GetUser", "GetUserPolicy", "ListAccessKeys", "ListAttachedRolePolicies",
		"ListGroups", "ListPolicies", "ListRoles", "ListUsers", "ListUserTags"}

	run(service, iamList)
}

func run(service string, list []string) {
	a := api.NewAPI()

	for _, action := range list {
		var params []byte

		paramsFile := fmt.Sprintf("params/%v/%v.json", service, action)
		_, err := os.Stat(paramsFile)
		if err == nil {
			params, err = os.ReadFile(paramsFile)
			if err != nil {
				api.PrintError(err)
				return
			}
		}

		result := a.Operate(service, action, params)
		api.DumpFile(result, fmt.Sprintf("%v/%v.json", service, action), true)
	}
}
