package value

import "strings"

// Quote returns v in the escaped form Git writes after a key assignment.
func Quote(v string) string {
	quoted := len(v) > 0 && (v[0] == ' ' || v[len(v)-1] == ' ')
	for _, r := range v {
		if r == ';' || r == '#' || r == '\r' {
			quoted = true
		}
	}

	var out strings.Builder
	if quoted {
		out.WriteByte('"')
	}
	for _, r := range v {
		switch r {
		case '\n':
			out.WriteString(`\n`)
		case '\t':
			out.WriteString(`\t`)
		case '"':
			out.WriteString(`\"`)
		case '\\':
			out.WriteString(`\\`)
		default:
			out.WriteRune(r)
		}
	}
	if quoted {
		out.WriteByte('"')
	}
	return out.String()
}
