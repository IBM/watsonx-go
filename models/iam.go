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
	IAMEndpoint = "https://iam.cloud.ibm.com/identity/token"
)

type IAMToken struct {
	value      string
	expiration time.Time
}

func GenerateToken(client Doer, ibmCloudAPIKey IBMCloudAPIKey) (IAMToken, error) {
	type TokenResponse struct {
		accessToken string `json:"access_token"`
		expiration  int64  `json:"expiration"`
	}

	values := url.Values{
		"grant_type": {"urn:ibm:params:oauth:grant-type:apikey"},
		"apikey":     {ibmCloudAPIKey},
	}

	payload := strings.NewReader(values.Encode())

	req, err := http.NewRequest(http.MethodPost, IAMEndpoint, payload)
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

	var tokenResponse TokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return IAMToken{}, err
	}

	return IAMToken{
		tokenResponse.accessToken,
		time.Unix(tokenResponse.expiration, 0),
	}, nil

}

func (t *IAMToken) Expired() bool {
	return t.expiration.Before(time.Now())
}
