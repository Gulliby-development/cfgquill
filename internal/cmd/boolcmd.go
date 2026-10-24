package cmd

import (
	"errors"
	"strconv"

	"github.com/Gulliby-development/cfgquill/internal/base"
	"github.com/Gulliby-development/cfgquill/internal/config"
	"github.com/Gulliby-development/cfgquill/internal/header"
	"github.com/Gulliby-development/cfgquill/internal/value"
)

// BoolUsage is the usage text for the bool command.
const BoolUsage = "bool KEY FILE"

// RunBool prints the boolean value of the last entry for a key.
func RunBool(args []string, io base.IO) (int, error) {
	parsed, err := base.Parse(args, nil)
	if err != nil {
		return 0, err
	}
	if err := base.Need(BoolUsage, parsed.Args, 2, 2); err != nil {
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

	entry := entries[len(entries)-1]
	result := true
	if entry.HasValue {
		boolean, valid := value.Bool(entry.Value)
		if !valid {
			return 0, errors.New(loaded.Display + ":" + strconv.Itoa(entry.Line) +
				": bad boolean " + base.Show(entry.Value) + " for " + canonical)
		}
		result = boolean
	}

	if result {
		io.Out("true\n")
	} else {
		io.Out("false\n")
	}
	return 0, nil
}
