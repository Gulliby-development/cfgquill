package problem

import (
	"sort"
	"strconv"

	"github.com/Gulliby-development/cfgquill/internal/base"
)

// Problem is a diagnostic attached to a source line.
type Problem struct {
	Line       int
	Code, Text string
}

// FromFault converts a syntax fault into a diagnostic.
func FromFault(f *base.Fault) Problem {
	return Problem{Line: f.Line, Code: f.Code, Text: FaultText(f.Code, f.Arg)}
}

// Sorted returns a stable, line-then-code ordered copy of found.
func Sorted(found []Problem) []Problem {
	result := make([]Problem, len(found))
	copy(result, found)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Line != result[j].Line {
			return result[i].Line < result[j].Line
		}
		return result[i].Code < result[j].Code
	})
	return result
}

// Render formats the diagnostic without adding a prefix or newline.
func (p Problem) Render(display string) string {
	return display + ":" + strconv.Itoa(p.Line) + ": " + p.Code + " " + p.Text
}
