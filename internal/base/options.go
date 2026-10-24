package base

import "errors"

// Parsed contains the recognized option flags and the remaining arguments.
type Parsed struct {
	Flags map[string]bool
	Args  []string
}

// Parse separates exact, recognized options from ordinary arguments.
func Parse(argv []string, specs []string) (Parsed, error) {
	known := make(map[string]struct{}, len(specs))
	for _, spec := range specs {
		known[spec] = struct{}{}
	}

	parsed := Parsed{
		Flags: make(map[string]bool),
		Args:  make([]string, 0, len(argv)),
	}
	for _, arg := range argv {
		if len(arg) >= 2 && arg[:2] == "--" {
			if _, ok := known[arg]; !ok {
				return Parsed{}, errors.New("unknown option " + Show(arg))
			}
			parsed.Flags[arg] = true
			continue
		}
		parsed.Args = append(parsed.Args, arg)
	}
	return parsed, nil
}

// Need verifies that args has a permitted number of ordinary arguments.
func Need(words string, args []string, min, max int) error {
	if len(args) < min || (max >= 0 && len(args) > max) {
		return errors.New("usage: " + words)
	}
	return nil
}
