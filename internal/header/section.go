package header

import (
	"strings"

	"github.com/Gulliby-development/cfgquill/internal/base"
)

// Parse reads a section header after its opening bracket and returns its stored name.
func Parse(c *base.Cursor) (string, *base.Fault) {
	var name strings.Builder

	for {
		r := c.Next()
		switch {
		case r == ']':
			if name.Len() == 0 {
				return "", c.Fault("CQ001", "")
			}
			return name.String(), nil
		case r == base.EOF:
			return "", c.Fault("CQ001", "")
		case base.IsAlnum(r) || r == '-' || r == '.':
			name.WriteString(base.Lower(string(r)))
		case r == ' ' || r == '\t' || r == '\r':
			return parseExtended(c, &name)
		default:
			return "", c.Fault("CQ001", "")
		}
	}
}

func parseExtended(c *base.Cursor, name *strings.Builder) (string, *base.Fault) {
	for {
		r := c.Next()
		if r == ' ' || r == '\t' || r == '\r' {
			continue
		}
		if r == base.EOF {
			return "", c.Fault("CQ001", "")
		}
		if r != '"' {
			return "", c.Fault("CQ001", "")
		}
		break
	}

	name.WriteByte('.')
	for {
		r := c.Next()
		switch r {
		case '\n':
			return "", c.Fault("CQ001", "")
		case base.EOF:
			return "", c.Fault("CQ001", "")
		case '\\':
			escaped := c.Next()
			if escaped == '\n' || escaped == base.EOF {
				return "", c.Fault("CQ001", "")
			}
			name.WriteRune(escaped)
		case '"':
			next := c.Next()
			if next == base.EOF {
				return "", c.Fault("CQ001", "")
			}
			if next != ']' {
				return "", c.Fault("CQ001", "")
			}
			return name.String(), nil
		default:
			name.WriteRune(r)
		}
	}
}
