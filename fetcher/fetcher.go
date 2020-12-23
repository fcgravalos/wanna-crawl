package fetcher

import (
	"context"
	"net/http"
	"time"

	logr "github.com/sirupsen/logrus"
)

// Fetcher interface just aims to make other packages easier to test. I don't expect, having multiple implementations
type Fetcher interface {
	Fetch(u string) ([]byte, error)
}

// NewHTTPFetcher returns a Fetcher given a `ctx` context and `reqTimeout` timeout parameters
func NewHTTPFetcher(ctx context.Context, logger *logr.Logger, reqTimeout time.Duration) Fetcher {
	return &httpFetcher{
		ctx,
		logger,
		&http.Client{Timeout: reqTimeout},
	}
}
