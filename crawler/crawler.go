package crawler

import (
	"bytes"
	neturl "net/url"

	"github.com/fcgravalos/wanna-crawl/fetcher"
	logr "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

// Config represents crawler configuration
type Config struct {
	// Whether or not to follow external links on a scraped page
	FollowExternalLinks bool
}

// Crawler holds the crawler data structure
type Crawler struct {
	fetcher.Fetcher
	*logr.Logger
	Config
}

func (c *Crawler) normalizeURL(url string, link string) (string, error) {
	u, err := neturl.Parse(link)
	if err != nil {
		return "", err
	}
	base, err := neturl.Parse(url)
	if err != nil {
		return "", err
	}

	return base.ResolveReference(u).String(), nil
}

func (c *Crawler) isInternal(urlA string, urlB string) bool {
	a, err := neturl.Parse(urlA)
	if err != nil {
		return false
	}

	b, err := neturl.Parse(urlB)
	if err != nil {
		return false
	}

	return a.Hostname() == b.Hostname()
}

func (c *Crawler) extractLinksFromPage(url string, page []byte) []string {
	links := []string{}
	extracted := map[string]bool{url: true} // Keep track of the already extracted links

	r := bytes.NewReader(page)
	it := html.NewTokenizer(r)

	for {
		token := it.Next()

		switch {
		case token == html.StartTagToken:
			token := it.Token()
			if token.Data == "a" {
				// found anchor tag, find href attr
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						l, err := c.normalizeURL(url, attr.Val)
						if err != nil {
							c.Warnf("malformed url %s", l)
							continue
						} else if extracted[l] {
							// If the same link is present in the page, ignore it
							continue
						} else if !c.FollowExternalLinks && !c.isInternal(url, l) {
							c.Debugf("discarding %s as it's an external link", l)
							continue
						}
						extracted[l] = true
						links = append(links, l)
					}
				}
			}
		case token == html.ErrorToken:
			return links
		}
	}
}

// Crawl receivers a string `url` and it will return the links ([]string) found
func (c *Crawler) Crawl(url string) ([]string, error) {
	page, err := c.Fetch(url)
	if err != nil {
		return nil, err
	}

	return c.extractLinksFromPage(url, page), nil
}

// NewCrawler builds a `Crawler` object
func NewCrawler(f fetcher.Fetcher, l *logr.Logger, cfg Config) *Crawler {
	return &Crawler{
		f,
		l,
		cfg,
	}
}
