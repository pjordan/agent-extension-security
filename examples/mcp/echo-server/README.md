# echo-server example

A minimal MCP server that echoes input. Demonstrates:

- **MCP server type** (`type: "mcp-server"`)
- **Subprocess permission** (`process.allow_subprocess: true`) — needed to run the server process
- **Minimal permissions** — no file access, no network, no shell

## Package, sign, and install

```bash
# 1. Package
agentsec package examples/mcp/echo-server --out echo-server.aext

# 2. Validate the manifest
agentsec manifest validate examples/mcp/echo-server/aem.json

# 3. Sign
agentsec sign echo-server.aext --key devkey.json --out echo-server.sig.json

# 4. Install
agentsec install echo-server.aext \
  --sig echo-server.sig.json --pub devkey.json \
  --aem examples/mcp/echo-server/aem.json \
  --policy examples/policies/permissive.json \
  --dest ./installed

# Or in dev mode
agentsec install echo-server.aext --dev \
  --aem examples/mcp/echo-server/aem.json \
  --dest ./installed
```

## Notes

- This example uses `allow_subprocess: true` because an MCP server runs as a
  subprocess of the host agent. The `allow_shell` permission is kept `false`
  since the server doesn't need arbitrary shell access.
- In production, MCP servers would typically run in a container for sandboxing.
