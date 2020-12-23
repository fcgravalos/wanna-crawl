package crawler

import (
	"testing"

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
	return []byte(fakeResponse), nil
}

func TestNormalizeURL(t *testing.T) {
	c := &Crawler{}

	testCases := []struct {
		urls     []string
		expected string
		err      bool
	}{
		{[]string{"https://wanna-crawl.com", "/about-us"}, "https://wanna-crawl.com/about-us", false},
		{[]string{"https://wanna-crawl.com/", "/about-us"}, "https://wanna-crawl.com/about-us", false},
		{[]string{"https://wanna-crawl.com/about-us", "/login"}, "https://wanna-crawl.com/login", false},
		{[]string{"https://wanna-crawl.com/contact", "index.html"}, "https://wanna-crawl.com/index.html", false},
		{[]string{"https://wanna-crawl.com/1/2/3/", ".index.html"}, "https://wanna-crawl.com/1/2/3/.index.html", false},
		{[]string{"https://wanna-crawl.com/1/2/3/", "../index.html"}, "https://wanna-crawl.com/1/2/index.html", false},
	}
	for _, tc := range testCases {
		failed := false
		normalized, err := c.normalizeURL(tc.urls[0], tc.urls[1])
		if err != nil {
			failed = true
		}
		assert.Equal(t, tc.expected, normalized)
		assert.True(t, tc.err == failed)
	}
}

func TestIsInternal(t *testing.T) {
	c := &Crawler{}
	urlA := "https://wanna-crawl.com/"
	urlB := "https://community.wanna-crawl.com/awesome"
	urlC := "https://wanna-crawl.com/login"

	assert.False(t, c.isInternal(urlA, urlB))
	assert.True(t, c.isInternal(urlA, urlC))
}

func TestExtractLinksFromPage(t *testing.T) {
	sampleHTML := []byte(fakeResponse)
	testCases := []struct {
		crawler       *Crawler
		url           string
		page          []byte
		expectedLinks []string
	}{
		{&Crawler{nil, new(logr.Logger), Config{FollowExternalLinks: false}}, "https://wanna-crawl.com/", sampleHTML, []string{"https://wanna-crawl.com/login", "https://wanna-crawl.com/about-us", "https://wanna-crawl.com/index.html"}},
		{&Crawler{nil, new(logr.Logger), Config{FollowExternalLinks: true}}, "https://wanna-crawl.com/", sampleHTML, []string{"https://wanna-crawl.com/login", "https://wanna-crawl.com/about-us", "https://wanna-crawl.com/index.html", "https://external.com/example"}},
	}

	for _, tc := range testCases {
		c := tc.crawler
		u := tc.url
		p := tc.page
		found := c.extractLinksFromPage(u, p)
		assert.Equal(t, tc.expectedLinks, found)
	}
}

func TestCrawl(t *testing.T) {
	cfg := Config{
		FollowExternalLinks: true,
	}

	c := NewCrawler(&testFetcher{}, new(logr.Logger), cfg)
	found, err := c.Crawl("https://wanna-crawl.com/")

	assert.Nil(t, err)
	assert.Equal(t,
		[]string{
			"https://wanna-crawl.com/login",
			"https://wanna-crawl.com/about-us",
			"https://wanna-crawl.com/index.html",
			"https://external.com/example",
		}, found)
}
