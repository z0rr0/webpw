#!/usr/bin/env bash

# cd golang
# docker build -t golang:webpw .
CONTAINER="golang:webpw"
BIN_FILE="webpw.alpine"
# run from docker folder
SOURCES="`realpath ..`"

rm -f ../${BIN_FILE}

/usr/bin/docker run --rm --user `id -u ${USER}`:`id -g ${USER}` \
    --volume ${SOURCES}:/usr/p \
    --workdir /usr/p \
    --env GOCACHE=/tmp/.cache \
    --env CGO_ENABLED=0 \
    ${CONTAINER} go build -o ${BIN_FILE}

if [[ $? -gt 0 ]]; then
	echo "ERROR: build container"
	exit 1
fi
