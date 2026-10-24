package config

import (
	"strings"

	"github.com/Gulliby-development/cfgquill/internal/base"
	"github.com/Gulliby-development/cfgquill/internal/header"
	"github.com/Gulliby-development/cfgquill/internal/value"
)

// Parse reads configuration text and returns its completed prefix and first fault.
func Parse(text string) (File, *base.Fault) {
	file := File{}
	cursor := base.NewCursor(text)
	current := cursor.Next()

	if current == '\ufeff' {
		current = cursor.Next()
	}

	name := ""
	currentHeader := -1

	for {
		switch {
		case base.IsSpace(current):
			current = cursor.Next()
		case current == base.EOF:
			return file, nil
		case current == '#' || current == ';':
			for current != '\n' && current != base.EOF {
				current = cursor.Next()
			}
			if current == '\n' {
				current = cursor.Next()
			}
		case current == '[':
			line := cursor.Line()
			parsedName, fault := header.Parse(cursor)
			if fault != nil {
				return file, fault
			}
			name = parsedName
			file.Headers = append(file.Headers, Header{Line: line, Name: name})
			currentHeader = len(file.Headers) - 1
			current = cursor.Next()
		case !base.IsAlpha(current):
			return file, cursor.Fault("CQ002", "")
		default:
			entry, fault := parseEntry(cursor, current, name)
			if fault != nil {
				return file, fault
			}
			file.Entries = append(file.Entries, entry)
			if currentHeader >= 0 {
				file.Headers[currentHeader].Count++
			}
			current = cursor.Next()
		}
	}
}

func parseEntry(cursor *base.Cursor, first rune, name string) (Entry, *base.Fault) {
	line := cursor.Line()
	var variable strings.Builder
	variable.WriteRune(first)

	current := cursor.Next()
	for base.IsKeyChar(current) {
		variable.WriteRune(current)
		current = cursor.Next()
	}

	for current == ' ' || current == '\t' {
		current = cursor.Next()
	}

	entryVar := base.Lower(variable.String())
	key := entryVar
	if name != "" {
		key = name + "." + entryVar
	}
	entry := Entry{
		Line: line,
		Name: name,
		Var:  entryVar,
		Key:  key,
	}

	switch current {
	case '\n':
		return entry, nil
	case '=':
		valueText, fault := value.Parse(cursor)
		if fault != nil {
			return Entry{}, fault
		}
		entry.Value = valueText
		entry.HasValue = true
		return entry, nil
	default:
		return Entry{}, cursor.Fault("CQ003", "")
	}
}
