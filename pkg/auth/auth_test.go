package auth

import (
	"os"
	"testing"

	"github.com/IBM/watsonx-go/pkg/constants"
	"github.com/IBM/watsonx-go/pkg/http"
)

func TestCPDAuthenticator(t *testing.T) {
	cpdHost := os.Getenv(constants.CPDHostEnvVarName)
	apiKey := os.Getenv(constants.CPDAPIKeyEnvVarName)
	cpdUsername := os.Getenv(constants.CPDUsernameEnvVarName)
	cpdPassword := os.Getenv(constants.CPDPasswordEnvVarName)

	if cpdHost == "" {
		t.Skip("CPD environment was not specified. Skipping...")
	}
	if cpdUsername == "" {
		t.Fatal("CPD username was not set")
	}
	if cpdPassword == "" && apiKey == "" {
		t.Fatal("CPD Auth missing either Password or APIKey")
	}
	if cpdPassword != "" && apiKey != "" {
		// default to password
		apiKey = ""
		t.Log("Both Password and APIKey was set; defaulting to Password")
	}

	authenticator, err := NewCPDAuthenticator(http.NewHttpClient(), cpdHost, cpdUsername, cpdPassword, apiKey)
	if err != nil {
		t.Fatal(err)
	}

	token, err := authenticator.GenerateToken()
	if err != nil {
		t.Fatal(err)
	}
	if token.Value == "" {
		t.Fatal("token was empty though generate token passed")
	}
}

func TestIAMAuthenticator(t *testing.T) {
	apiKey := os.Getenv(constants.WatsonxAPIKeyEnvVarName)
	iam := os.Getenv(constants.WatsonxIAMEnvVarName)

	if iam == "" {
		iam = constants.DefaultIAMCloudHost
	}
	if apiKey == "" {
		t.Fatal("Unable to test IAM auth; please set required environment variables")
	}

	authenticator, err := NewIAMAuthenticator(http.NewHttpClient(), apiKey, iam)
	if err != nil {
		t.Fatal(err)
	}

	token, err := authenticator.GenerateToken()
	if err != nil {
		t.Fatal(err)
	}

	if token.Value == "" {
		t.Fatal("token was empty though generate token passed")
	}
}

func TestExtractExpIfAvailable(t *testing.T) {
	token := "eyJ0eXAiOiJhdCtqd3QiLCJhbGciOiJFUzI1NiIsImtpZCI6IjIwMjAzMjNlOTI2YTMyN2QwOWJiNDg4MDU1NGE2NDMwIn0.eyJzdWIiOiI1YmU4NjM1OTA3M2M0MzRiYWQyZGEzOTMyMjIyZGFiZSIsImV4cCI6MTc2MTkyOTU5NiwiaWF0IjoxNzYxOTI1OTk2fQ.GI6kTr3z_Rov92ssKFEHy9B2W5MLVqNRR7i1Gj7qyc_pV6jNyouQ9rXhfCd4ZWOb9WgQRVV5i42_LCRQWfJJdw"

	exp, err := extractExpIfAvailable(token)
	if err != nil {
		t.Fatal(err)
	}
	if exp.Unix() != 1761929596 {
		t.Fatal("the expiration time was not the expected time")
	}

}
