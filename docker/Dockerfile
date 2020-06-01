# Build Stage
FROM golang:1.14.3-alpine3.11 AS builder
RUN apk add git
WORKDIR /go/src/lazy-rest
ADD ./go.* ./
RUN go mod download
ADD ./ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server

# Deploy Stage
FROM scratch
ARG VERSION
ARG VCS_REF
LABEL maintainer="arnoud@kleinloog.ch" \
    description="Lazy REST in Go" \
    version=${VERSION} \
    vcs_url="https://github.com/akleinloog/lazy-rest" \
    vcs-ref=${VCS_REF}
COPY --from=builder /go/src/lazy-rest/server .
EXPOSE 8080
CMD ["./server"]