#!/bin/sh
JUKEBOX_SCRIPT_DIR=$(dirname "$0")
JUKEBOX_FRONTEND_DIR=${JUKEBOX_SCRIPT_DIR}/frontend
JUKEBOX_BIN_DIR=${JUKEBOX_SCRIPT_DIR}/bin
JUKEBOX_BIN=${JUKEBOX_BIN_DIR}/jukebox

# Check if either podman-compose or docker is installed
if command -v podman-compose &> /dev/null
then
  echo "Using podman-compose"
  DOCKER_COMPOSE="podman-compose"
elif command -v docker &> /dev/null
then
  echo "Using docker compose"
  DOCKER_COMPOSE="docker compose"
else
  echo "Neither podman-compose nor docker is installed. Exiting."
  exit 1
fi

cd ${JUKEBOX_FRONTEND_DIR} && \
npm install && \
rm -rf ${JUKEBOX_FRONTEND_DIR}/dist && \
npm run build && \
cd .. && \
${DOCKER_COMPOSE} up --build
