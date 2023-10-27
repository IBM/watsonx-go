package utils

import (
	"errors"
)

const IAM_URL = "https://iam.cloud.ibm.com"
const TOKENS_URL = IAM_URL + "/identity/token"

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func GetIAMToken(apiKey string) (string, error) {
	return "Not implemented", errors.New("Not implemented")
}
