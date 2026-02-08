#!/bin/sh
# read-config.sh — Read a config file and print its contents.
# This script is called by the file-reader skill.

set -e

if [ -z "$1" ]; then
  echo "usage: read-config.sh <path>" >&2
  exit 1
fi

FILE="$1"

if [ ! -f "$FILE" ]; then
  echo "error: file not found: $FILE" >&2
  exit 1
fi

cat "$FILE"
