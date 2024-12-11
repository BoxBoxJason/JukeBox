# Stop script execution if any command fails
$ErrorActionPreference = 'Stop'

# Define the directory variables
$JUKEBOX_SCRIPT_DIR = Split-Path -Parent -Path $MyInvocation.MyCommand.Path
Write-Host "Le chemin du script est : $JUKEBOX_SCRIPT_DIR"

# Build the Docker image from the Dockerfile (Containerfile)
docker build $JUKEBOX_SCRIPT_DIR -t jukebox:dev -f Containerfile

# Run the Docker container, mounting the directory without spaces and exposing port 3000
docker run -d --name jukebox -v "${JUKEBOX_SCRIPT_DIR}:/opt/jukebox" -p 3000:3000 jukebox:dev tail -f /dev/null

