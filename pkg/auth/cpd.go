package auth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/IBM/watsonx-go/pkg/types"
)

const (
	CpdTokenEndpointPath = "/icp4d-api/v1/authorize"
)

type CPDAuthPairRequest struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	ApiKey   string `json:"api_key,omitempty"`
}

type CPDAuthResponse struct {
	Token     string `json:"token"`
	Exception string `json:"exception"`
}

type CPDAuthenticator struct {
	client   types.Doer
	cpdHost  string
	username string
	password string
	apiKey   string
}

func NewCPDAuthenticator(client types.Doer, cpdHost, username, password, apiKey string) (*CPDAuthenticator, error) {
	if password != "" && apiKey != "" {
		return nil, errors.New("password and apiKey should not both be passed")
	}
	return &CPDAuthenticator{
		client:   client,
		cpdHost:  cpdHost,
		username: username,
		password: password,
		apiKey:   apiKey,
	}, nil
}

func (a *CPDAuthenticator) GenerateToken() (*AuthToken, error) {

	request := &CPDAuthPairRequest{
		Username: a.username,
		Password: a.password,
		ApiKey:   a.apiKey,
	}

	payloadJSON, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	endpoint := url.URL{
		Scheme: "https",
		Host:   a.cpdHost,
		Path:   CpdTokenEndpointPath,
	}
	req, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tokenRes CPDAuthResponse
	err = json.Unmarshal(body, &tokenRes)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed authenticating with status code [%d]: %s", resp.StatusCode, tokenRes.Exception)
	}
	exp, err := extractExpIfAvailable(tokenRes.Token)
	if err != nil {
		return nil, err
	}
	return &AuthToken{
		tokenRes.Token,
		exp,
	}, nil
}

func extractExpIfAvailable(token string) (time.Time, error) {
	// instead of importing a whole jwt library for checking exp,
	// we simply extract it ourselves

	payloadEncoded := strings.Split(token, ".")[1]
	payloadDecoded, err := base64.RawURLEncoding.DecodeString(payloadEncoded)
	if err != nil {
		return time.Unix(0, 0), err
	}
	var exp struct {
		Exp float64 `json:"exp"`
	}
	if err := json.Unmarshal(payloadDecoded, &exp); err != nil {
		return time.Unix(0, 0), err
	}
	if exp.Exp == 0 {
		// If there is no expiration, just refresh it anyway every 12 hours
		return time.Now().Add(time.Hour * 12), nil
	}

	seconds := int64(exp.Exp)
	nanoseconds := int64((exp.Exp - float64(seconds)) * float64(time.Second))

	return time.Unix(seconds, nanoseconds), nil
}
