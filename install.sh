#! /bin/bash

set -o errexit

if ! which git || ! which wget; then
  if which apk; then
    apk add --no-cache git wget
  elif which apt-get; then
    apt-get update && apt-get install -y git wget
  fi
fi

INSTALLATION_DIR=/usr/local/cachr
git clone https://github.com/scripbox/cachr.git $INSTALLATION_DIR
ln -s $INSTALLATION_DIR/cache_run /usr/bin/cache_run
ln -s $INSTALLATION_DIR/cache_get /usr/bin/cache_get
wget https://github.com/scripbox/cachr/releases/download/v1.1/cachr -O /usr/bin/cachr
chmod +x /usr/bin/cache_run /usr/bin/cache_get /usr/bin/cachr
