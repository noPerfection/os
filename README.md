# os-lib

`github.com/noPerfection/os` is a small Go utility library for the SDS framework. It wraps common operating-system concerns—paths, CLI flags, environment files, processes, and ports—so application code does not have to repeat that logic.

## Mushroom substrate

This library also exposes a [Mushroom](https://github.com/ahmetson/mushroom) substrate so OS helpers can be resolved through Mushroom URLs instead of direct Go imports.

Substrate URL: `pkg:os/$`, any url that satisfies it will be implemented using this substrate.

Registered packages (each maps to an internal module):

| Mushroom URL package | Go package |
|----------------------|------------|
| `pkg:os/path`        | `path`     |
| `pkg:os/env`         | `env`      |
| `pkg:os/net`         | `net`      |
| `pkg:os/process`     | `process`  |
| `pkg:os/arg`         | `arg`      |

Packages are addressed as `pkg:os/<package>` (not `pkg:os#<package>`). A `#module` suffix is not supported, but maybe in the future it might have modules, for example pkg:os/path#windows might deal with Windows OS related path functions. But for now nothing like that.

Good to use it for reading environment variables, to link the environment variables not just in the code, but in the configuration files as well: `*pkg:os/env?var=MY_SECRET_KEY` in the json file might embed it but not store in the configuration itself.

### Quick start

```go
import "github.com/noPerfection/os/substrate"

mycelium, err := substrate.Root("pkg:os/path")
if err != nil {
    log.Fatal(err)
}

dir, err := mycelium.Spore("*pkg:os/path?func=CurrentDir()")
```

Register the substrate in a shared soil:

```go
soil := &mushroom.Soil{}
_ = soil.AddSubstrate(substrate.New())
```

### `path` — `?func=…`

Supported functions (must include `()` in the URL):

- `CurrentDir()` — returns the current directory where app is running
- `FileName(path)` — strips away the directory and returns file, *e.g. /path/dir/file.txt -> file.txt*
- `NoExtension(filename)` — strips away the directory and file name and returns extension only, *e.g. /path/file.txt -> txt*
- `DirAndFileName(fileDir)` — returns `[dir, name]`
- `FileExist(path)` — returns true if the file exist, if passed directory it returns error
- `DirExist(path)`— returns true if the directory exist

`MakeDir()` and other helpers in the Go package are **not** registered.

Examples:

| URL | Result |
|-----|--------|
| `*pkg:os/path?func=CurrentDir()` | executable directory |
| `*pkg:os/path?func=FileName(/tmp/a.txt)` | `"a.txt"` |
| `*pkg:os/path?func=FileName()` | error — parameter required |
| `*pkg:os/path?func=MakeDir()` | error — not registered |
| `*pkg:os/path?var` | error — `path` has no `var` resources |
| `pkg:os#path?func=CurrentDir` | error — wrong package syntax |

### `env` — `?var=…`

Read an environment variable:

```
*pkg:os/env?var=MY_VAR
```

Optional query parameters run `env.LoadAnyEnv` before reading the variable:

| Parameter | Effect |
|-----------|--------|
| `LoadAnyEnv=true` | `env.LoadAnyEnv(true)` |
| `LoadAnyEnv=false` | `env.LoadAnyEnv(false)` |
| `LoadAnyEnv` (no value) | `env.LoadAnyEnv()` |
| `arg=true` | `env.LoadAnyEnv(true)` |

Returns an empty string when the variable is unset.

### `process` — `?func=…`

- `*pkg:os/process?func=CurrentPid` — current PID (`()` optional for this call)
- `*pkg:os/process?func=PortToPid(8080)` — PID listening on a TCP port

### `net` — `?func=…`

- `*pkg:os/net?func=GetFreePort()`
- `*pkg:os/net?func=IsPortUsed(127.0.0.1, 8080)`

### `arg` — `?func=…` and `?var=…`

Functions: `NewFlag`, `Flags`, `IsFlag`, `FlagExist`, `ExtractFlagName`, `ExtractFlagValue`, `FlagValue`, `EnvPaths`.

Variables:

- `*pkg:os/arg?var=prefix` → `"--"`
- `*pkg:os/arg?var=sep` → `"="`

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

- [github.com/ahmetson/mushroom](https://github.com/ahmetson/mushroom) — Mushroom URL substrate (`substrate` package)
- [github.com/joho/godotenv](https://github.com/joho/godotenv) — `.env` file loading and writing
- [github.com/phayes/freeport](https://github.com/phayes/freeport) — free port allocation
- [github.com/cakturk/go-netstat](https://github.com/cakturk/go-netstat) — TCP socket / process lookup
## Requirements

- Go 1.19+
