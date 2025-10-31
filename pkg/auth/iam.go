package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/IBM/watsonx-go/pkg/constants"
	"github.com/IBM/watsonx-go/pkg/types"
)

const (
	IAMTokenPath string = "/identity/token"
)

type IAMAuthenticator struct {
	client        types.Doer
	watsonxApiKey constants.WatsonxAPIKey
	iamCloudHost  string
}

type IAMTokenResponse struct {
	AccessToken  string `json:"access_token"`
	Expiration   int64  `json:"expiration"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

func NewIAMAuthenticator(client types.Doer, watsonxApiKey constants.WatsonxAPIKey, iamCloudHost string) (*IAMAuthenticator, error) {
	return &IAMAuthenticator{
		client:        client,
		watsonxApiKey: watsonxApiKey,
		iamCloudHost:  iamCloudHost,
	}, nil
}

func (a *IAMAuthenticator) GenerateToken() (*AuthToken, error) {
	values := url.Values{
		"grant_type": {"urn:ibm:params:oauth:grant-type:apikey"},
		"apikey":     {a.watsonxApiKey},
	}

	payload := strings.NewReader(values.Encode())

	iamTokenEndpoint := url.URL{
		Scheme: "https",
		Host:   a.iamCloudHost,
		Path:   IAMTokenPath,
	}
	req, err := http.NewRequest(http.MethodPost, iamTokenEndpoint.String(), payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tokenRes IAMTokenResponse
	err = json.Unmarshal(body, &tokenRes)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed authenticating with status code [%d]: %s", resp.StatusCode, tokenRes.ErrorMessage)
	}

	return &AuthToken{
		tokenRes.AccessToken,
		time.Unix(tokenRes.Expiration, 0),
	}, nil
}
