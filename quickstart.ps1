# Stop script execution if any command fails
$ErrorActionPreference = 'Stop'

# Define the directory variables
$JUKEBOX_SCRIPT_DIR = Split-Path -Parent -Path $MyInvocation.MyCommand.Path
$JUKEBOX_FRONTEND_DIR = Join-Path -Path $JUKEBOX_SCRIPT_DIR -ChildPath "frontend"
$JUKEBOX_BIN_DIR = Join-Path -Path $JUKEBOX_SCRIPT_DIR -ChildPath "bin"
$JUKEBOX_MAIN_GO = Join-Path -Path $JUKEBOX_SCRIPT_DIR -ChildPath "cmd\server\main.go"
$JUKEBOX_EXECUTABLE = Join-Path -Path $JUKEBOX_BIN_DIR -ChildPath "jukebox"

# Change directory to the frontend directory
Set-Location -Path $JUKEBOX_FRONTEND_DIR

# Run npm install and build
npm --prefix $JUKEBOX_FRONTEND_DIR install
npm --prefix $JUKEBOX_FRONTEND_DIR run build

# Change directory back to the script directory
Set-Location -Path $JUKEBOX_SCRIPT_DIR

# Tidy up Go modules
go mod tidy

# Create the bin directory if it doesn't exist
if (!(Test-Path -Path $JUKEBOX_BIN_DIR)) {
    New-Item -ItemType Directory -Path $JUKEBOX_BIN_DIR
}

# Build the Go application
go build -o $JUKEBOX_EXECUTABLE $JUKEBOX_MAIN_GO

# Run the jukebox application
& $JUKEBOX_EXECUTABLE
