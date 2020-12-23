package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fcgravalos/wanna-crawl/crawler"
	"github.com/fcgravalos/wanna-crawl/fetcher"
	"github.com/fcgravalos/wanna-crawl/frontier"
	"github.com/fcgravalos/wanna-crawl/seen"
	"github.com/fcgravalos/wanna-crawl/storage"
	logr "github.com/sirupsen/logrus"
)

var version string

func main() {
	var frontierCfg frontier.Config
	var crawlerCfg crawler.Config
	var storageEngine string
	var seenCacheEngine string
	var seedFile string
	var logLevel string
	var printVersion bool
	var fetcherRequestTimeout time.Duration

	flag.BoolVar(&printVersion, "version", false, "Print wanna-crawl version")
	flag.DurationVar(&fetcherRequestTimeout, "fetcher.request-timeout", 3*time.Second, "HTTP Request connection timeout.")
	flag.BoolVar(&crawlerCfg.FollowExternalLinks, "crawler.follow-external-links", true, "Whether or not to extract links outside the subdomain of the root url.")
	flag.IntVar(&frontierCfg.MaxConcurrency, "frontier.max-concurrency", 8, "Max number of workers attending to crawling jobs.")
	flag.IntVar(&frontierCfg.MaxDepth, "frontier.max-depth", 2, "The max number of links a single url can  be reached from.")
	flag.IntVar(&frontierCfg.MaxPoolSize, "frontier.max-pool-size", 4, "Max number of frontier servers that can be started concurrently.")
	flag.IntVar(&frontierCfg.PublishQueueSize, "frontier.publish-queue-size", 1024, "Size for the queue where workers will store results.")
	flag.StringVar(&storageEngine, "storage.engine", "in-memory", "Storage engine to use to ingest crawling results.")
	flag.StringVar(&seenCacheEngine, "seen_cache.engine", "in-memory", "Seen cache engine to use to track already seen urls.")
	flag.StringVar(&seedFile, "seeds.file", "seeds.txt", "Source for seed urls")
	flag.StringVar(&logLevel, "log.level", "error", "Logging level: error, warning, info or debug")
	flag.Parse()

	// Print version
	if printVersion {
		fmt.Printf("Wanna Crawl %s\n", version)
		os.Exit(0)
	}
	// Read seeds from seedFile
	fd, err := os.Open(seedFile)
	if err != nil {
		fmt.Printf("Failed to read seed file %s: %v\n", seedFile, err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(fd)
	seeds := []string{}
	for scanner.Scan() {
		seeds = append(seeds, scanner.Text())
	}

	fd.Close()

	// Create wanna-crawl components objects
	ctx, cancel := context.WithCancel(context.Background())

	var log logr.Logger
	log.SetFormatter(&logr.JSONFormatter{})
	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)
	switch logLevel {
	case "error":
		log.SetLevel(logr.ErrorLevel)
	case "warning":
		log.SetLevel(logr.WarnLevel)
	case "info":
		log.SetLevel(logr.InfoLevel)
	case "debug":
		log.SetLevel(logr.DebugLevel)
	default:
		log.SetLevel(logr.ErrorLevel)
	}

	db, _ := storage.NewStorage(storageEngine)
	seenCache, _ := seen.NewCache(seenCacheEngine)

	c := crawler.NewCrawler(fetcher.NewHTTPFetcher(ctx, &log, fetcherRequestTimeout), &log, crawlerCfg)
	f := frontier.NewFrontier(ctx, seenCache, db, c, &log, frontierCfg)

	done := make(chan struct{}, 1)
	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)

	go f.StartManager(seeds, done)

	for {
		select {
		case <-done:
			break
		case <-sig:
			cancel()
			// canceling will make all frontiers to send to the done channel too
		}
		break
	}
	// Print sitemap
	sitemap, err := db.Dump()
	if err != nil {
		fmt.Printf("failed to print sitemap")
		os.Exit(1)
	}
	fmt.Println(sitemap)
}
