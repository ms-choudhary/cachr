#! /bin/bash

set -o errexit

if ! which git || ! which wget; then
  if which apk; then
    apk add --no-cache git wget
  elif which apt-get; then
    apt-get update && apt-get install -y git wget
  fi
fi

wget https://raw.githubusercontent.com/scripbox/cachr/master/cache_run -O /usr/bin/cache_run
wget https://raw.githubusercontent.com/scripbox/cachr/master/cache_get -O /usr/bin/cache_get
wget https://github.com/scripbox/cachr/releases/download/v1.1/cachr-ubuntu -O /usr/bin/cachr
chmod +x /usr/bin/cache_run /usr/bin/cache_get /usr/bin/cachr
