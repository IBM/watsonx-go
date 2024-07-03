# watsonx-go

`watsonx-go` is a [watsonx](https://www.ibm.com/watsonx) Client for Go

## Install

```sh
go get -u github.com/IBM/watsonx-go
```

## Usage

```go
import (
  wx "github.com/IBM/watsonx-go/pkg/models"
)
```

### Example Usage

```sh
export WATSONX_API_KEY="YOUR WATSONX API KEY"
export WATSONX_PROJECT_ID="YOUR WATSONX PROJECT ID"
```

Create a client:

```go
client, _ := wx.NewClient()
```

Or pass in the required secrets directly:

```go
client, err := wx.NewClient(
  wx.WithWatsonxAPIKey(apiKey),
  wx.WithWatsonxProjectID(projectID),
)
```

Generation:

```go
result, _ := client.GenerateText(
  "meta-llama/llama-3-70b-instruct",
  "Hi, who are you?",
  wx.WithTemperature(0.9),
  wx.WithTopP(.5),
  wx.WithTopK(10),
  wx.WithMaxNewTokens(512),
)

println(result.Text)
```

### Customization
If you want to use Watsonx test environment, choose one of the following methods:

#### Option 1: Using Environment Variables

Specify the Watsonx URL and IAM endpoint using environment variables:
```sh
export WATSONX_URL_HOST="us-south.ml.test.cloud.ibm.com"
export WATSONX_IAM_HOST="iam.test.cloud.ibm.com"
```

#### Option 2: Using the `NewClient` Function Parameters

Specify the Watsonx URL and IAM endpoint through the parameters of the NewClient function:
```go
client, err := wx.NewClient(
  wx.WithURL("us-south.ml.test.cloud.ibm.com"),
  wx.WithIAM("iam.test.cloud.ibm.com"),
  wx.WithWatsonxAPIKey(apiKey),
  wx.WithWatsonxProjectID(projectID),
)
```

## Development Setup

### Tests

#### Setup

```sh
export WATSONX_API_KEY="YOUR WATSONX API KEY"
export WATSONX_PROJECT_ID="YOUR WATSONX PROJECT ID"
```

#### Run

```sh
go test ./...
```

### Pre-commit Hooks

Run the following command to run pre-commit formatting:

```sh
git config --local core.hooksPath .githooks/
```

## Resources

- [watsonx REST API Docs](https://cloud.ibm.com/apidocs/watsonx-ai)
- [watsonx Python SDK Docs](https://ibm.github.io/watson-machine-learning-sdk)
