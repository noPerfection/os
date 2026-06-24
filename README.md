# os-lib

`github.com/noPerfection/os` is a small Go utility library for the SDS framework. It wraps common operating-system concerns—paths, CLI flags, environment files, processes, and ports—so application code does not have to repeat that logic.

## Packages

### `path` — filesystem helpers

- Resolve the executable's directory (`CurrentDir`) and build absolute paths relative to it (`AbsDir`)
- File and directory existence checks that distinguish files from directories (`FileExist`, `DirExist`)
- Path parsing: basename, strip extension, split directory and name (`FileName`, `NoExtension`, `DirAndFileName`)
- Create nested directories (`MakeDir`)
- Platform-specific binary paths: `.exe` on Windows, plain name on Unix (`BinPath`)

### `arg` — SDS-style CLI flags

- Flags use `--name` or `--name=value` (not the standard library `flag` package)
- Helpers to list flags, test existence, read values, and build flag strings
- Positional arguments ending in `.env` are treated as environment file paths (`EnvPaths`)

### `env` — `.env` loading and writing

- `LoadAnyEnv()` reads `.env` paths from the command line (via `arg.EnvPaths()`), resolves them relative to the executable directory, and loads them with `godotenv` into process environment variables (for use with `app/config.Config` in the wider framework)
- `WriteEnv()` writes key/value data to a `.env` file (any type implementing `env.KeyValue`)

### `process` — process identity and port mapping

- `CurrentPid()` returns the current process ID
- `PortToPid()` finds which process is listening on a given TCP port (via `go-netstat`)

### `net` — port availability

- `GetFreePort()` picks an unused TCP port
- `IsPortUsed()` checks whether something is listening on a host:port via a TCP dial

## Summary

In short, this library is a thin **OS/platform utility layer** for SDS applications: where the binary lives, how CLI flags and `.env` files are parsed, and how to work with ports and processes.

## Dependencies

- [github.com/joho/godotenv](https://github.com/joho/godotenv) — `.env` file loading and writing
- [github.com/phayes/freeport](https://github.com/phayes/freeport) — free port allocation
- [github.com/cakturk/go-netstat](https://github.com/cakturk/go-netstat) — TCP socket / process lookup
## Requirements

- Go 1.19+
