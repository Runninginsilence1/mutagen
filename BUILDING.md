# Building

Mutagen's build is slightly unique because it needs to cross-compile agent
binaries for remote platforms (with cgo support in the case of macOS) and then
generate a bundle of these binaries to ship alongside the Mutagen CLI. As such,
using `go get` or `go install` to acquire Mutagen will result in an incomplete
installation, and users should instead download the release builds from the
[releases page](https://github.com/mutagen-io/mutagen/releases/latest) or
[install Mutagen](https://mutagen.io/documentation/introduction/installation)
via [Homebrew](https://brew.sh/).

However, Mutagen can be built locally for testing and development. Mutagen
relies on the Go toolchain's module support, so make sure that you have Go
module support enabled.

Individual Mutagen executables can be built normally using the Go toolchain, but
a script is provided to ensure a normalized build, manage cross-compiled builds
and agent bundle creation, and perform code signing on macOS. To see information
about the build script, run:

    go run scripts/build.go --help

The build script can do four different types of builds: `local` (with support
for the local system only), `slim` (the default - with support for a selection
of common platforms used in testing), `release` (used for generating complete
release artifacts), and `release-slim` (used for generating complete release
artifacts for a selection of common platforms used in testing). macOS is
currently the only platform that supports doing `release` builds, because the
macOS binaries require cgo support for filesystem monitoring.

All artifacts from the build are placed in a `build` directory at the root of
the Mutagen source tree. As a convenience, artifacts built for the current
platform are placed in the root of the build directory for easy testing, e.g.:

    go run scripts/build.go
    build/mutagen --help


## Quick build

Individual binaries can be built directly with `go build`, provided the correct
build tag is specified:

    go build -tags mutagencli -o build/mutagen ./cmd/mutagen/
    go build -tags mutagenagent -o build/mutagen-agent ./cmd/mutagen-agent/

All build output goes into the `build/` directory (which is gitignored).


## Build tags

| Tag | Purpose | Required for |
|-----|---------|-------------|
| `mutagencli` | Registers built-in transport protocols (SSH/Docker/Local) and disables the tagcheck panic guard | Building the standalone CLI binary |
| `mutagenagent` | Disables the agent tagcheck panic guard | Building the agent binary |
| `mutagensspl` | Enables SSPL-licensed enhancements (xxh128 hashing, zstandard compression, Linux fanotify watching) | Optional |
| `mutagenfanotify` | Enables Linux fanotify filesystem watching (requires `mutagensspl`) | Optional, Linux only |
| `mutagensidecar` | Builds the sidecar binary | Building the sidecar |

**What happens without the required tag?** Each binary entry point contains a
`tagcheck.go` guard file with a build constraint of `!mutagencli` or
`!mutagenagent`. Without the tag, this file is compiled in and the program
panics at startup. This is by design — Mutagen's command-line code can be
embedded into third-party tools that register their own protocol handlers, and
the tagcheck prevents such embedded binaries from being mistakenly run as a
standalone CLI.


## SSPL-enhanced builds

    go build -tags mutagencli,mutagensspl -o build/mutagen ./cmd/mutagen/
    go build -tags mutagenagent,mutagensspl,mutagenfanotify -o build/mutagen-agent ./cmd/mutagen-agent/


## Windows cross-compilation

Two helper scripts are provided in the `build/` directory for producing Windows
binaries from macOS or Linux:

**`build/build_with_windows_cli.sh`** — Runs the official slim build (native CLI
\+ multi-platform agent bundle), then cross-compiles a Windows CLI. This is the
recommended approach:

    ./build/build_with_windows_cli.sh

Output in `build/`: `mutagen` (native), `mutagen.exe` (Windows), and
`mutagen-agents.tar.gz` (shared agent bundle).

**`build/windows_cross_compile.sh`** — Standalone Windows cross-compilation that
builds both the Windows CLI and a self-contained agent bundle from scratch,
without running the full official build:

    ./build/windows_cross_compile.sh                 # amd64 (default)
    GOARCH=arm64 ./build/windows_cross_compile.sh    # arm64

Output in `build/windows_<arch>/`: `mutagen.exe` and `mutagen-agents.tar.gz`.

> **Important:** `mutagen-agents.tar.gz` must be placed in the same directory as
> `mutagen.exe`. The CLI locates this bundle at runtime to deploy agents to
> remote endpoints.


## Known build warnings

On macOS, you may see the following compiler warning, which can be safely
ignored:

    'FSEventStreamScheduleWithRunLoop' is deprecated: first deprecated in macOS 13.0

This originates from the upstream `github.com/mutagen-io/fsevents` dependency
using a deprecated macOS API. It does not affect functionality.


## Protocol Buffers code generation

Mutagen uses Protocol Buffers extensively, and as such needs to generate Go code
from `.proto` files. To avoid the need for developers (and CI systems) to have
the Protocol Buffers compiler installed, generated code is checked into the
repository. If a `.proto` file is modified, code can be regenerated by running

    go generate ./pkg/...

in the root of the Mutagen source tree.

The `go generate` commands used by Mutagen rely on Go module support being
enabled. You will also need to have the `protoc` compiler (with support for
Protocol Buffers 3) available in your path, but not the Go generator, which will
be built as part of the `go generate` command.
