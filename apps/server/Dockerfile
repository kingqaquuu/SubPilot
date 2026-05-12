FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/subpilot-server ./cmd/server

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /bin/subpilot-server /usr/local/bin/subpilot-server
COPY configs/config.example.yaml ./configs/config.example.yaml
COPY docs ./docs
COPY scripts/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh

RUN chmod +x /usr/local/bin/docker-entrypoint.sh

EXPOSE 18080

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["subpilot-server"]
