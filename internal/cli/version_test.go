package cli

import (
	"errors"
	"os"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

type versionTestIO struct {
	stdout []string
	stderr []string
}

func (io *versionTestIO) ReadStdin() ([]byte, error) {
	return nil, nil
}

func (io *versionTestIO) ReadFile(string) ([]byte, error) {
	return nil, errors.New("file not found")
}

func (io *versionTestIO) Out(text string) {
	io.stdout = append(io.stdout, text)
}

func (io *versionTestIO) Err(text string) {
	io.stderr = append(io.stderr, text)
}

func TestVersionSingleton(t *testing.T) {
	io := &versionTestIO{}
	if status := Run([]string{"--version"}, io); status != 0 {
		t.Fatalf("status = %d, want 0", status)
	}
	want := "cfgquill built with " + runtime.Version() + "\n"
	if len(io.stdout) != 1 || io.stdout[0] != want {
		t.Fatalf("stdout writes = %q, want exactly [%q]", io.stdout, want)
	}
	matches, err := regexp.MatchString("^cfgquill built with go1\\.[0-9]+(\\.[0-9]+)?\\n$", io.stdout[0])
	if err != nil {
		t.Fatal(err)
	}
	if !matches {
		t.Errorf("stdout = %q, want a Go release version line", io.stdout[0])
	}
	if len(io.stderr) != 0 {
		t.Errorf("stderr writes = %q, want none", io.stderr)
	}
}

func TestVersionInvalidArgumentLists(t *testing.T) {
	beforeCommand := &versionTestIO{}
	if status := Run([]string{"--version", "list"}, beforeCommand); status != 2 {
		t.Errorf("status for --version list = %d, want 2", status)
	}
	if len(beforeCommand.stdout) != 0 {
		t.Errorf("stdout for --version list = %q, want none", beforeCommand.stdout)
	}
	if len(beforeCommand.stderr) != 1 || beforeCommand.stderr[0] != "cfgquill: unknown command \"--version\"\n" {
		t.Errorf("stderr for --version list = %q, want exact unknown-command diagnostic", beforeCommand.stderr)
	}

	repeatedVersion := &versionTestIO{}
	if status := Run([]string{"--version", "--version"}, repeatedVersion); status != 2 {
		t.Errorf("status for repeated --version = %d, want 2", status)
	}
	if len(repeatedVersion.stdout) != 0 {
		t.Errorf("stdout for repeated --version = %q, want none", repeatedVersion.stdout)
	}
	if len(repeatedVersion.stderr) != 1 || repeatedVersion.stderr[0] != "cfgquill: unknown command \"--version\"\n" {
		t.Errorf("stderr for repeated --version = %q, want exact unknown-command diagnostic", repeatedVersion.stderr)
	}

	afterCommand := &versionTestIO{}
	if status := Run([]string{"list", "--version"}, afterCommand); status != 2 {
		t.Errorf("status for list --version = %d, want 2", status)
	}
	if len(afterCommand.stdout) != 0 {
		t.Errorf("stdout for list --version = %q, want none", afterCommand.stdout)
	}
	if len(afterCommand.stderr) != 1 || afterCommand.stderr[0] != "cfgquill: unknown option \"--version\"\n" {
		t.Errorf("stderr for list --version = %q, want exact unknown-option diagnostic", afterCommand.stderr)
	}

	bareVersion := &versionTestIO{}
	if status := Run([]string{"version"}, bareVersion); status != 2 {
		t.Errorf("status for version = %d, want 2", status)
	}
	if len(bareVersion.stdout) != 0 {
		t.Errorf("stdout for version = %q, want none", bareVersion.stdout)
	}
	if len(bareVersion.stderr) != 1 || bareVersion.stderr[0] != "cfgquill: unknown command \"version\"\n" {
		t.Errorf("stderr for version = %q, want exact unknown-command diagnostic", bareVersion.stderr)
	}
}

func TestVersionModuleContract(t *testing.T) {
	module, err := os.ReadFile("../../go.mod")
	if err != nil {
		t.Fatal(err)
	}
	moduleText := string(module)
	want := "module github.com/Gulliby-development/cfgquill\n\ngo 1.27\n\ntoolchain go1.27.0\n"
	if moduleText != want {
		t.Fatalf("../../go.mod = %q, want %q", moduleText, want)
	}
	if !strings.Contains(moduleText, "\ngo 1.27\n") {
		t.Errorf("../../go.mod does not contain whole line %q", "go 1.27")
	}
	if !strings.Contains(moduleText, "\ntoolchain go1.27.0\n") {
		t.Errorf("../../go.mod does not contain whole line %q", "toolchain go1.27.0")
	}
}

func TestVersionREADMEContract(t *testing.T) {
	readme, err := os.ReadFile("../../README.md")
	if err != nil {
		t.Fatal(err)
	}
	text := string(readme)
	if !strings.Contains(text, "`cfgquill --version` prints one line") {
		t.Errorf("../../README.md does not contain %q", "`cfgquill --version` prints one line")
	}
	if !strings.Contains(text, "cfgquill built with go1.27.0") {
		t.Errorf("../../README.md does not contain %q", "cfgquill built with go1.27.0")
	}
	if !strings.Contains(text, "`toolchain go1.27.0`") {
		t.Errorf("../../README.md does not contain %q", "`toolchain go1.27.0`")
	}
}
