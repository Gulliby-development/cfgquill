package invisible

import (
	"reflect"
	"testing"
)

func TestMarks(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  []Mark
	}{
		{name: "empty", value: "", want: nil},
		{name: "ordinary text", value: "plain text", want: nil},
		{name: "excluded whitespace", value: " \t\n\r\v\f", want: nil},
		{
			name:  "all marks in table order",
			value: "\u00a0\u200b\u0085\u2028\u2029\u3000\ufeff",
			want: []Mark{
				{Rune: '\u00a0', Code: "U+00A0", Name: "NBSP"},
				{Rune: '\u200b', Code: "U+200B", Name: "ZWSP"},
				{Rune: '\u0085', Code: "U+0085", Name: "NEL"},
				{Rune: '\u2028', Code: "U+2028", Name: "LS"},
				{Rune: '\u2029', Code: "U+2029", Name: "PS"},
				{Rune: '\u3000', Code: "U+3000", Name: "IDEOSP"},
				{Rune: '\ufeff', Code: "U+FEFF", Name: "BOM"},
			},
		},
		{
			name:  "duplicates retain first appearance only",
			value: "\ufeff\u00a0\ufeff\u0085\u00a0",
			want: []Mark{
				{Rune: '\ufeff', Code: "U+FEFF", Name: "BOM"},
				{Rune: '\u00a0', Code: "U+00A0", Name: "NBSP"},
				{Rune: '\u0085', Code: "U+0085", Name: "NEL"},
			},
		},
		{
			name:  "repeated mark",
			value: "\u200b\u200b",
			want:  []Mark{{Rune: '\u200b', Code: "U+200B", Name: "ZWSP"}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Marks(test.value); !reflect.DeepEqual(got, test.want) {
				t.Errorf("Marks(%q) = %#v, want %#v", test.value, got, test.want)
			}
		})
	}
}

func TestMarkDescriptions(t *testing.T) {
	tests := []struct {
		mark Mark
		want string
	}{
		{mark: Mark{Rune: '\u00a0', Code: "U+00A0", Name: "NBSP"}, want: "U+00A0 (NBSP)"},
		{mark: Mark{Rune: '\u200b', Code: "U+200B", Name: "ZWSP"}, want: "U+200B (ZWSP)"},
		{mark: Mark{Rune: '\u0085', Code: "U+0085", Name: "NEL"}, want: "U+0085 (NEL)"},
		{mark: Mark{Rune: '\u2028', Code: "U+2028", Name: "LS"}, want: "U+2028 (LS)"},
		{mark: Mark{Rune: '\u2029', Code: "U+2029", Name: "PS"}, want: "U+2029 (PS)"},
		{mark: Mark{Rune: '\u3000', Code: "U+3000", Name: "IDEOSP"}, want: "U+3000 (IDEOSP)"},
		{mark: Mark{Rune: '\ufeff', Code: "U+FEFF", Name: "BOM"}, want: "U+FEFF (BOM)"},
	}
	for _, test := range tests {
		if got := test.mark.Describe(); got != test.want {
			t.Errorf("Describe() = %q, want %q", got, test.want)
		}
	}
}

func TestMarksCallsAreIndependent(t *testing.T) {
	first := Marks("\u00a0")
	second := Marks("\u00a0")
	want := []Mark{{Rune: '\u00a0', Code: "U+00A0", Name: "NBSP"}}
	if !reflect.DeepEqual(first, want) {
		t.Errorf("first Marks call = %#v, want %#v", first, want)
	}
	if !reflect.DeepEqual(second, want) {
		t.Errorf("second Marks call = %#v, want %#v", second, want)
	}
}

func TestBareCRs(t *testing.T) {
	tests := []struct {
		name string
		text string
		want []Pos
	}{
		{name: "empty", text: "", want: nil},
		{name: "ordinary text", text: "plain text", want: nil},
		{name: "CRLF is not bare", text: "\r\n", want: nil},
		{name: "bare CR at beginning and end", text: "\r", want: []Pos{{Line: 1, Col: 1}}},
		{name: "consecutive CRs before LF", text: "\r\r\n", want: []Pos{{Line: 1, Col: 1}}},
		{
			name: "multiple lines, tabs and trailing CR",
			text: "ab\r\rx\n\t\r",
			want: []Pos{{Line: 1, Col: 3}, {Line: 1, Col: 4}, {Line: 2, Col: 2}},
		},
		{
			name: "CR before CRLF and CRLF between lines",
			text: "a\r\r\nb\r\n",
			want: []Pos{{Line: 1, Col: 2}},
		},
		{
			name: "BOM and astral rune count as one column",
			text: "\ufeff😀\r\n\u0085\u2028\u2029\r",
			want: []Pos{{Line: 2, Col: 4}},
		},
		{
			name: "non-LF separators do not advance line",
			text: "x\u0085\u2028\u2029\r",
			want: []Pos{{Line: 1, Col: 5}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := BareCRs(test.text); !reflect.DeepEqual(got, test.want) {
				t.Errorf("BareCRs(%q) = %#v, want %#v", test.text, got, test.want)
			}
		})
	}
}
