#!/bin/sh
JUKEBOX_SCRIPT_DIR=$(dirname "$0")
JUKEBOX_FRONTEND_DIR=${JUKEBOX_SCRIPT_DIR}/frontend

cd ${JUKEBOX_FRONTEND_DIR}
npm install
npm run build

cd ${JUKEBOX_SCRIPT_DIR}
go mod tidy
go run ./cmd/server/main.go
