FROM golang:1.25.8-alpine AS build

WORKDIR /src/Server

COPY Server/go.mod Server/go.sum ./
RUN go mod download

COPY Server/ ./

ARG VERSION=dev
ARG TARGETARCH=amd64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build \
    -trimpath \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o /out/owncord-server .

FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata wget \
    && addgroup -S owncord \
    && adduser -S -D -H -h /app -s /sbin/nologin -G owncord owncord \
    && mkdir -p /app/data/uploads /app/data/backups /app/data/acme_certs \
    && chown -R owncord:owncord /app

WORKDIR /app

COPY --from=build /out/owncord-server /usr/local/bin/owncord-server

USER owncord:owncord

EXPOSE 8080
VOLUME ["/app/data"]

HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
  CMD wget -qO- http://127.0.0.1:8080/api/v1/health >/dev/null || exit 1

ENTRYPOINT ["owncord-server"]
