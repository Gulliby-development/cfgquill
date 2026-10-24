package invisible

import "github.com/Gulliby-development/cfgquill/internal/base"

// Mark describes an invisible character recognized by the checker.
type Mark struct {
	Rune rune
	Code string
	Name string
}

var markTable = []Mark{
	{Rune: '\u00a0', Code: "U+00A0", Name: "NBSP"},
	{Rune: '\u200b', Code: "U+200B", Name: "ZWSP"},
	{Rune: '\u0085', Code: "U+0085", Name: "NEL"},
	{Rune: '\u2028', Code: "U+2028", Name: "LS"},
	{Rune: '\u2029', Code: "U+2029", Name: "PS"},
	{Rune: '\u3000', Code: "U+3000", Name: "IDEOSP"},
	{Rune: '\ufeff', Code: "U+FEFF", Name: "BOM"},
}

// Marks returns the distinct recognized invisible characters in first-appearance order.
func Marks(value string) []Mark {
	found := make(map[rune]bool)
	var marks []Mark
	cursor := base.NewCursor(value)
	for r := cursor.Next(); r != base.EOF; r = cursor.Next() {
		if found[r] {
			continue
		}
		for _, mark := range markTable {
			if r == mark.Rune {
				found[r] = true
				marks = append(marks, mark)
				break
			}
		}
	}
	return marks
}

// Describe returns a mark's code and short name.
func (m Mark) Describe() string {
	return m.Code + " (" + m.Name + ")"
}
