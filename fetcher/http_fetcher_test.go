package fetcher

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	logr "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const fakeResponse = `
<!DOCTYPE html>
<html>
<body>

<h2>HTML Links</h2>
<p>HTML links are defined with the a tag:</p>

<a href="https://www.wanna-crawl.com/login">This is a link</a>
<a href="https://www.wanna-crawl.com/login">This is a duplicated link</a>
<a href="/about-us">This is a relative link</a>
<a href="/index.html">This is a link</a>
<a href="https://external.com/example">External link</a>
</body>
</html>`

var fakeURL string

func TestFetch(t *testing.T) {

	logger := new(logr.Logger)
	f := NewHTTPFetcher(context.Background(), logger, 1*time.Second)

	response, err := f.Fetch(fakeURL)

	assert.Nil(t, err)
	assert.Equal(t, []byte(fakeResponse), response)
}

func TestMain(m *testing.M) {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fakeResponse))
	})

	server := httptest.NewServer(r)
	fakeURL = server.URL
	defer server.Close()
	os.Exit(m.Run())
}
