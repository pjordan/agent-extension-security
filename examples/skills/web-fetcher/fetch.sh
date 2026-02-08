#!/bin/sh
# fetch.sh — Fetch data from an approved API endpoint.
# This script is called by the web-fetcher skill.
# NOTE: The scanner will flag this curl usage — that is intentional
# to demonstrate the scanning feature.

set -e

BASE_URL="https://api.example.com"

if [ -z "$1" ]; then
  echo "usage: fetch.sh <path>" >&2
  exit 1
fi

curl -sS "${BASE_URL}$1"
