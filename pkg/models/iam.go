package models

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	TokenPath string = "/identity/token"
)

type IAMToken struct {
	value      string
	expiration time.Time
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	Expiration  int64  `json:"expiration"`
}

func GenerateToken(client Doer, watsonxApiKey WatsonxAPIKey, iamCloudHost string) (IAMToken, error) {
	values := url.Values{
		"grant_type": {"urn:ibm:params:oauth:grant-type:apikey"},
		"apikey":     {watsonxApiKey},
	}

	payload := strings.NewReader(values.Encode())

	iamTokenEndpoint := url.URL{
		Scheme:   "https",
		Host:     iamCloudHost,
		Path:     TokenPath,
	}
	req, err := http.NewRequest(http.MethodPost, iamTokenEndpoint.String(), payload)
	if err != nil {
		return IAMToken{}, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return IAMToken{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return IAMToken{}, err
	}

	var tokenRes TokenResponse
	err = json.Unmarshal(body, &tokenRes)
	if err != nil {
		return IAMToken{}, err
	}

	return IAMToken{
		tokenRes.AccessToken,
		time.Unix(tokenRes.Expiration, 0),
	}, nil

}

func (t *IAMToken) Expired() bool {
	return t.expiration.Before(time.Now())
}
