package base

import "unicode/utf8"

// EOF is returned after the cursor has emitted its synthetic final newline.
const EOF rune = -1

// Fault describes a parse fault at a cursor position.
type Fault struct {
	Code string
	Line int
	Arg  string
}

// Error returns the fault code.
func (f *Fault) Error() string {
	return f.Code
}

// Cursor reads text one rune at a time, folding CRLF and adding a final LF.
type Cursor struct {
	text     string
	offset   int
	nextLine int
	lastLine int
	terminal bool
}

// NewCursor creates a cursor over text.
func NewCursor(text string) *Cursor {
	return &Cursor{text: text, nextLine: 1, lastLine: 1}
}

// Next returns the next rune, one synthetic LF after the original text, and then EOF.
func (c *Cursor) Next() rune {
	if c.offset < len(c.text) {
		r, size := utf8.DecodeRuneInString(c.text[c.offset:])
		c.offset += size
		if r == '\r' && c.offset < len(c.text) && c.text[c.offset] == '\n' {
			c.offset++
			r = '\n'
		}

		c.lastLine = c.nextLine
		if r == '\n' {
			c.nextLine++
		}
		return r
	}

	if !c.terminal {
		c.terminal = true
		c.lastLine = c.nextLine
		c.nextLine++
		return '\n'
	}

	return EOF
}

// Line returns the line of the most recently returned rune.
func (c *Cursor) Line() int {
	return c.lastLine
}

// AtEnd reports whether all original text has been consumed.
func (c *Cursor) AtEnd() bool {
	return c.offset >= len(c.text)
}

// Fault creates a fault at the current cursor line.
func (c *Cursor) Fault(code, arg string) *Fault {
	return &Fault{Code: code, Line: c.Line(), Arg: arg}
}

// IsSpace reports whether r is one of the four recognized whitespace runes.
func IsSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

// IsAlpha reports whether r is an ASCII letter.
func IsAlpha(r rune) bool {
	return r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z'
}

// IsAlnum reports whether r is an ASCII letter or digit.
func IsAlnum(r rune) bool {
	return IsAlpha(r) || r >= '0' && r <= '9'
}

// IsKeyChar reports whether r is an ASCII letter, digit, or hyphen.
func IsKeyChar(r rune) bool {
	return IsAlnum(r) || r == '-'
}

// Lower converts ASCII uppercase letters in s to lowercase.
func Lower(s string) string {
	lowered := []byte(s)
	for i, b := range lowered {
		if b >= 'A' && b <= 'Z' {
			lowered[i] = b + ('a' - 'A')
		}
	}
	return string(lowered)
}
