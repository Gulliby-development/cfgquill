<p align="center"><img src="docs/logo.svg" width="96" alt="" /></p>

# cfgquill

cfgquill is an offline command-line tool for people who read or edit Git configuration files such as `.gitconfig`, `.git/config`, and `.gitmodules`. It helps inspect values and identify syntax errors or questionable entries without asking Git to interpret the file. It reads only the files you name (or standard input when you use `-`); it does not discover repositories or modify configuration files.

[![CI](https://github.com/Gulliby-development/cfgquill/actions/workflows/ci.yml/badge.svg)](https://github.com/Gulliby-development/cfgquill/actions/workflows/ci.yml) [![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Install

Install the Go toolchain required by the project, then build from the project directory:

```sh
go build -o cfgquill .
```

On Windows, build a Windows executable with:

```powershell
go build -o cfgquill.exe .
```

The project uses only the Go standard library. The build needs the source checkout and Go toolchain; the cfgquill executable itself is offline and makes no network requests. Run `./cfgquill` on macOS or Linux. In PowerShell use `.\cfgquill.exe`, and in Command Prompt use `cfgquill.exe`; alternatively, put the executable on your `PATH` and use `cfgquill` as in the examples below. The `go.mod` file names the Go release (`go 1.27`) and the toolchain (`toolchain go1.27.0`) the project is built with; no other change to `go.mod` is needed.

## Usage

These examples use POSIX-style shell syntax and assume `cfgquill` is on your `PATH`. A shell removes the surrounding quotes used to pass a single argument; use the equivalent quoting syntax in PowerShell or Command Prompt. Replace the sample paths with the configuration files you want to inspect.

```sh
cfgquill list .gitconfig
cfgquill sections .git/config
cfgquill get --all remote.origin.fetch .gitconfig
cfgquill bool core.filemode .git/config
cfgquill unquote value
cfgquill quote "two words"
cfgquill check .gitconfig .gitmodules
```

For example, with this `t.cfg`:

```ini
[core]
	editor = vim
	Pager = "less -R"   ; the pager
[remote "origin"]
	fetch = +refs/heads/*:refs/remotes/origin/*
[Alias]
	co
```

`cfgquill list t.cfg` prints (exit 0):

```text
LINE  KEY                  VALUE
2     core.editor          "vim"
3     core.pager           "less -R"
5     remote.origin.fetch  "+refs/heads/*:refs/remotes/origin/*"
7     alias.co             -
```

`cfgquill get core.missing t.cfg` prints nothing and exits 1; `cfgquill bool core.editor t.cfg` fails with `cfgquill: t.cfg:2: bad boolean "vim" for core.editor` and exit 2.

A file argument of `-` reads standard input instead. For example, in a POSIX-style shell, `cfgquill list - < .gitconfig` sends that file to standard input. Standard input can be named only once in a `check` command.

cfgquill accepts strict UTF-8 configuration text. Its blank characters are space, TAB, LF, and CR; other Unicode characters that may look like spaces or line breaks remain ordinary text. In values, `#` and `;` start comments outside double quotes, quoted text can retain blanks, and a backslash followed by LF continues a value onto the next line. These are the configuration rules implemented by cfgquill; the tool does not claim to implement every Git configuration feature.

## Commands

`cfgquill --version` prints one line, `cfgquill built with` and the Go release that built the executable (for example `cfgquill built with go1.27.0`), and exits 0; it is accepted only as the single argument.

The command names and argument forms are printed by `cfgquill --help`:

```text
usage: cfgquill list FILE
usage: cfgquill sections FILE
usage: cfgquill get [--all] KEY FILE
usage: cfgquill bool KEY FILE
usage: cfgquill unquote TEXT
usage: cfgquill quote TEXT
usage: cfgquill check FILE...
```

- `list FILE` prints a rune-aligned table of entries with their line, stored key, and quoted display value. An entry without a value is shown as `-`; an explicitly empty value is shown as `""`. A file with no entries produces no table.
- `sections FILE` prints each section header, its stored section and optional subsection, and its entry count. It produces no table when there are no headers.
- `get [--all] KEY FILE` prints the last value for KEY followed by LF, or every matching value in file order with `--all`. An entry without a value and an explicitly empty value each produce an empty output line. A missing key is not an error: the command exits 1 and prints nothing. Section and variable parts of a command-line key are ASCII case-insensitive; a subsection is case-sensitive.
- `bool KEY FILE` prints `true` or `false` for the last matching entry. An entry without a value is true, while an explicitly empty value is false. A missing key exits 1 without output; a present value that is not a recognized boolean is an error.
- `unquote TEXT` interprets value text and writes the resulting value without adding LF. `quote TEXT` writes the Git-style representation of a value, also without adding LF.
- `check FILE...` reports syntax errors and findings for duplicate entries, empty sections, blanks at value edges, selected invisible characters, and bare CR characters. A syntax error is the sole finding for that file. Reports follow file argument order and findings are ordered by line and code. If a file cannot be read, contains invalid UTF-8, or standard input is named twice, the command discards all buffered report output and exits 2.
- Successful commands exit 0. `get` or `bool` exits 1 when the key is absent, and `check` exits 1 when it reports findings. Usage, input, parsing, or invalid-boolean errors exit 2 and write a diagnostic to standard error.

## Development

Run these checks from the project directory:

```sh
gofmt -l .
go vet ./...
go test ./...
```

## License

MIT, see [LICENSE](LICENSE).

---

Built with BAKED
