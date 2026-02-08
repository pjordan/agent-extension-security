# Installation

The `agentsec` CLI is a single, statically linked binary with no runtime dependencies. All three installation methods below produce the same binary — choose the one that best fits your environment.

!!! note "Prerequisites"
    - **Go 1.23+** is required for Option 1 (`go install`) and Option 3 (build from source).
    - **No Go toolchain needed** for Option 2 (pre-built release binaries).
    - **Supported platforms:** Linux (amd64, arm64), macOS (amd64, arm64), Windows (amd64).

---

## Option 1: Install with `go install`

The fastest path if you already have a Go toolchain:

```bash
go install github.com/pjordan/agent-extension-security/cmd/agentsec@latest
```

!!! tip "Binary not found?"
    If your shell cannot find `agentsec` after installation, ensure `$(go env GOPATH)/bin` is on your `PATH`:

    ```bash
    export PATH="$(go env GOPATH)/bin:$PATH"
    ```

    Add this line to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.) to make it permanent.

Confirm the installation:

```bash
agentsec version
```

---

## Option 2: Download pre-built binaries

Use this method in environments without Go, in CI pipelines, or on air-gapped systems.

Pre-built binaries for every supported platform are published on the [GitHub Releases](https://github.com/pjordan/agent-extension-security/releases) page.

**Artifact naming convention:**

| Platform | Archive name |
|----------|-------------|
| Linux / macOS | `agentsec_<version>_<os>_<arch>.tar.gz` |
| Windows | `agentsec_<version>_windows_amd64.zip` |
| Checksums | `checksums.txt` |

!!! tip "Replace `<version>`"
    In the commands below, replace `<version>` with the release tag you downloaded (e.g., `v0.1.0`).

### Verify checksums

Before installing, verify the archive integrity against the published checksums. This confirms the download has not been tampered with or corrupted in transit.

=== "Linux"

    ```bash
    sha256sum -c checksums.txt
    ```

=== "macOS"

    ```bash
    shasum -a 256 -c checksums.txt
    ```

=== "Windows PowerShell"

    ```powershell
    # Compare the output against the corresponding entry in checksums.txt
    Get-FileHash .\agentsec_<version>_windows_amd64.zip -Algorithm SHA256
    ```

!!! warning "Use the official checksums"
    Always download `checksums.txt` from the [GitHub Releases](https://github.com/pjordan/agent-extension-security/releases) page — the same source as the archive itself. Do not rely on checksums hosted elsewhere.

### Extract and install

=== "Linux / macOS"

    ```bash
    tar -xzf agentsec_<version>_<os>_<arch>.tar.gz
    sudo install -m 0755 agentsec /usr/local/bin/agentsec
    agentsec version
    ```

=== "Windows PowerShell"

    ```powershell
    Expand-Archive .\agentsec_<version>_windows_amd64.zip -DestinationPath .
    Move-Item .\agentsec.exe "$env:USERPROFILE\bin\agentsec.exe"
    agentsec version
    ```

---

## Option 3: Build from source

Build from source when you want to audit the code, contribute changes, or customize the build.

=== "Using make"

    ```bash
    git clone https://github.com/pjordan/agent-extension-security.git
    cd agent-extension-security
    make build
    ./bin/agentsec version
    ```

=== "Using go build"

    ```bash
    git clone https://github.com/pjordan/agent-extension-security.git
    cd agent-extension-security
    go build -trimpath -o ./bin/agentsec ./cmd/agentsec
    ./bin/agentsec version
    ```

!!! tip "`make install`"
    Running `make install` copies the binary to `$(go env GOPATH)/bin` so it is available system-wide without a path prefix.

---

## Verify your installation

Regardless of which method you used, confirm that `agentsec` is available and working:

```bash
agentsec version
```

You should see output containing the version number and build metadata. If the command is not found, revisit the PATH instructions for your chosen installation method.

Having trouble? Check the [Troubleshooting](troubleshooting.md) guide.

---

## Next steps

- **Building an extension?** Start with the [Creator Quickstart](creating/quickstart.md).
- **Installing an extension?** Start with the [Consumer Quickstart](consuming/quickstart.md).
- **Exploring commands?** See the [CLI Reference](reference/cli.md) for the full command list.
- **Looking for examples?** Browse [Examples & Policies](examples.md) for sample manifests and policy templates.
