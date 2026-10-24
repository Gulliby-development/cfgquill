package cli

import (
	"github.com/Gulliby-development/cfgquill/internal/base"
	"github.com/Gulliby-development/cfgquill/internal/cmd"
	"github.com/Gulliby-development/cfgquill/internal/lint"
)

// Usage is the help text for cfgquill.
const Usage = "usage: cfgquill list FILE\n" +
	"usage: cfgquill sections FILE\n" +
	"usage: cfgquill get [--all] KEY FILE\n" +
	"usage: cfgquill bool KEY FILE\n" +
	"usage: cfgquill unquote TEXT\n" +
	"usage: cfgquill quote TEXT\n" +
	"usage: cfgquill check FILE...\n"

// Run dispatches argv to a command and reports command errors consistently.
func Run(argv []string, io base.IO) int {
	if len(argv) == 1 && argv[0] == "--help" {
		io.Out(Usage)
		return 0
	}
	if len(argv) == 0 {
		io.Err("cfgquill: missing command (list, sections, get, bool, unquote, quote or check)\n")
		return 2
	}

	var (
		status int
		err    error
	)
	switch argv[0] {
	case "list":
		status, err = cmd.RunList(argv[1:], io)
	case "sections":
		status, err = cmd.RunSections(argv[1:], io)
	case "get":
		status, err = cmd.RunGet(argv[1:], io)
	case "bool":
		status, err = cmd.RunBool(argv[1:], io)
	case "unquote":
		status, err = cmd.RunUnquote(argv[1:], io)
	case "quote":
		status, err = cmd.RunQuote(argv[1:], io)
	case "check":
		status, err = lint.Run(argv[1:], io)
	default:
		io.Err("cfgquill: unknown command " + base.Show(argv[0]) + "\n")
		return 2
	}

	if err != nil {
		io.Err("cfgquill: " + err.Error() + "\n")
		return 2
	}
	return status
}
