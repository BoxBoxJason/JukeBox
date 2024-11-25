# Stop script execution if any command fails
$ErrorActionPreference = 'Stop'

# Define the directory variables
$JUKEBOX_SCRIPT_DIR = Split-Path -Parent -Path $MyInvocation.MyCommand.Path

docker build $JUKEBOX_SCRIPT_DIR -t jukebox:dev -f Containerfile

docker run -d -p 3000:3000 jukebox:dev
