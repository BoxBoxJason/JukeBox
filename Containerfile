FROM alpine:latest AS build

WORKDIR /opt/jukebox

RUN apk add --no-cache npm go

COPY ./cmd/ /opt/jukebox/cmd/
COPY ./frontend/ /opt/jukebox/frontend/
COPY ./internal/ /opt/jukebox/internal/
COPY ./pkg/ /opt/jukebox/pkg/
COPY ./go.mod /opt/jukebox/

RUN npm --prefix /opt/jukebox/frontend install && \
    npm --prefix /opt/jukebox/frontend run build && \
    go mod tidy && \
    go build -o jukebox ./cmd/server && \
    chmod +x jukebox

FROM alpine:latest AS runtime

WORKDIR /opt/jukebox

COPY --from=build /opt/jukebox/jukebox /opt/jukebox/jukebox
COPY --from=build /opt/jukebox/frontend/dist /opt/jukebox/frontend/dist

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 CMD [ "curl", "-f", "localhost:3000/api/health" ]

EXPOSE 3000

CMD ["/opt/jukebox/jukebox"]
