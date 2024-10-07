#!/bin/bash

jukebox() {
    CONTAINER_OPERATOR="podman"
    IMAGE_NAME="jukebox"
    IMAGE_TAG="latest"
    CONTAINER_NAME="jukebox"
    CONTAINER_PORT="3000"
    NETWORK_NAME="jukebox"
    CONTAINER_FILE="./Containerfile"

    go() {
        # Delete the already existing container & image
        $CONTAINER_OPERATOR rm -f $CONTAINER_NAME
        $CONTAINER_OPERATOR rmi $IMAGE_NAME:$IMAGE_TAG

        # Build the image
        $CONTAINER_OPERATOR build -t $IMAGE_NAME:$IMAGE_TAG -f $CONTAINER_FILE .
        if [ $? -ne 0 ]; then
            echo "Failed to build the image"
            return 1
        fi

        # Create the network if it doesn't exist
        if [[ -z $($CONTAINER_OPERATOR network ls --filter name=^$NETWORK_NAME$ --format="{{ .Name }}") ]]; then
            $CONTAINER_OPERATOR network create $NETWORK_NAME
        fi

        # Run the container
        $CONTAINER_OPERATOR run -d --name $CONTAINER_NAME -p $CONTAINER_PORT:$CONTAINER_PORT --network $NETWORK_NAME $IMAGE_NAME:$IMAGE_TAG
    }

    stop() {
        $CONTAINER_OPERATOR stop $CONTAINER_NAME
    }

    resume() {
        $CONTAINER_OPERATOR start $CONTAINER_NAME
    }
}

# Check arguments and call the respective function
if [ "$1" == "go" ]; then
    jukebox
    go
elif [ "$1" == "stop" ]; then
    jukebox
    stop
elif [ "$1" == "resume" ]; then
    jukebox
    resume
else
    echo "Usage: $0 {quickstart|stop|resume}"
fi
