package cmd

import (
	"errors"

	"github.com/Gulliby-development/cfgquill/internal/base"
	"github.com/Gulliby-development/cfgquill/internal/problem"
	"github.com/Gulliby-development/cfgquill/internal/value"
)

const UnquoteUsage = "unquote TEXT"

// RunUnquote transforms one value-text argument into its unquoted value.
func RunUnquote(args []string, io base.IO) (int, error) {
	parsed, err := base.Parse(args, nil)
	if err != nil {
		return 0, err
	}
	if err := base.Need(UnquoteUsage, parsed.Args, 1, 1); err != nil {
		return 0, err
	}

	result, fault := value.Unquote(parsed.Args[0])
	if fault != nil {
		return 0, errors.New("unquote: " + fault.Code + " " + problem.FaultText(fault.Code, fault.Arg))
	}

	io.Out(result)
	return 0, nil
}
