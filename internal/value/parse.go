package value

import "github.com/Gulliby-development/cfgquill/internal/base"

// Parse reads a value from c through its terminating line feed.
func Parse(c *base.Cursor) (string, *base.Fault) {
	var value []rune
	var pending []rune
	quoted := false
	comment := false

	flush := func() {
		if len(value) != 0 {
			value = append(value, pending...)
		}
		pending = pending[:0]
	}

	for {
		r := c.Next()
		if r == base.EOF {
			if quoted {
				return string(value), c.Fault("CQ004", "")
			}
			return string(value), nil
		}

		if r == '\n' {
			if quoted {
				return string(value), c.Fault("CQ004", "")
			}
			return string(value), nil
		}
		if comment {
			continue
		}
		if !quoted {
			switch r {
			case ' ', '\t', '\r':
				pending = append(pending, r)
				continue
			case ';', '#':
				comment = true
				continue
			}
		}

		flush()
		switch r {
		case '\\':
			next := c.Next()
			if next == '\n' {
				continue
			}
			switch next {
			case 't':
				value = append(value, '\t')
			case 'b':
				value = append(value, '\b')
			case 'n':
				value = append(value, '\n')
			case '\\':
				value = append(value, '\\')
			case '"':
				value = append(value, '"')
			default:
				return string(value), c.Fault("CQ005", string([]rune{'\\', next}))
			}
		case '"':
			quoted = !quoted
		default:
			value = append(value, r)
		}
	}
}

// Unquote reads a value and rejects text beyond its terminating line feed.
func Unquote(text string) (string, *base.Fault) {
	cursor := base.NewCursor(text)
	value, fault := Parse(cursor)
	if fault != nil {
		return value, fault
	}
	if !cursor.AtEnd() {
		return value, cursor.Fault("CQ006", "")
	}
	return value, nil
}
