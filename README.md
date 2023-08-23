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

#### 1. `describe`

Describe the specific resource and dump to `output_file`.

Accepts arbitrary amount of resource ids to describe certain items, defaults to retrieve all items of the type if not specified.

Refer to the official document for more available resources.

```
describe <resource> <output_file> [ids...]
```

## Examples

See `samples` folder for more sample outputs.

**Describe**:

`describe instances instances.json`

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
      ...
    ]
  },
  {
    "Groups": [],
    "Instances": [
      {
        "AmiLaunchIndex": 0,
        "Architecture": "x86_64"
      }
      ...
    ]
  },
  ...
]
```

## References

1. **AWS SDK v2**:
    - [Official Documentation](https://aws.github.io/aws-sdk-go-v2/docs)
    - [GitHub Repository](https://github.com/aws/aws-sdk-go-v2)
2. **Interactive Shell Library (ishell)**:
    - [GitHub Repository](https://github.com/abiosoft/ishell)
