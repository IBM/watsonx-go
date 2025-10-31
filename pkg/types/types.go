package types

import (
	"net/http"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
	DoWithRetry(req *http.Request) (*http.Response, error)
}
