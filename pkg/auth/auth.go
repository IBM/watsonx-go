package auth

import (
	"time"
)

type Authenticator interface {
	GenerateToken() (*AuthToken, error)
}

type AuthToken struct {
	Value      string
	Expiration time.Time
}

func (t *AuthToken) Expired() bool {
	return t.Expiration.Before(time.Now())
}
