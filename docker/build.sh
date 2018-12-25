#!/usr/bin/env bash

BUILD_CONTAINER="golang:webpw"
TAG="zorro/webpw"
BIN_FILE="webpw.alpine"
SCRIPT_DIR=`realpath $0 | xargs dirname`
SOURCES=`realpath ${SCRIPT_DIR}/..`

# create custom golang image if it doesn't exist
docker inspect ${BUILD_CONTAINER} > /dev/null
if [[ $? -gt 0 ]]; then
	echo "INFO: image ${BUILD_CONTAINER} does not exist"
	cd ${SOURCES}/docker/golang
	docker build -t ${BUILD_CONTAINER} .
fi

# delete previous result file
rm -f ${SOURCES}/${BIN_FILE}

echo "INFO: build a program inside golang docker image"
/usr/bin/docker run --rm --user `id -u ${USER}`:`id -g ${USER}` \
    --volume ${SOURCES}:/usr/p \
    --workdir /usr/p \
    --env GOCACHE=/tmp/.cache \
    --env CGO_ENABLED=0 \
    ${BUILD_CONTAINER} go build -o ${BIN_FILE}

if [[ $? -gt 0 ]]; then
	echo "ERROR: build container"
	exit 1
fi

echo "INFO: build docker image with prepared binary file"
cd ${SOURCES}
docker build -t ${TAG} -f docker/Dockerfile .
