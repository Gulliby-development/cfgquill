package cmd

import (
	"strconv"
	"strings"

	"github.com/Gulliby-development/cfgquill/internal/base"
	"github.com/Gulliby-development/cfgquill/internal/config"
)

const ListUsage = "list FILE"

// RunList lists entries in a configuration file.
func RunList(args []string, io base.IO) (int, error) {
	parsed, err := base.Parse(args, nil)
	if err != nil {
		return 0, err
	}
	if err := base.Need(ListUsage, parsed.Args, 1, 1); err != nil {
		return 0, err
	}

	loaded, err := config.Load(base.NewInputs(io), parsed.Args[0])
	if err != nil {
		return 0, err
	}
	if len(loaded.File.Entries) == 0 {
		return 0, nil
	}

	rows := make([][]string, 0, len(loaded.File.Entries))
	for _, entry := range loaded.File.Entries {
		value := "-"
		if entry.HasValue {
			value = base.Show(entry.Value)
		}
		rows = append(rows, []string{
			strconv.Itoa(entry.Line),
			entry.Key,
			value,
		})
	}

	lines := base.Table([]string{"LINE", "KEY", "VALUE"}, rows)
	io.Out(strings.Join(lines, "\n") + "\n")
	return 0, nil
}
