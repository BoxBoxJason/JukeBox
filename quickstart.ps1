function Jukebox {
    $global:CONTAINER_OPERATOR = "docker"
    $global:IMAGE_NAME = "jukebox"
    $global:IMAGE_TAG = "latest"
    $global:CONTAINER_NAME = "jukebox"
    $global:CONTAINER_PORT = 3000
    $global:NETWORK_NAME = "jukebox"
    $global:CONTAINER_FILE = "./Containerfile"

    function Go {
        # Remove the existing container and image if they exist
        & $CONTAINER_OPERATOR rm -f $CONTAINER_NAME
        & $CONTAINER_OPERATOR rmi "$IMAGE_NAME:$IMAGE_TAG"

        # Build the image
        & $CONTAINER_OPERATOR build -t "$IMAGE_NAME:$IMAGE_TAG" -f $CONTAINER_FILE .
        if ($LASTEXITCODE -ne 0) {
            Write-Host "Failed to build the image"
            return 1
        }

        # Create the network if it doesn't exist
        if (-not (& $CONTAINER_OPERATOR network ls --filter "name=^$NETWORK_NAME$" --format "{{ .Name }}")) {
            & $CONTAINER_OPERATOR network create $NETWORK_NAME
        }

        # Run the container
        & $CONTAINER_OPERATOR run -d --name $CONTAINER_NAME -p "$CONTAINER_PORT:$CONTAINER_PORT" --network $NETWORK_NAME "$IMAGE_NAME:$IMAGE_TAG"
    }

    function Stop {
        & $CONTAINER_OPERATOR stop $CONTAINER_NAME
    }

    function Resume {
        & $CONTAINER_OPERATOR start $CONTAINER_NAME
    }
}

# Check arguments and call the respective function
param (
    [string]$Action
)

switch ($Action) {
    "go" {
        Jukebox
        Go
    }
    "stop" {
        Jukebox
        Stop
    }
    "resume" {
        Jukebox
        Resume
    }
    default {
        Write-Host "Usage: script.ps1 {go|stop|resume}"
    }
}
