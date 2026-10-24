package header

// BEGIN acceptance rows imports: host-written from the audited spec table. Edit acceptanceActual, not this block.
import (
	acceptancejson "encoding/json"
	acceptancereflect "reflect"
	acceptancetesting "testing"
)

// END acceptance rows imports

import (
	"testing"

	"github.com/Gulliby-development/cfgquill/internal/base"
)

// BEGIN acceptance rows: host-written from the audited spec table. Edit acceptanceActual, not this block.
func TestAcceptanceRowsHeader(t *acceptancetesting.T) {
	var rows []struct {
		ID       string `json:"id"`
		Input    any    `json:"input"`
		Expected any    `json:"expected"`
	}
	if err := acceptancejson.Unmarshal([]byte("[{\"id\":\"AC-003:1\",\"input\":\"A.B\",\"expected\":[\"a.b\",true]},{\"id\":\"AC-003:2\",\"input\":\"a.S.b\",\"expected\":[\"a.S.b\",true]},{\"id\":\"AC-003:3\",\"input\":\"A.S.B\",\"expected\":[\"a.S.b\",true]},{\"id\":\"AC-003:4\",\"input\":\"a.s.b\",\"expected\":[\"a.s.b\",true]},{\"id\":\"AC-003:5\",\"input\":\"a..c\",\"expected\":[\"a..c\",true]},{\"id\":\"AC-003:6\",\"input\":\"a.x.y.B\",\"expected\":[\"a.x.y.b\",true]},{\"id\":\"AC-003:7\",\"input\":\"1-.B2-C\",\"expected\":[\"1-.b2-c\",true]},{\"id\":\"AC-003:8\",\"input\":\"a.😀 .B\",\"expected\":[\"a.😀 .b\",true]},{\"id\":\"AC-003:9\",\"input\":\"a.x\\ry\\tz.B\",\"expected\":[\"a.x\\ry\\tz.b\",true]}]"), &rows); err != nil {
		t.Fatal(err)
	}
	for _, row := range rows {
		row := row
		t.Run(row.ID, func(t *acceptancetesting.T) {
			got, err := acceptanceActual(row.Input)
			if err != nil {
				t.Fatal(err)
			}
			text, err := acceptancejson.Marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			var value any
			if err := acceptancejson.Unmarshal(text, &value); err != nil {
				t.Fatal(err)
			}
			if !acceptancereflect.DeepEqual(value, row.Expected) {
				want, _ := acceptancejson.Marshal(row.Expected)
				t.Fatalf("got %s, want %s", text, want)
			}
		})
	}
}

// END acceptance rows

func acceptanceActual(input any) (any, error) {
	key, _ := input.(string)
	canonical, valid := Canonical(key)
	return []any{canonical, valid}, nil
}

func TestParseValidHeaders(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
		next  rune
	}{
		{"old form", "[Core]", "core", 0},
		{"extended form preserves case", "[remote \"Origin\"]", "remote.Origin", 0},
		{"old form lowercases subsection", "[remote.Origin]", "remote.origin", 0},
		{"extended subsection after dotted prefix", "[a.b \"C\"]", "a.b.C", 0},
		{"empty subsection", "[a \"\"]", "a.", 0},
		{"digits and hyphen", "[1-.Z]", "1-.z", 0},
		{"extended whitespace", "[a \t\r\"B\"]", "a.B", 0},
		{"empty section prefix", "[ \"x\"]", ".x", 0},
		{"empty prefix and subsection", "[ \"\"]", ".", 0},
		{"unknown and escaped characters", "[a \"x\\q\\\"\\\\y\"]", "a.xq\"\\y", 0},
		{"astral subsection", "[a \"😀\"]", "a.😀", 0},
		{"NBSP subsection", "[a \"\u00a0\"]", "a.\u00a0", 0},
		{"NEL subsection", "[a \"\u0085\"]", "a.\u0085", 0},
		{"line separator subsection", "[a \"\u2028\"]", "a.\u2028", 0},
		{"paragraph separator subsection", "[a \"\u2029\"]", "a.\u2029", 0},
		{"bare CR subsection", "[a \"x\ry\"]", "a.x\ry", 0},
		{"leaves text after closing bracket", "[s] k = v", "s", ' '},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := base.NewCursor(test.input)
			if got := c.Next(); got != '[' {
				t.Fatalf("opening rune = %q, want '['", got)
			}
			got, fault := Parse(c)
			if fault != nil {
				t.Fatalf("Parse() fault = %q, want nil", fault.Code)
			}
			if got != test.want {
				t.Fatalf("Parse() = %q, want %q", got, test.want)
			}
			if test.next != 0 {
				if got := c.Next(); got != test.next {
					t.Fatalf("next rune = %q, want %q", got, test.next)
				}
			}
		})
	}
}

func TestParseMalformedHeaders(t *testing.T) {
	tests := []struct {
		name  string
		input string
		line  int
	}{
		{"empty old form", "[]", 1},
		{"underscore", "[a_]", 1},
		{"non-ASCII old-form character", "[é]", 1},
		{"subsection quote required", "[a ]", 1},
		{"no text after closing quote", "[a \"x\" ]", 1},
		{"no character after subsection quote", "[a \"x\"q]", 1},
		{"synthetic final LF", "[a", 1},
		{"unterminated subsection at end", "[a \"x", 1},
		{"LF in old form", "[a\n", 1},
		{"CRLF in old form", "[a\r\n", 1},
		{"LF after subsection backslash", "[a \"x\\\n", 1},
		{"physical trailing subsection backslash", "[a \"x\\", 1},
		{"NBSP outside subsection", "[a\u00a0\"x\"]", 1},
		{"VT outside subsection", "[a\v\"x\"]", 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := base.NewCursor(test.input)
			if got := c.Next(); got != '[' {
				t.Fatalf("opening rune = %q, want '['", got)
			}
			_, fault := Parse(c)
			if fault == nil {
				t.Fatal("Parse() fault = nil, want CQ001")
			}
			if fault.Code != "CQ001" {
				t.Fatalf("fault code = %q, want CQ001", fault.Code)
			}
			if fault.Line != test.line {
				t.Fatalf("fault line = %d, want %d", fault.Line, test.line)
			}
		})
	}
}

func TestParseMalformedHeaderLineAndExhaustedCursor(t *testing.T) {
	c := base.NewCursor("# ok\n[a \"b\\\"]\n")
	for _, want := range "# ok\n[" {
		if got := c.Next(); got != want {
			t.Fatalf("prefix rune = %q, want %q", got, want)
		}
	}
	_, fault := Parse(c)
	if fault == nil {
		t.Fatal("Parse() fault = nil, want CQ001")
	}
	if fault.Code != "CQ001" || fault.Line != 2 {
		t.Fatalf("fault = (%q, line %d), want (CQ001, line 2)", fault.Code, fault.Line)
	}

	exhausted := base.NewCursor("")
	exhausted.Next()
	if got := exhausted.Next(); got != base.EOF {
		t.Fatalf("second rune from empty cursor = %q, want EOF", got)
	}
	_, fault = Parse(exhausted)
	if fault == nil {
		t.Fatal("Parse() on exhausted cursor fault = nil, want CQ001")
	}
	if fault.Code != "CQ001" || fault.Line != 1 {
		t.Fatalf("exhausted cursor fault = (%q, line %d), want (CQ001, line 1)", fault.Code, fault.Line)
	}
}

func TestCanonicalInvalidKeys(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{"empty", ""},
		{"no dot", "b"},
		{"empty section", ".b"},
		{"empty variable", "a."},
		{"underscore in section", "a_b.c"},
		{"variable starts with digit", "a.1b"},
		{"underscore in variable", "a.b_c"},
		{"non-ASCII section", "é.b"},
		{"non-ASCII variable", "a.é"},
		{"LF in subsection", "a.x\ny.b"},
		{"space in section", "a .b"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, valid := Canonical(test.key)
			if valid {
				t.Fatalf("Canonical(%q) validity = true, want false", test.key)
			}
		})
	}
}

func TestSplitStoredNames(t *testing.T) {
	tests := []struct {
		label   string
		input   string
		section string
		sub     string
		hasSub  bool
	}{
		{"plain name", "core", "core", "", false},
		{"empty name", "", "", "", false},
		{"empty subsection", "a.", "a", "", true},
		{"dotted subsection", "a.b.C", "a", "b.C", true},
		{"case preserved", "remote.Origin", "remote", "Origin", true},
		{"empty section", ".x", "", "x", true},
	}

	for _, test := range tests {
		section, sub, hasSub := Split(test.input)
		if section != test.section || sub != test.sub || hasSub != test.hasSub {
			t.Fatalf("Split(%q) = (%q, %q, %t), want (%q, %q, %t)",
				test.input, section, sub, hasSub, test.section, test.sub, test.hasSub)
		}
	}
}
