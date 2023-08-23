package api

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"reflect"
	"strings"
)

type AwsAPI struct {
	Config    *aws.Config
	EC2Client *ec2.Client
}

func NewAPI() *AwsAPI {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		PrintError(err.Error())
		return nil
	}

	api := &AwsAPI{}
	api.Config = &cfg
	api.EC2Client = ec2.NewFromConfig(*api.Config)

	return api
}

func (a *AwsAPI) Describe(operation string, ids []string) interface{} {
	client := reflect.ValueOf(a.EC2Client)

	/*
		switch module {
		case "ec2":
			client = reflect.ValueOf(a.EC2Client)
		default:
			PrintError("Unknown module")
			return nil
		}
	*/

	operation = strings.Title(operation)

	method := client.MethodByName("Describe" + operation)
	if !method.IsValid() {
		PrintError("Invalid operation")
		return nil
	}

	paramType := method.Type().In(1)
	paramPtr := reflect.New(paramType.Elem())
	param := paramPtr.Elem()

	inIdsField, _ := param.Type().FieldByNameFunc(func(s string) bool {
		return strings.HasSuffix(s, "Ids")
	})

	inTokenField, hasNextToken := param.Type().FieldByName("NextToken")

	outResultField, _ := method.Type().Out(0).Elem().FieldByNameFunc(func(s string) bool {
		field, ok := method.Type().Out(0).Elem().FieldByName(s)
		return ok && field.Type.Kind().String() == "slice"
	})

	outTokenField, _ := method.Type().Out(0).Elem().FieldByName("NextToken")

	var result reflect.Value

	for nextToken := new(string); nextToken != nil; {
		if *nextToken == "" {
			param.FieldByIndex(inIdsField.Index).Set(reflect.ValueOf(ids))
		} else {
			param.FieldByIndex(inIdsField.Index).Set(reflect.Zero(reflect.TypeOf([]string{})))
			param.FieldByIndex(inTokenField.Index).Set(reflect.ValueOf(nextToken))
		}

		ret := method.Call([]reflect.Value{reflect.ValueOf(context.Background()), paramPtr})

		if err := ret[1].Interface(); err != nil {
			PrintError(err.(error).Error())
			return nil
		}

		e := ret[0].Elem()

		if hasNextToken {
			if *nextToken == "" {
				result = e.FieldByIndex(outResultField.Index)
			} else {
				reflect.AppendSlice(result, e.FieldByIndex(outResultField.Index))
			}
			nextToken = e.FieldByIndex(outTokenField.Index).Interface().(*string)
		} else {
			return e.FieldByIndex(outResultField.Index).Interface()
		}
	}
	return result.Interface()
}

func PrintError(msg string) {
	println("Error: " + msg)
}
