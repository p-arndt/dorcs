FROM golang:1.25 AS builder
WORKDIR /app

ARG VERSION

# Install musl for static linking
RUN apt-get update && apt-get install -y musl-tools && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CC=musl-gcc CGO_ENABLED=1 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w -linkmode=external -extldflags=-static -X dorcs-v2/cmd/dorcs.Version=${VERSION}" \
    -o dorcs ./cmd/dorcs

FROM scratch
COPY --from=builder /app/dorcs /app/dorcs
WORKDIR /app
EXPOSE 8080
CMD ["/app/dorcs"]