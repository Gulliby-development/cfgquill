package lint

import (
	"strconv"

	"github.com/Gulliby-development/cfgquill/internal/base"
	"github.com/Gulliby-development/cfgquill/internal/config"
	"github.com/Gulliby-development/cfgquill/internal/invisible"
	"github.com/Gulliby-development/cfgquill/internal/problem"
)

type entryState struct {
	key      string
	hasValue bool
	value    string
}

// Findings returns semantic findings for a parsed configuration and its raw text.
func Findings(text string, f config.File) []problem.Problem {
	var found []problem.Problem
	first := make(map[entryState]int)

	for _, entry := range f.Entries {
		value := entry.Value
		if !entry.HasValue {
			value = ""
		}
		state := entryState{key: entry.Key, hasValue: entry.HasValue, value: value}
		if line, ok := first[state]; ok {
			found = append(found, problem.Problem{
				Line: entry.Line,
				Code: "CQ101",
				Text: "duplicate entry " + base.Show(entry.Key) + " (first at line " + strconv.Itoa(line) + ")",
			})
		} else {
			first[state] = entry.Line
		}

		if len(value) > 0 && (value[0] == ' ' || value[0] == '\t' || value[0] == '\r' || value[len(value)-1] == ' ' || value[len(value)-1] == '\t' || value[len(value)-1] == '\r') {
			found = append(found, problem.Problem{
				Line: entry.Line,
				Code: "CQ103",
				Text: "blank at the edge of value " + base.Show(value),
			})
		}

		for _, mark := range invisible.Marks(value) {
			found = append(found, problem.Problem{
				Line: entry.Line,
				Code: "CQ104",
				Text: "invisible character " + mark.Describe() + " in the value of " + base.Show(entry.Key),
			})
		}
	}

	for _, header := range f.Headers {
		if header.Count == 0 {
			found = append(found, problem.Problem{
				Line: header.Line,
				Code: "CQ102",
				Text: "empty section " + base.Show(header.Name),
			})
		}
	}

	for _, position := range invisible.BareCRs(text) {
		found = append(found, problem.Problem{
			Line: position.Line,
			Code: "CQ105",
			Text: "bare CR at column " + strconv.Itoa(position.Col),
		})
	}

	return problem.Sorted(found)
}
