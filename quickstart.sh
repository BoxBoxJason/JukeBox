#!/bin/sh
JUKEBOX_SCRIPT_DIR=$(dirname "$0")
JUKEBOX_FRONTEND_DIR=${JUKEBOX_SCRIPT_DIR}/frontend
JUKEBOX_BIN_DIR=${JUKEBOX_SCRIPT_DIR}/bin
JUKEBOX_BIN=${JUKEBOX_BIN_DIR}/jukebox

cd ${JUKEBOX_FRONTEND_DIR} && \
npm install && \
npm run build && \
cd .. && \
go mod tidy && \
mkdir -p ${JUKEBOX_BIN_DIR} && \
go build -o ${JUKEBOX_BIN} ${JUKEBOX_SCRIPT_DIR}/cmd/server/main.go && \
chmod +x ${JUKEBOX_BIN} && \
${JUKEBOX_BIN}
