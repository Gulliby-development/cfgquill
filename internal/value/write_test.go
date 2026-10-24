package value

import "testing"

func TestQuote(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"ordinary text", "plain text", "plain text"},
		{"leading ASCII space wraps", " leading", "\" leading\""},
		{"trailing ASCII space wraps", "trailing ", "\"trailing \""},
		{"both edge spaces wrap", " edge ", "\" edge \""},
		{"edge tab does not wrap", "\t", "\\t"},
		{"interior tab does not wrap", "a\tb", "a\\tb"},
		{"semicolon wraps", "a;b", "\"a;b\""},
		{"hash wraps", "a#b", "\"a#b\""},
		{"interior CR wraps and stays raw", "a\rb", "\"a\rb\""},
		{"trailing CR wraps", "a\r", "\"a\r\""},
		{"LF is escaped without wrapper", "a\nb", "a\\nb"},
		{"TAB is escaped without wrapper", "a\tb", "a\\tb"},
		{"quote is escaped", "a\"b", "a\\\"b"},
		{"backslash is escaped", "a\\b", "a\\\\b"},
		{"backspace stays raw", "a\bb", "a\bb"},
		{"NBSP is not an edge-space trigger", " x ", " x "},
		{"Unicode non-blanks stay raw", "a  　​\ufeffb", "a  　​\ufeffb"},
		{"astral rune stays raw", "a😀b", "a😀b"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Quote(test.in); got != test.want {
				t.Errorf("Quote(%q) = %q, want %q", test.in, got, test.want)
			}
		})
	}
}

func TestBool(t *testing.T) {
	tests := []struct {
		in   string
		want bool
		ok   bool
	}{
		{"", false, true},
		{"true", true, true},
		{"TRUE", true, true},
		{"TrUe", true, true},
		{"yes", true, true},
		{"YeS", true, true},
		{"on", true, true},
		{"ON", true, true},
		{"false", false, true},
		{"FALSE", false, true},
		{"no", false, true},
		{"No", false, true},
		{"off", false, true},
		{"OfF", false, true},
		{"0", false, true},
		{"000", false, true},
		{"000000000", false, true},
		{"1", true, true},
		{"8", true, true},
		{"000000001", true, true},
		{"123456789", true, true},
		{"0123456789", false, false},
		{"0000000000", false, false},
		{"-1", false, false},
		{"+1", false, false},
		{" true", false, false},
		{"true ", false, false},
		{"1.0", false, false},
		{"١", false, false},
		{" ", false, false},
		{"TRUE ", false, false},
	}
	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			got, ok := Bool(test.in)
			if got != test.want || ok != test.ok {
				t.Errorf("Bool(%q) = (%t, %t), want (%t, %t)", test.in, got, ok, test.want, test.ok)
			}
		})
	}
}
