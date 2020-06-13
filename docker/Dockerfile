# Build Stage
FROM golang:1.14.3-alpine3.11 AS builder
RUN apk add git
WORKDIR /go/src/lazy-rest
ADD ./go.* ./
RUN go mod download
ADD ./ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o lazy-rest

# Deploy Stage
FROM scratch
ARG VERSION
ARG VCS_REF
LABEL maintainer="arnoud@kleinloog.ch" \
    description="Lazy REST Server in Go" \
    version=${VERSION} \
    vcs_url="https://github.com/akleinloog/lazy-rest" \
    vcs-ref=${VCS_REF}
COPY --from=builder /go/src/lazy-rest/lazy-rest .
VOLUME ["/data"]
EXPOSE 8080
CMD ["./lazy-rest", "serve"]