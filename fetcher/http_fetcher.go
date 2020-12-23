package fetcher

import (
	"context"
	"io/ioutil"
	"net/http"

	logr "github.com/sirupsen/logrus"
)

type httpFetcher struct {
	ctx context.Context
	*logr.Logger
	*http.Client
}

func (f *httpFetcher) Fetch(url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(f.ctx, "GET", url, nil)
	if err != nil {
		f.Debugf("failed to build HTTP GET request: %v", err)
		return nil, err
	}

	resp, err := f.Do(req)
	if err != nil {
		f.Errorf("failed HTTP GET request: %v", err)
		return nil, err
	}
	body := resp.Body
	defer body.Close()

	return ioutil.ReadAll(body)
}
