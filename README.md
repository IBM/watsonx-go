# go-watsonx

A watsonx API Wrapper Client for Go

## Install

Install:

```sh
go get -u github.com/h0rv/go-watsonx
```

Import:

```go
import (
  wx "github.com/h0rv/go-watsonx/models"
)
```

## Example Usage

```go
	model, _ := wx.NewModel(
		yourIBMCloudAPIKey,
		yourWatsonxProjectID,
		wx.WithModel(wx.LLAMA_2_70B_CHAT),
	)

	result, _ := model.GenerateText(
    "Hi, who are you?",
		wx.WithTemperature(0.9),
		wx.WithTopP(.5),
		wx.WithTopK(10),
		wx.WithMaxNewTokens(512),
		wx.WithDecodingMethod(wx.Greedy),
	)

  println(result)
```

## Setup

### Pre-commit Hooks

Run the following command to run pre-commit formatting:

```sh
git config --local core.hooksPath .githooks/
```

## Resources

- [watsonx Python SDK Docs](https://ibm.github.io/watson-machine-learning-sdk)
- [watsonx REST API Docs (Internal)](https://test.cloud.ibm.com/apidocs/watsonx-ai)
