# Install & Release Verification

Both extension creators and consumers need the `agentsec` CLI. Choose the installation method that fits your environment.

## Option 1: Install with `go install`

```bash
go install github.com/pjordan/agent-extension-security/cmd/agentsec@latest
agentsec version
```

If `agentsec` is not found, add `$(go env GOPATH)/bin` to your `PATH`.

## Option 2: Download release archives

Release archives are published at:

- [GitHub Releases](https://github.com/pjordan/agent-extension-security/releases)

Artifact naming:

- Linux/macOS: `agentsec_<version>_<os>_<arch>.tar.gz`
- Windows: `agentsec_<version>_windows_amd64.zip`
- Checksums: `checksums.txt`

### Verify checksums

=== "Linux"

    ```bash
    sha256sum -c checksums.txt
    ```

=== "macOS"

    ```bash
    shasum -a 256 agentsec_v0.1.0_darwin_arm64.tar.gz
    ```

=== "Windows PowerShell"

    ```powershell
    Get-FileHash .\agentsec_v0.1.0_windows_amd64.zip -Algorithm SHA256
    ```

### Install from archive

=== "Linux / macOS"

    ```bash
    tar -xzf agentsec_v0.1.0_darwin_arm64.tar.gz
    install -m 0755 agentsec /usr/local/bin/agentsec
    ```

=== "Windows PowerShell"

    ```powershell
    Expand-Archive .\agentsec_v0.1.0_windows_amd64.zip -DestinationPath .
    Move-Item .\agentsec.exe "$env:USERPROFILE\bin\agentsec.exe"
    ```

## Option 3: Build from source

```bash
git clone https://github.com/pjordan/agent-extension-security.git
cd agent-extension-security
go build -trimpath -o ./bin/agentsec ./cmd/agentsec
./bin/agentsec version
```

`make build` and `make install` are convenience wrappers around the same Go toolchain flow.

## Next steps

- **Building an extension?** Start with the [Creator Quickstart](creating/quickstart.md)
- **Installing an extension?** Start with the [Consumer Quickstart](consuming/quickstart.md)
