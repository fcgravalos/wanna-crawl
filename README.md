# Wanna Crawl

Wanna Crawl is a simple and concurrent web crawler written in Go. It is compound of the typical web crawler components:

- Frontiers: It acts as a kind of a "crawling director", deciding what needs to be crawl.
- Seen: In charge of keeping track of visited links, so we don't crawl the same url twice.
- Storage: To store and dump crawling results data.
- Crawler: It will fetch a url, given by the `Frontier`, parse it and extract it links.

The `Frontier` can scale to:

- Multiple servers, so seeds are processed concurrently. Use `-frontier.max-pool-size` to tune this.

- Multiple workers: long lived go routines that will wait for crawling jobs. Use `-frontier.max-concurrency` to tune this.

The following diagram illustrates this architecture:

```text

  +---------------+            +---------------+            +----------------+           +-------------+
  |               +<-----------+               <------------+                +---------->+             |
  |      SEEN     |            |  FRONTIERS    |            |    CRAWLER     |           |  INTERNET   |
  |               +----------->+               +----------->+                +<----------+             |
  +---------------+            +---------------+            +----------------+           +-------------+
                                       |
                                       |
                                       |
                               +-------v-------+
                               |               |
                               |   STORAGE     |
                               |               |
                               +---------------+
```

## How to build it

### Run unit tests

`make test`

### Build a single binary

`make build`

This will compile the source code and leave the binary at `./bin/wanna-crawl`

### Build Docker image

`make docker-build`

## How to run it

The only requirement to run Wanna Crawl is to provide it with a text file containing the list of seed urls.

Wanna Crawl does not parse this file, meaning that each line of the seed file should be a well formatted canonical url, for instance [https://www.example.com/](https://www.example.com)

The list of flags provided are listed in the table below:

| Flag | Go Type | Default | Description |
| ---- | ---- | ------- | ------ |
|`-crawler.follow-external-links`| `bool` | true | Whether or not to extract links outside the subdomain of the root url.|
|`-fetcher.request-timeout duration`| `time.Duration` | 3s | HTTP Request connection timeout.|
|`-frontier.max-concurrency` | `int` | 8 | Max number of workers attending to crawling jobs.|  
|`-frontier.max-depth`| `int`| 2 | The max number of links a single url can  be reached from|
|`-frontier.max-pool-size`| `int` | 4 | Max number of frontier servers that can be started concurrently |
|`-frontier.publish-queue-size` | `int` | 1024 | Size for the queue where workers will store results.|
|`-log.level` | `string` | "error" | Logging level: error, warning, info or debug. |
|`-seeds.file`| `string` | "seeds.txt" | Source for seed urls |
|`-seen_cache.engine` | `string` | "in-memory" | Seen cache engine to use to track already seen urls|
|`-storage.engine`| `string` | "in-memory" | Storage engine to use to ingest crawling results.|
|`-version`| `bool`| false | Print Wanna Crawl version |

Run `wanna-crawl [flags]` to override the defaults.

To run it as a docker container:

`docker run -v ${PATH_TO_SEEDS_FILE}:/seeds.txt --rm wanna-crawl:${WANNA_CRAWL_VERSION} [flags]`

## TODO

- Implement a robots.txt parser. I have leaned towards implementing `max-depth` first, to make sure the crawler terminates quickly. But a parser it's a must.

- Create implementations of `seen` and `storage` for a more scalable data structures. Some proposals:

  - `seen`: Redis, Memcached, Etcd, S3, etc.

  - `store`: DynamoDB, S3, etc.

- Use caches to store results, so we avoid crawling recently visited pages.

- Improve url filtering. Maybe all the logic of treating the HTML result should go to a different package.
