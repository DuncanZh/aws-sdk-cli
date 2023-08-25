package api

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"reflect"
)

func (a *AWSAPI) fetchClient(service string) reflect.Value {
	var client reflect.Value
	switch service {
	case "ec2":
		client = reflect.ValueOf(ec2.NewFromConfig(*a.Config))
	case "eks":
		client = reflect.ValueOf(eks.NewFromConfig(*a.Config))
	case "iam":
		client = reflect.ValueOf(iam.NewFromConfig(*a.Config))
	}
	return client
}
