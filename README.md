# HexLoader

![CI/CD](https://github.com/TEgor1212/HexLoader/actions/workflows/ci.yml/badge.svg)

A **GUI hex file viewer/uploader for Arduino boards**, built with **Go** and **Fyne**.

HexLoader lets you select a `.hex` file, choose a connected Arduino board (with automatic FQBN detection via `arduino-cli`), and upload the firmware — all from a desktop application.

---

## Run

```bash
go run .
```

## Build

```bash
go build -o hexloader .
```

> **Note:** Fyne requires CGO and system GUI libraries (OpenGL, X11, etc.). See the [Fyne prerequisites](https://docs.fyne.io/started/) for your OS.

---

## CI/CD Pipeline

The workflow in `.github/workflows/ci.yml` runs on pushes to `main`, PRs, and tag pushes (`v*`).

### CI jobs (push / PR)

| Job | What it does |
|---|---|
| **Lint** | `gofmt -l .` (formatting check), `go vet ./...` — with GUI dependencies installed on Linux |
| **Test** | `go test -v -count=1 ./...` — unit tests with GUI dependencies |

### Release job (tag push `v*`)

After both CI jobs pass, the **release** job:

- Builds binaries for **Linux** (`hexloader-linux-x64`), **macOS** (`hexloader-macos-arm64`), **Windows** (`hexloader-windows-x64.exe`)
- Uploads each binary as a workflow artifact
- Publishes all binaries to **GitHub Releases** via `softprops/action-gh-release@v2`

Example:
```bash
git tag v1.0.0
git push origin v1.0.0
```

This triggers the release pipeline and creates a Release with the compiled binaries.