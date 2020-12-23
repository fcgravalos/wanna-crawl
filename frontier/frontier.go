package frontier

import (
	"context"
	"sync"

	"github.com/fcgravalos/wanna-crawl/crawler"
	"github.com/fcgravalos/wanna-crawl/seen"
	"github.com/fcgravalos/wanna-crawl/storage"
	logr "github.com/sirupsen/logrus"
)

// Config represents Frontier configuration
type Config struct {
	MaxPoolSize      int
	MaxConcurrency   int
	MaxDepth         int
	PublishQueueSize int
}

// Frontier will tell the crawler what to crawl next
type Frontier struct {
	// To smoothly shutdown workers
	ctx context.Context
	// To track already seen  urls
	seen.Cache
	// Where crawling results will be stored
	storage.Storage
	// The actual crawler
	*crawler.Crawler
	// Logger
	*logr.Logger
	// Frontier configuration
	Config
}

func (f *Frontier) spawnCrawlingWorkers(wg *sync.WaitGroup, next chan string, publish chan []string, limits chan struct{}) {
	for i := 0; i < f.MaxConcurrency; i++ {
		wg.Add(1)
		go func(workerId int, next chan string, publish chan []string, limits chan struct{}) {
			log := f.WithFields(logr.Fields{
				"frontier_role": "worker",
				"worker_id":     workerId,
			})

			for {
				select {
				case l := <-next:
					found, err := f.Crawl(l)
					if err != nil {
						log.Error(err, "failed to crawl", "url", l)
						break
					}
					f.Store(l, found)
					publish <- found
				case <-limits:
					log.Debug("max depth reached, shutting down")
					wg.Done()
					return
				case <-f.ctx.Done():
					log.Debug("context canceled shutting down")
					wg.Done()
					return
				}
			}
		}(i, next, publish, limits)
	}
}

func (f *Frontier) run(seeds []string) {
	log := f.WithFields(logr.Fields{
		"frontier_role": "manager",
	})
	next := make(chan string)
	publish := make(chan []string, f.PublishQueueSize)

	// Initialize the frontier
	log.Info("initializing frontier with seeds")
	go func(s []string, p chan []string) { p <- s }(seeds, publish)

	// To cancel smoothly all crawling jobs when depth limit has been reached
	limits := make(chan struct{}, 1)

	// Start workers
	log.Infof("starting %d workers", f.MaxConcurrency)
	var wg sync.WaitGroup
	go f.spawnCrawlingWorkers(&wg, next, publish, limits)

	// Dispatch crawling jobs; stop when limit hit
	for depth := f.MaxDepth; depth >= 0; depth-- {
		batch := <-publish
		for _, link := range batch {
			if !f.Seen(link) {
				if err := f.Add(link); err != nil {
					log.Warnf("failed to add %s to seen cache, might be revisited", link)
				}
				next <- link
			}
		}
	}
	// Max Depth hit
	close(limits)
	wg.Wait()

	// Cleanup channels
	close(next)
	close(publish)
	return
}

// StartManager will start all Frontier servers ans will wait for the result
func (f *Frontier) StartManager(seeds []string, done chan struct{}) {
	// Start frontier pool
	var wg sync.WaitGroup
	frontierPool := make(chan struct{}, f.MaxPoolSize)
	for _, seed := range seeds {
		frontierPool <- struct{}{}
		wg.Add(1)
		go func(f *Frontier, seed string, pool chan struct{}) {
			defer wg.Done()
			f.run([]string{seed})
			<-pool
		}(f, seed, frontierPool)
	}
	wg.Wait()
	done <- struct{}{}
}

// NewFrontier will return a Frontier object
func NewFrontier(ctx context.Context, seenCache seen.Cache, db storage.Storage, c *crawler.Crawler, l *logr.Logger, cfg Config) *Frontier {
	return &Frontier{
		ctx,
		seenCache,
		db,
		c,
		l,
		cfg,
	}
}
