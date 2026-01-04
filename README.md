# LogDrift

[![GitHub release](https://img.shields.io/github/v/release/adermanbuilds/logdrift)](https://github.com/adermanbuilds/logdrift/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/adermanbuilds/logdrift)](https://goreportcard.com/report/github.com/adermanbuilds/logdrift)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Fast, lightweight log anomaly detection for the command line — single binary, no servers.

![LogDrift demo](./logdrift.gif)

## Features

- Real-time streaming from stdin
- ANSI color-coded output
- Heuristic anomaly scoring
- Structured log parsing (e.g., `ERROR [component]`)
- Configurable JSON patterns
- Component filtering
- Built-in runtime statistics
- Zero-configuration defaults

## Status

- **Version:** v0.1.0 (CLI only)
- GUI/TUI not implemented yet — roadmap in this README.

---

## Quick start

### Prerequisites

- Go 1.21+ (recommended)

### Installation

#### Download Binary

Download the latest release for your platform:
- [macOS (Intel)](https://github.com/adermanbuilds/logdrift/releases/latest/download/logdrift-darwin-amd64)
- [macOS (Apple Silicon)](https://github.com/adermanbuilds/logdrift/releases/latest/download/logdrift-darwin-arm64)
- [Linux (x64)](https://github.com/adermanbuilds/logdrift/releases/latest/download/logdrift-linux-amd64)
- [Windows (x64)](https://github.com/adermanbuilds/logdrift/releases/latest/download/logdrift-windows-amd64.exe)

#### Install via Go

```bash
go install github.com/adermanbuilds/logdrift@latest
```

#### Build from Source

```bash
git clone https://github.com/adermanbuilds/logdrift.git
cd logdrift
make build
```

#### Alternatively, for module-aware builds without make

1. Initialize modules (only once per repo):

```bash
go mod init github.com/adermanbuilds/logdrift
go mod tidy
```

2. Build:

```bash
go build -o logdrift *.go
```

3. Run:

```bash
tail -f /var/log/nginx/access.log | ./logdrift
cat /var/log/system.log | ./logdrift
cat app.log | ./logdrift --all
cat app.log | ./logdrift --component=parser
tail -f app.log | ./logdrift --anomalies-only
```

### Generate default config

```bash
./logdrift --generate-config > my-config.json
cat app.log | ./logdrift --config=my-config.json
```

---

## What it detects

- **HTTP 5xx status patterns**
- **Error keywords:** error, exception, fatal, panic, failed
- **Slow responses** (default threshold ~500ms)
- **Connection issues:** refused, timeout, broken pipe
- **Suspicious activity:** blocked/denied, injection, brute-force indicators

---

## Usage

```text
Usage:
    logdrift [options] < logfile.log
    tail -f /var/log/app.log | logdrift [options]

Options:
    --config string         Path to config file (JSON)
    --all                   Show all logs (including DEBUG/INFO)
    --debug                 Show DEBUG logs
    --info                  Show INFO logs
    --no-color              Disable colored output
    --compact               Compact output (no anomaly highlights)
    --component string      Filter by component (e.g., 'parser')
    --min-level string      Minimum log level (default: WARN)
    --anomalies-only        Show only anomalies
    --generate-config       Generate default config file
    --version               Show version
```

---

## Example output

```text
ANOMALY (score: 110) [ERROR] 2026-01-03 23:06:01 ERROR  [notifier] failed to parse line 429
[WARN] 2026-01-03 23:06:04 WARN  [indexer] connection slow; retrying in 11255s
ANOMALY (score: 130) [FATAL] 2026-01-03 23:06:05 FATAL  [watcher] fatal: unrecoverable error

────────────────────────────────────────────────────────────
LogDrift Statistics
────────────────────────────────────────────────────────────
Runtime:         1s
Lines Processed: 10000
Anomalies:       794 (7.9%)
Errors:          794 (7.9%)
Warnings:        1176 (11.8%)
Rate:            10613 lines/sec
────────────────────────────────────────────────────────────
```

---

## Notes & limitations

- **Go modules:** add a `go.mod` (required for builds outside GOPATH).
- **Buffer limit:** `bufio.Scanner` max line length is set to ~1MB in `main.go`. Very long single log lines may be truncated — adjust scanner buffer or pre-process logs if needed.
- **Performance:** level-detection regexp is currently compiled per line in `ParseLine`; compiling once improves throughput (suggested patch available).
- **Compatibility:** code uses `slices` package (Go 1.21); ensure `go.mod` sets `go 1.21+` or replace with a small helper.

---

## Roadmap

### v0.2 — TUI

- Interactive terminal UI (bubbletea)
- Filtering, searching, realtime graphs

### v0.3 — Web dashboard

- Web UI (Go + HTMX)
- Local storage (SQLite), alerts & webhooks

---

## Contributing

**PRs welcome. Areas that need help:**

- More log format parsers (Apache, MySQL, PostgreSQL)
- Pattern library expansion
- Unit tests and CI workflow
- Improved sampling/learning for anomaly scoring

### Development tips

- Run linters and tests:

```bash
go vet ./...
go test ./...
```

- Add a lightweight GitHub Actions workflow to run `go test` and `go vet` on PRs.

---

## License

MIT — use freely for personal or commercial projects.

---

## Contact / repo

[LogDrift](https://github.com/adermanbuilds/logdrift)
