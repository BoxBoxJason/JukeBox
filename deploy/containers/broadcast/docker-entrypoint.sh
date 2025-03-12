#!/bin/sh

set -e

# Replace the environment variables in the icecast.xml file
envsubst < /etc/icecast_template.xml > /etc/icecast.xml

# Start the icecast server as the icecast user
su icecast -c "icecast -c /etc/icecast.xml" &

DIRECTORY_TO_WATCH="/opt/musicgpt/generated"
MUSIC_FILE_EXTENSION=".wav"
FFMPEG_PID_FILE="/tmp/ffmpeg.pid"
CURRENT_FILE=""

inotifywait -m -e create --format "%f" "$DIRECTORY_TO_WATCH" | while read NEW_FILE
do
  if [[ "$NEW_FILE" == *"$MUSIC_FILE_EXTENSION" ]]; then
    # If ffmpeg is already running, kill it and delete the former file
    if [ -f "$FFMPEG_PID_FILE" ]; then
      FFMPEG_PID=$(cat "$FFMPEG_PID_FILE")
      if kill -0 "$FFMPEG_PID" 2>/dev/null; then
        kill -9 "$FFMPEG_PID"
      fi
      rm -f "$FFMPEG_PID_FILE"
      if [ -n "$CURRENT_FILE" ]; then
        rm -f "$CURRENT_FILE"
      fi
    fi

    # Set the current file to the new file
    CURRENT_FILE="$DIRECTORY_TO_WATCH/$NEW_FILE"

    # Stream the new file using ffmpeg and save the PID
    echo "Playing $CURRENT_FILE" >> /tmp/ffmpeg.log
    ffmpeg -re -stream_loop -1 -i "$CURRENT_FILE" -f mp3 icecast://${ICECAST_USER}:${ICECAST_PASSWORD}@localhost:3001/stream &
    echo $! > "$FFMPEG_PID_FILE"
  fi
done
