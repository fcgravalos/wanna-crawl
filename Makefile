WANNA_CRAWL_VERSION=v0.0.1
BUILD_FLAGS=-ldflags "-X main.version=${WANNA_CRAWL_VERSION}"
DOCKER_IMG_TAG=wanna-crawl:${WANNA_CRAWL_VERSION}

fmt:
	go fmt ./...

vet:
	go vet ./...

test: fmt vet 
	go test -cover -v ./fetcher/... ./frontier/... ./crawler/... ./seen/... ./storage/... -coverprofile cover.out 

build: fmt vet
	go build ${BUILD_FLAGS} -o bin/wanna-crawl wanna-crawl.go

run: fmt vet
	go run ./wanna-crawl.go

docker-build: test
	docker build --build-arg WANNA_CRAWL_VERSION=${WANNA_CRAWL_VERSION} -t ${DOCKER_IMG_TAG} .