# AWS SDK CLI

`aws-sdk-cli` is a command-line interface (CLI) tool designed to interact with the AWS SDK written in Go. It provides an interactive shell, allowing users to authenticate with AWS credentials and retrieve specified resources.

## Features

- Authenticate with AWS shared credential files
- Describe specific resources
- Provide Interactive shell for ease of use
- Store results to JSON files for future use

## Usage

Make sure you have available AWS shared credentials on your machine.

Run the program by executing:

```
go run main.go
```

### Commands

#### 1. `run`

Perform the specific `action` from the `service` and dump the result to `output_file`.

Accepts an optional `params_file` in JSON that adds required and optional parameters to the action.

If you are not sure about required parameters, you may perform the action without parameters and refer to the error information.

Note that the auto-complete only lists common services and actions, but you may specify any of the available actions.

By default, the cli supports `ec2` and `iam` services.
You can add supports for any services by appending necessary information to function `fetchClient` in `clients.go`.

Refer to the official document for all available services, actions, and parameters.

```
run <service> <action> <output_file> [params_file]
```

## Examples

See `test` folder for more sample usages.

If you want to perform the test on your own, make sure to use valid parameter files.

**Run**:

`run ec2 DescribeInstances instances.json params.json`

*instances.json*

```json
[
  {
    "Groups": [],
    "Instances": [
      {
        "AmiLaunchIndex": 0,
        "Architecture": "x86_64"
      }
    ]
  },
  {
    "Groups": [],
    "Instances": [
      {
        "AmiLaunchIndex": 0,
        "Architecture": "x86_64"
      }
    ]
  }
]
```

## References

1. **AWS SDK v2**:
    - [Official Documentation](https://aws.github.io/aws-sdk-go-v2/docs)
    - [GitHub Repository](https://github.com/aws/aws-sdk-go-v2)
    - [API Reference](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2)
2. **Interactive Shell Library (ishell)**:
    - [GitHub Repository](https://github.com/abiosoft/ishell)
