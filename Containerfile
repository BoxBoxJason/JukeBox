FROM alpine:latest AS build

WORKDIR /opt/jukebox

COPY ./cmd/ /opt/jukebox/cmd/
COPY ./frontend/ /opt/jukebox/frontend/
COPY ./internal/ /opt/jukebox/internal/
COPY ./pkg/ /opt/jukebox/pkg/
COPY ./go.mod /opt/jukebox/

RUN apk add --no-cache npm go && \
    cd frontend && \
    npm install && \
    npm run build && \
    cd .. && \
    go mod tidy && \
    go build -o jukebox ./cmd/server \
    && chmod +x jukebox

FROM alpine:latest AS runtime

WORKDIR /opt/jukebox

COPY --from=build /opt/jukebox/jukebox /opt/jukebox/jukebox
COPY --from=build /opt/jukebox/frontend/dist /opt/jukebox/frontend/dist

CMD ["/opt/jukebox/jukebox"]
