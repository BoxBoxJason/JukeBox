FROM alpine:latest AS build

WORKDIR /opt/jukebox

COPY ./cmd/ /opt/jukebox/cmd/
COPY ./frontend/ /opt/jukebox/frontend/
COPY ./internal/ /opt/jukebox/internal/
COPY ./pkg/ /opt/jukebox/pkg/
COPY ./go.mod /opt/jukebox/

RUN apk add --no-cache npm go && \
    npm --prefix /opt/jukebox/frontend install

RUN npm --prefix /opt/jukebox/frontend run build

RUN go mod tidy && \
    go build -o jukebox ./cmd/server && \
    chmod +x jukebox

FROM alpine:latest AS runtime

WORKDIR /opt/jukebox

COPY --from=build /opt/jukebox/jukebox /opt/jukebox/jukebox
COPY --from=build /opt/jukebox/frontend/dist /opt/jukebox/frontend/dist

CMD ["/opt/jukebox/jukebox"]

# HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 CMD [ "curl", "-f", "http://localhost:3000/api/health" ]
