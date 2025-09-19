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
  wx.WithClientRetryConfig(wx.NewRetryConfig(
    wx.WithReturnHTTPStatusAsErr(false)),
  ),
  wx.WithWatsonxAPIKey(apiKey),
  wx.WithWatsonxProjectID(projectID),
)
```

#### Generate Text

Generation:

```go
result, _ := client.GenerateText(
  "meta-llama/llama-3-3-70b-instruct",
  "Hi, who are you?",
  wx.WithTemperature(0.4),
  wx.WithMaxNewTokens(512),
)

println(result.Text)
```

Stream Generation:

```go
dataChan, _ := client.GenerateTextStream(
  "meta-llama/llama-3-3-70b-instruct",
  "Hi, who are you?",
  wx.WithTemperature(0.4),
  wx.WithMaxNewTokens(512),
)

for data := range dataChan {
  print(data.Text) // print the result as it's being generated
}
```

#### Generate Embeddings

Embedding | Single query:

```go
result, _ := client.EmbedQuery(
	"ibm/slate-30m-english-rtrvr",
	"Hello, world!",
	wx.WithEmbeddingTruncateInputTokens(2), 
	wx.WithEmbeddingReturnOptions(true),
)

embeddingVector := result.Results[0].Embedding
```

Embedding | Multiple docs:

```go
result, _ := clientl.EmbedDocuments(
    "ibm/slate-30m-english-rtrvr",
    []string{"Hello, world!", "Goodbye, world!"},
    wx.WithEmbeddingTruncateInputTokens(2), 
    wx.WithEmbeddingReturnOptions(true),
)

for _, doc := range result.Results {
    fmt.Println(doc.Embedding)
}
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

### Using Test Environment

There are two methods for configuring the watsonx client to be used with the test environment:

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
  	wx.WithClientRetryConfig(wx.NewRetryConfig(
			wx.WithReturnHTTPStatusAsErr(false)),
		),
  wx.WithURL("us-south.ml.test.cloud.ibm.com"),
  wx.WithIAM("iam.test.cloud.ibm.com"),
  wx.WithWatsonxAPIKey(apiKey),
  wx.WithWatsonxProjectID(projectID),
)
```

---

## Resources

- [watsonx REST API Docs](https://cloud.ibm.com/apidocs/watsonx-ai)
- [watsonx Python SDK Docs](https://ibm.github.io/watson-machine-learning-sdk)
