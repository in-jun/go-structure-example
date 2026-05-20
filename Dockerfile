FROM golang:1.24-alpine AS builder

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=dev
ARG GIT_COMMIT=unknown
RUN BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ) && \
    CGO_ENABLED=0 go build -ldflags "-s -w -X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT" -o /api ./cmd/api

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /api /api
COPY --from=builder /app/migrations /migrations
EXPOSE 8080
ENTRYPOINT ["/api"]
