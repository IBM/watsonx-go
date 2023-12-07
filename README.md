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
  wx "github.com/h0rv/go-watsonx/foundation_models"
)
```

## Example Usage

### Builder Pattern

```go
	model, err := wx.NewModelBuilder().
		SetModelId(wx.ModelTypesEnum.LLAMA_2_70B_CHAT).
		SetApiKey(yourWatsonxApiKey).
		SetProjectId(yourWatsonxProjectID).
		SetTemperature(yourtemperature).
		SetMaxNewTokens(yourMaxNewTokens).
		SetDecodingMethod(wx.DecodingMethodsEnum.GREEDY).
		Build()
	if err != nil {
		// Failed to get watsonx model
		return err
	}

	result, err := model.GenerateText(
		"Hi, how are you?",
		nil, /* or your Generation Params */
	)
	if err != nil {
		// Failed to call generate on model
		return err
	}
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
