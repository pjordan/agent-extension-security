# Install And Release Verification

## Download from GitHub Releases

Release binaries are published by OS/arch under:

- [GitHub Releases](https://github.com/pjordan/agent-extension-security/releases)

Each release includes per-artifact checksum files and a `checksums.txt` aggregate file.

## Verify checksums

On Linux/macOS:

```bash
sha256sum -c checksums.txt --ignore-missing
```

On macOS with `shasum`:

```bash
shasum -a 256 agentsec_darwin_arm64
```

On Windows PowerShell:

```powershell
Get-FileHash .\agentsec_windows_amd64.exe -Algorithm SHA256
```

## Optional: build from source

```bash
git clone https://github.com/pjordan/agent-extension-security.git
cd agent-extension-security
make build
./bin/agentsec version
```
