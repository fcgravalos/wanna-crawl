FROM golang:1.13.1 as builder

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY wanna-crawl.go wanna-crawl.go
COPY crawler/ crawler/
COPY frontier/ frontier/
COPY seen/ seen/
COPY storage/ storage/
COPY fetcher/ fetcher/

ARG WANNA_CRAWL_VERSION

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=${WANNA_CRAWL_VERSION}" -a -o wanna-crawl wanna-crawl.go

FROM gcr.io/distroless/static:latest
WORKDIR /
COPY --from=builder /workspace/wanna-crawl .
ENTRYPOINT ["/wanna-crawl"]
