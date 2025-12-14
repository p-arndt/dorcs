FROM golang:1.25 AS builder
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o dorcs ./cmd/dorcs

FROM scratch
COPY --from=builder /app/dorcs /app/dorcs
WORKDIR /app
EXPOSE 8080
CMD ["/app/dorcs"]