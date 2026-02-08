#!/bin/sh
# server.sh — Minimal echo responder for MCP server example.
# Reads lines from stdin and echoes them back.

while IFS= read -r line; do
  echo "$line"
done
