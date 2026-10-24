package invisible

import "github.com/Gulliby-development/cfgquill/internal/base"

// Pos is a one-based line and rune-column position.
type Pos struct {
	Line int
	Col  int
}

// BareCRs returns the positions of carriage returns that are not part of CRLF.
func BareCRs(text string) []Pos {
	var positions []Pos
	line, col := 1, 0
	cursor := base.NewCursor(text)
	for r := cursor.Next(); r != base.EOF; r = cursor.Next() {
		if r == '\r' {
			positions = append(positions, Pos{Line: line, Col: col + 1})
		}
		if r == '\n' {
			line++
			col = 0
		} else {
			col++
		}
	}
	return positions
}
