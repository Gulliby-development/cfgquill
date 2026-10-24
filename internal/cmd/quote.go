package cmd

import (
	"github.com/Gulliby-development/cfgquill/internal/base"
	"github.com/Gulliby-development/cfgquill/internal/value"
)

// QuoteUsage describes the arguments accepted by the quote command.
const QuoteUsage = "quote TEXT"

// RunQuote validates and quotes its single text argument.
func RunQuote(args []string, io base.IO) (int, error) {
	parsed, err := base.Parse(args, nil)
	if err != nil {
		return 0, err
	}
	if err := base.Need(QuoteUsage, parsed.Args, 1, 1); err != nil {
		return 0, err
	}

	io.Out(value.Quote(parsed.Args[0]))
	return 0, nil
}
