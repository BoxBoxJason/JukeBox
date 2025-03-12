#!/bin/sh

set -e

/usr/bin/musicgpt-api &

/bin/sh -c musicgpt "$@" --
