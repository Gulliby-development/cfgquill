package value

import (
	"testing"

	"github.com/Gulliby-development/cfgquill/internal/base"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		want     string
		code     string
		line     int
		arg      string
		next     rune
		nextLine int
	}{
		{"empty", "", "", "", 0, "", base.EOF, 1},
		{"only pending blanks", "  \t\r", "", "", 0, "", base.EOF, 1},
		{"leading blanks dropped", "  x", "x", "", 0, "", base.EOF, 1},
		{"interior blanks retained", "x \ty", "x \ty", "", 0, "", base.EOF, 1},
		{"hash comment", "# ignored", "", "", 0, "", base.EOF, 1},
		{"semicolon comment ignores escapes", "; ignored \\q", "", "", 0, "", base.EOF, 1},
		{"comment after value", "x  # ignored \\q", "x", "", 0, "", base.EOF, 1},
		{"quoted blanks", "\"a b\"", "a b", "", 0, "", base.EOF, 1},
		{"quoted edge blanks", "\"  \"", "  ", "", 0, "", base.EOF, 1},
		{"quoted tab", "\"\\t\"", "\t", "", 0, "", base.EOF, 1},
		{"empty quotes", "\"\"", "", "", 0, "", base.EOF, 1},
		{"empty quotes drop following leading blank", "\"\" a", "a", "", 0, "", base.EOF, 1},
		{"empty quotes retain preceding blank", "a \"\"", "a ", "", 0, "", base.EOF, 1},
		{"empty quotes retain blanks around text", "a \"\" b", "a  b", "", 0, "", base.EOF, 1},
		{"all allowed escapes", "a\\tb\\bc\\nd\\\\e\\\"f", "a\tb\bc\nd\\e\"f", "", 0, "", base.EOF, 1},
		{"escaped quote does not toggle", "\"a\\\"b\"", "a\"b", "", 0, "", base.EOF, 1},
		{"unknown escape", "a\\q", "a", "CQ005", 1, "\\q", 0, 0},
		{"escaped bare CR is unknown", "a\\\rb", "a", "CQ005", 1, "\\\r", 0, 0},
		{"escaped NBSP is unknown", "a\\ b", "a", "CQ005", 1, "\\ ", 0, 0},
		{"unterminated quote at end", "\"a", "a", "CQ004", 1, "", 0, 0},
		{"unterminated quote at physical LF", "\"a\nx", "a", "CQ004", 1, "", 0, 0},
		{"ordinary LF ends the value", "a\nnext", "a", "", 0, "", 'n', 2},
		{"continuation outside quotes", "a\\\nb", "ab", "", 0, "", base.EOF, 2},
		{"continuation keeps following interior blanks", "a\\\n  b", "a  b", "", 0, "", base.EOF, 2},
		{"continuation inside quotes", "\"a\\\n  b\"", "a  b", "", 0, "", base.EOF, 2},
		{"terminal backslash continues to EOF", "a\\", "a", "", 0, "", base.EOF, 1},
		{"terminal backslash in quote reports quote fault", "\"a\\", "a", "CQ004", 1, "", 0, 0},
		{"CRLF is a line ending", "x\r\ny", "x", "", 0, "", 'y', 2},
		{"bare CR inside is retained", "x\ry", "x\ry", "", 0, "", base.EOF, 1},
		{"bare CR at edge is dropped", "x\r", "x", "", 0, "", base.EOF, 1},
		{"escaped newline adds a newline", "a\\nb", "a\nb", "", 0, "", base.EOF, 1},
		{"escaped backslash", "a\\\\b", "a\\b", "", 0, "", base.EOF, 1},
		{"Unicode non-blanks are ordinary", "x   　​\ufeffy", "x   　​\ufeffy", "", 0, "", base.EOF, 1},
		{"astral rune is preserved", "😀x", "😀x", "", 0, "", base.EOF, 1},
		{"continued fault uses next line", "a\\\n\\q", "a", "CQ005", 2, "\\q", 0, 0},
		{"pending blanks disappear before comment", "x \t; tail", "x", "", 0, "", base.EOF, 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cursor := base.NewCursor(test.input)
			got, fault := Parse(cursor)
			if got != test.want {
				t.Errorf("Parse value = %q, want %q", got, test.want)
			}
			if test.code == "" {
				if fault != nil {
					t.Fatalf("Parse fault = %#v, want none", fault)
				}
				if next := cursor.Next(); next != test.next {
					t.Errorf("next rune = %U, want %U", next, test.next)
				}
				if line := cursor.Line(); line != test.nextLine {
					t.Errorf("line after next rune = %d, want %d", line, test.nextLine)
				}
				return
			}
			if fault == nil {
				t.Fatalf("Parse fault = nil, want %s", test.code)
			}
			if fault.Code != test.code || fault.Line != test.line || fault.Arg != test.arg {
				t.Errorf("Parse fault = %#v, want code %s, line %d, arg %q", fault, test.code, test.line, test.arg)
			}
		})
	}
}

func TestUnquote(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
		code  string
		line  int
		arg   string
	}{
		{"empty input", "", "", "", 0, ""},
		{"only blanks", "  ", "", "", 0, ""},
		{"trailing LF is allowed", "a\n", "a", "", 0, ""},
		{"text after LF", "a\nb", "a", "CQ006", 1, ""},
		{"text after folded CRLF", "a\r\nb", "a", "CQ006", 1, ""},
		{"continued LF is part of value", "a\\\nb", "ab", "", 0, ""},
		{"continued quote is part of value", "\"a\\\n b\"", "a b", "", 0, ""},
		{"second LF is trailing text", "a\n\n", "a", "CQ006", 1, ""},
		{"unterminated quote is not overridden", "\"a\nb\"", "a", "CQ004", 1, ""},
		{"unknown escape is not overridden", "a\\q\nrest", "a", "CQ005", 1, "\\q"},
		{"terminal backslash continuation", "a\\", "a", "", 0, ""},
		{"terminal backslash leaves quote open", "\"a\\", "a", "CQ004", 1, ""},
		{"bare CR is data except at edge", "a\rb", "a\rb", "", 0, ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, fault := Unquote(test.input)
			if got != test.want {
				t.Errorf("Unquote value = %q, want %q", got, test.want)
			}
			if test.code == "" {
				if fault != nil {
					t.Errorf("Unquote fault = %#v, want none", fault)
				}
				return
			}
			if fault == nil {
				t.Fatalf("Unquote fault = nil, want %s", test.code)
			}
			if fault.Code != test.code || fault.Line != test.line || fault.Arg != test.arg {
				t.Errorf("Unquote fault = %#v, want code %s, line %d, arg %q", fault, test.code, test.line, test.arg)
			}
		})
	}
}
