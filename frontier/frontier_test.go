package frontier

import (
	"context"
	"encoding/json"
	neturl "net/url"
	"testing"

	"github.com/fcgravalos/wanna-crawl/crawler"
	"github.com/fcgravalos/wanna-crawl/seen"
	"github.com/fcgravalos/wanna-crawl/storage"
	logr "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const fakeResponse = `
<!DOCTYPE html>
<html>
<body>

<h2>HTML Links</h2>
<p>HTML links are defined with the a tag:</p>

<a href="https://wanna-crawl.com/login">This is a link</a>
<a href="https://wanna-crawl.com/login">This is a duplicated link</a>
<a href="/about-us">This is a relative link</a>
<a href="/index.html">This is a link</a>
<a href="https://external.com/example">External link</a>
</body>
</html>`

type testFetcher struct{}

func (t *testFetcher) Fetch(url string) ([]byte, error) {
	u, _ := neturl.Parse(url)
	resp := []byte("")
	if u.Path == "" || u.Path == "/" {
		resp = []byte(fakeResponse)
	}
	return resp, nil
}

func TestStartManager(t *testing.T) {
	cfg := Config{
		MaxPoolSize:      1,
		MaxConcurrency:   1,
		MaxDepth:         1,
		PublishQueueSize: 1024,
	}

	crawlerCfg := crawler.Config{
		FollowExternalLinks: true,
	}

	db, _ := storage.NewStorage("in-memory")
	seenCache, _ := seen.NewCache("in-memory")

	logger := new(logr.Logger)
	ctx := context.TODO()

	c := crawler.NewCrawler(&testFetcher{}, logger, crawlerCfg)
	f := NewFrontier(ctx, seenCache, db, c, logger, cfg)

	done := make(chan struct{}, 1)
	expectedSiteMap := map[string][]string{
		"https://wanna-crawl.com/":           []string{"https://wanna-crawl.com/login", "https://wanna-crawl.com/about-us", "https://wanna-crawl.com/index.html", "https://external.com/example"},
		"https://wanna-crawl.com/login":      []string{},
		"https://wanna-crawl.com/about-us":   []string{},
		"https://wanna-crawl.com/index.html": []string{},
		"https://external.com/example":       []string{},
	}
	expectedJSON, _ := json.MarshalIndent(expectedSiteMap, "", "\t")

	f.StartManager([]string{"https://wanna-crawl.com/"}, done)

	<-done
	sitemap, _ := db.Dump()
	assert.EqualValues(t, string(expectedJSON), sitemap)
}
