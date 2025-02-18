#!/bin/sh

# Start the icecast server
icecast -c /etc/icecast.xml &

# Wait for the icecast server to start
sleep 5 # yeah this is a hack

# Start the ffmpeg loop stream (test purpose)
ffmpeg -re -stream_loop -1 -i test.m4a -content_type audio/aac -f adts icecast://source:source_password@localhost:3001/stream
