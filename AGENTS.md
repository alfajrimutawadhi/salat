# AGENTS.md

Small Go CLI (module `github.com/alfajrimutawadhi/salat`, Go 1.19) that shows Islamic prayer times and Hijri dates via the Aladhan API.

## Commands

- Build: `go build ./cmd` — entrypoint is `cmd/main.go` (single `main` package; there is no `Makefile`, tests, lint, or CI test workflow).
- Verify changes by running the built binary (`./salat now`, `date`, `calendar`, etc.). All commands except `version` call the live `api.aladhan.com/v1` API, so they need internet and a valid location in `~/.salat/config.json`.
- Releases: CI (`.github/workflows/release.yml`) builds `project_path: ./cmd` per OS/arch (linux/darwin/windows × amd64/arm64) when a GitHub release is created. Assets are `salat-<tag>-<os>-<arch>.tar.gz` for unix and `.zip` (containing `salat.exe`) for Windows — that naming is what the install scripts download. Nothing else runs in CI.
- Installers: `install.sh`/`uninstall.sh` (bash, Linux/macOS) and `install.ps1`/`uninstall.ps1` (PowerShell, Windows). The bash scripts detect Windows (MINGW/MSYS in Git Bash) and point the user to the `.ps1` scripts; `install.ps1` adds `~\.salat` to the user PATH.

## Gotchas

- Config is read/written at `<user home>/.salat/config.json` (resolved via `os.UserHomeDir()`; on Windows that is `%USERPROFILE%\.salat\config.json`), not the repo-root `config.json`. The repo-root `config.json` is only a default template shipped with release archives.
- The app version lives in **three places** that must be kept in sync and match the release tag: `cmd/main.go` (`var version`), `install.sh` (`version=`), and `install.ps1` (`$version =`).
- `common.HandleError` (common/error.go) prints and calls `os.Exit` — nothing after it runs. Do not expect error returns from the main flow.
- Command names are string constants in `constant/command.go`; dispatch is a `switch` on `os.Args[1]` in `cmd/main.go`. Update the `Help` text in the same file when adding commands.
- `set-location` and `time-mode` are interactive (read stdin in `cmd/main.go`), so they can't be run non-interactively.
- `salat.RequestAPI` (salat/api.go) has two modes via the `nec` param: `"s"` → `timingsByCity`, `"c"` → `calendarByCity`. Calculation method `method=3` is hardcoded.
