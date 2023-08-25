package api

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"reflect"
)

type AWSAPI struct {
	Config *aws.Config
}

func NewAPI() *AWSAPI {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		PrintError(err)
		return nil
	}

	api := &AWSAPI{}
	api.Config = &cfg

	return api
}

func (a *AWSAPI) Operate(service string, action string, input []byte) interface{} {
	client := a.fetchClient(service)
	if !client.IsValid() {
		PrintErrorString("Invalid service " + service)
		return nil
	}

	ctx := reflect.ValueOf(context.Background())

	tokenList := []string{"nextToken", "NextToken", "marker", "Marker", "PaginationToken", "NextPageMarker"}

	method := client.MethodByName(action)
	if !method.IsValid() {
		PrintErrorString("Invalid action " + action)
		return nil
	}

	paramsPtrType := method.Type().In(1)
	paramsType := paramsPtrType.Elem()
	paramsPtr := reflect.New(paramsType)

	if input != nil && len(input) > 0 {
		err := json.Unmarshal(input, paramsPtr.Interface())
		if err != nil {
			PrintError(err)
			return nil
		}
	}

	ret := method.Call([]reflect.Value{ctx, paramsPtr})

	if err := ret[1].Interface(); err != nil {
		PrintError(err.(error))
		return nil
	}

	e := ret[0].Elem()

	tokenFunc := func(s string) bool {
		for _, v := range tokenList {
			if s == v {
				return true
			}
		}
		return false
	}

	tokenField, hasToken := e.Type().FieldByNameFunc(tokenFunc)

	if hasToken {
		var result reflect.Value

		outputField, _ := e.Type().FieldByNameFunc(func(s string) bool {
			field, ok := e.Type().FieldByName(s)
			return ok && field.Type.Kind().String() == "slice"
		})

		result = e.FieldByIndex(outputField.Index)
		nextToken := e.FieldByIndex(tokenField.Index)

		for !nextToken.IsZero() {
			paramsPtr.Elem().FieldByNameFunc(tokenFunc).Set(nextToken)

			ret = method.Call([]reflect.Value{ctx, paramsPtr})

			if err := ret[1].Interface(); err != nil {
				PrintError(err.(error))
				return nil
			}

			e = ret[0].Elem()

			result = reflect.AppendSlice(result, e.FieldByIndex(outputField.Index))
			nextToken = e.FieldByIndex(tokenField.Index)
		}
		return result.Interface()
	} else {
		result := map[string]interface{}{}
		for i := 0; i < e.Type().NumField(); i++ {
			field := e.Type().Field(i)
			if field.Type.String() != "middleware.Metadata" && field.Type.String() != "document.NoSerde" {
				result[field.Name] = e.FieldByIndex(field.Index).Interface()
			}
		}
		return result
	}
}
