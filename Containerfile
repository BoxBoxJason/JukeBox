FROM alpine:latest as base

WORKDIR /opt/jukebox

RUN apk update && \
    apk upgrade && \
    apk add --no-cache go git && \
    rm -rf /var/cache/apk/*

COPY ./go.mod ./go.mod
COPY ./cmd/ ./cmd/
COPY ./internal/ ./internal/
COPY ./pkg/ ./pkg/

RUN go mod tidy && \
    go build -o /opt/jukebox/bin/jukebox ./cmd/server/main.go

FROM alpine:latest as production

RUN addgroup -S jukebox && \
    adduser -S -G jukebox -D -u 2300 jukebox

USER jukebox
WORKDIR /opt/jukebox

COPY --from=base --chown=jukebox:jukebox --chmod=550 /opt/jukebox/bin/jukebox /opt/jukebox/bin/jukebox

WORKDIR /home/jukebox

EXPOSE 3000

CMD ["/opt/jukebox/bin/jukebox"]
