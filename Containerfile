FROM alpine:latest

WORKDIR /opt/jukebox

RUN apk add --no-cache npm go && \
    rm -rf /var/cache/apk/*

EXPOSE 3000

CMD ["tail","-f","/dev/null"]
