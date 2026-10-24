package cli

import (
	"os"
	"strings"
	"testing"
)

func TestREADMEContract(t *testing.T) {
	readme, err := os.ReadFile("../../README.md")
	if err != nil {
		t.Fatal(err)
	}
	text := string(readme)

	required := []string{
		"## Install",
		"## Usage",
		"## Commands",
		"## Development",
		"cfgquill is an offline command-line tool",
		"people who read or edit Git configuration files",
		"`.gitconfig`, `.git/config`, and `.gitmodules`",
		"it does not discover repositories or modify configuration files",
		"go build -o cfgquill .",
		"go build -o cfgquill.exe .",
		"Run `./cfgquill` on macOS or Linux",
		"In PowerShell use `.\\cfgquill.exe`",
		"in Command Prompt use `cfgquill.exe`",
		"cfgquill list .gitconfig",
		"cfgquill sections .git/config",
		"cfgquill get --all remote.origin.fetch .gitconfig",
		"cfgquill bool core.filemode .git/config",
		"cfgquill unquote value",
		`cfgquill quote "two words"`,
		"cfgquill check .gitconfig .gitmodules",
		"cfgquill list - < .gitconfig",
		"strict UTF-8 configuration text",
		"space, TAB, LF, and CR",
		"`#` and `;` start comments outside double quotes",
		"usage: cfgquill list FILE",
		"usage: cfgquill sections FILE",
		"usage: cfgquill get [--all] KEY FILE",
		"usage: cfgquill bool KEY FILE",
		"usage: cfgquill unquote TEXT",
		"usage: cfgquill quote TEXT",
		"usage: cfgquill check FILE...",
		"Section and variable parts of a command-line key are ASCII case-insensitive",
		"a subsection is case-sensitive",
		"An entry without a value is true, while an explicitly empty value is false",
		"discards all buffered report output and exits 2",
		"gofmt -l .",
		"go vet ./...",
		"go test ./...",
	}
	for _, expected := range required {
		if !strings.Contains(text, expected) {
			t.Errorf("README.md does not contain %q", expected)
		}
	}
}
