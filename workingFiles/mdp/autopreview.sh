#! /bin/bash

# preview the markdown file every time it changes (the md5 checksum changes)
# we recalculate the MD5 checksum every 5 seconds, and trigger the mdp tool if the content has changed

FHASH=`md5sum $1`
while true; do
  NHASH=`md5sum $1`
  if [ "$NHASH" != "$FHASH" ]; then
    ./mdp -file $1
    FHASH=$NHASH
  fi
  sleep 5
done