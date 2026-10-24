package cmd

import (
	"strconv"
	"strings"

	"github.com/Gulliby-development/cfgquill/internal/base"
	"github.com/Gulliby-development/cfgquill/internal/config"
	"github.com/Gulliby-development/cfgquill/internal/header"
)

const SectionsUsage = "sections FILE"

// RunSections lists section headers in a configuration file.
func RunSections(args []string, io base.IO) (int, error) {
	parsed, err := base.Parse(args, nil)
	if err != nil {
		return 0, err
	}
	if err := base.Need(SectionsUsage, parsed.Args, 1, 1); err != nil {
		return 0, err
	}

	loaded, err := config.Load(base.NewInputs(io), parsed.Args[0])
	if err != nil {
		return 0, err
	}
	if len(loaded.File.Headers) == 0 {
		return 0, nil
	}

	rows := make([][]string, 0, len(loaded.File.Headers))
	for _, section := range loaded.File.Headers {
		name, subsection, hasSubsection := header.Split(section.Name)
		if hasSubsection {
			subsection = base.Show(subsection)
		} else {
			subsection = "-"
		}
		rows = append(rows, []string{
			strconv.Itoa(section.Line),
			name,
			subsection,
			strconv.Itoa(section.Count),
		})
	}

	lines := base.Table([]string{"LINE", "SECTION", "SUBSECTION", "ENTRIES"}, rows)
	io.Out(strings.Join(lines, "\n") + "\n")
	return 0, nil
}
