FROM golang:1.25 AS builder
WORKDIR /app

ARG VERSION

# Static build deps
RUN apt-get update && apt-get install -y musl-tools \
    && rm -rf /var/lib/apt/lists/*

# deps
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# source
COPY . .

# build
RUN CC=musl-gcc CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags "-s -w -linkmode=external -extldflags=-static -X main.Version=${VERSION}" \
    -o dorcs ./cmd/dorcs

# runtime
FROM scratch
WORKDIR /app
# Copy CA certificates for TLS verification (needed for GitHub API)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/dorcs /app/dorcs
EXPOSE 8080
ENV SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt
CMD ["/app/dorcs"]
