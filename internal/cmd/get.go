package cmd

import (
	"errors"

	"github.com/Gulliby-development/cfgquill/internal/base"
	"github.com/Gulliby-development/cfgquill/internal/config"
	"github.com/Gulliby-development/cfgquill/internal/header"
)

// GetUsage is the usage text for the get command.
const GetUsage = "get [--all] KEY FILE"

// RunGet prints the last value for a key, or every value when --all is set.
func RunGet(args []string, io base.IO) (int, error) {
	parsed, err := base.Parse(args, []string{"--all"})
	if err != nil {
		return 0, err
	}
	if err := base.Need(GetUsage, parsed.Args, 2, 2); err != nil {
		return 0, err
	}

	canonical, ok := header.Canonical(parsed.Args[0])
	if !ok {
		return 0, errors.New("invalid key " + base.Show(parsed.Args[0]))
	}

	loaded, err := config.Load(base.NewInputs(io), parsed.Args[1])
	if err != nil {
		return 0, err
	}
	entries := loaded.File.Find(canonical)
	if len(entries) == 0 {
		return 1, nil
	}

	var output string
	if parsed.Flags["--all"] {
		for _, entry := range entries {
			output += entry.Value + "\n"
		}
	} else {
		output = entries[len(entries)-1].Value + "\n"
	}

	io.Out(output)
	return 0, nil
}
