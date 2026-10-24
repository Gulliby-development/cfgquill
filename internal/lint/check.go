package lint

import (
	"strings"

	"github.com/Gulliby-development/cfgquill/internal/base"
	"github.com/Gulliby-development/cfgquill/internal/config"
	"github.com/Gulliby-development/cfgquill/internal/problem"
)

// Usage is the command's usage text.
const Usage = "check FILE..."

// Run checks each named configuration file and emits one buffered report.
func Run(args []string, io base.IO) (int, error) {
	parsed, err := base.Parse(args, nil)
	if err != nil {
		return 0, err
	}
	if err := base.Need(Usage, parsed.Args, 1, -1); err != nil {
		return 0, err
	}

	inputs := base.NewInputs(io)
	var report strings.Builder
	total := 0

	for _, arg := range parsed.Args {
		source, err := inputs.Read(arg)
		if err != nil {
			return 0, err
		}

		file, fault := config.Parse(source.Text)
		var found []problem.Problem
		if fault != nil {
			found = []problem.Problem{problem.FromFault(fault)}
		} else {
			found = Findings(source.Text, file)
		}

		for _, item := range found {
			report.WriteString(item.Render(source.Display))
			report.WriteByte('\n')
			total++
		}
	}

	status := 0
	if total == 0 {
		report.WriteString("ok\n")
	} else {
		status = 1
		report.WriteString(base.Plural(total, "problem"))
		report.WriteByte('\n')
	}

	io.Out(report.String())
	return status, nil
}
