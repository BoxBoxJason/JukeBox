# Get the directory of the script
$JUKEBOX_SCRIPT_DIR = Split-Path -Path $MyInvocation.MyCommand.Definition -Parent
$JUKEBOX_FRONTEND_DIR = Join-Path -Path $JUKEBOX_SCRIPT_DIR -ChildPath "frontend"

# Change to the frontend directory, install dependencies, and build
Set-Location -Path $JUKEBOX_FRONTEND_DIR
npm install
npm run build

# Change back to the script's directory, tidy Go modules, and run the server
Set-Location -Path $JUKEBOX_SCRIPT_DIR
go mod tidy
go run ./cmd/server/main.go
