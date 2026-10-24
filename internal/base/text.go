package base

import "strconv"

// Plain returns s with every code point outside printable ASCII replaced by '?'.
func Plain(s string) string {
	var result []byte
	for _, r := range s {
		if r < ' ' || r > '~' {
			result = append(result, '?')
			continue
		}
		result = append(result, byte(r))
	}
	return string(result)
}

// Show returns Plain(s) surrounded by double quotes.
func Show(s string) string {
	return "\"" + Plain(s) + "\""
}

// Plural formats n and word, adding an s unless n is exactly one.
func Plural(n int, word string) string {
	suffix := "s"
	if n == 1 {
		suffix = ""
	}
	return strconv.Itoa(n) + " " + word + suffix
}

// Table renders a left-aligned table without trailing spaces.
func Table(header []string, rows [][]string) []string {
	columnCount := len(header)
	for _, row := range rows {
		if len(row) > columnCount {
			columnCount = len(row)
		}
	}

	widths := make([]int, columnCount)
	for column, cell := range header {
		widths[column] = runeWidth(cell)
	}
	for _, row := range rows {
		for column, cell := range row {
			width := runeWidth(cell)
			if width > widths[column] {
				widths[column] = width
			}
		}
	}

	lines := make([]string, 0, len(rows)+1)
	lines = append(lines, tableLine(header, widths))
	for _, row := range rows {
		lines = append(lines, tableLine(row, widths))
	}
	return lines
}

func runeWidth(s string) int {
	width := 0
	for range s {
		width++
	}
	return width
}

func tableLine(cells []string, widths []int) string {
	line := make([]byte, 0)
	for column, width := range widths {
		if column > 0 {
			line = append(line, ' ', ' ')
		}
		cell := ""
		if column < len(cells) {
			cell = cells[column]
		}
		line = append(line, cell...)
		padding := width - runeWidth(cell)
		for range padding {
			line = append(line, ' ')
		}
	}
	for len(line) > 0 && line[len(line)-1] == ' ' {
		line = line[:len(line)-1]
	}
	return string(line)
}
