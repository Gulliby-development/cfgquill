package config

import (
	"testing"

	"github.com/Gulliby-development/cfgquill/internal/base"
)

func TestParseModels(t *testing.T) {
	tests := []struct {
		name string
		text string
		want File
	}{
		{
			name: "empty text",
			text: "",
			want: File{},
		},
		{
			name: "initial BOM only",
			text: "\ufeff",
			want: File{},
		},
		{
			name: "whitespace and line comments",
			text: " \t\r\n# x\\\n; y\n",
			want: File{},
		},
		{
			name: "headerless variable without value",
			text: "B",
			want: File{
				Entries: []Entry{{Line: 1, Name: "", Var: "b", Key: "b", Value: "", HasValue: false}},
			},
		},
		{
			name: "headerless empty value",
			text: "B = ",
			want: File{
				Entries: []Entry{{Line: 1, Name: "", Var: "b", Key: "b", Value: "", HasValue: true}},
			},
		},
		{
			name: "header without entries",
			text: "[a]",
			want: File{
				Headers: []Header{{Line: 1, Name: "a", Count: 0}},
			},
		},
		{
			name: "CRLF, bare CR and terminal entries",
			text: "\ufeff[a]\r\nB=1\r\nc=2",
			want: File{
				Entries: []Entry{
					{Line: 2, Name: "a", Var: "b", Key: "a.b", Value: "1", HasValue: true},
					{Line: 3, Name: "a", Var: "c", Key: "a.c", Value: "2", HasValue: true},
				},
				Headers: []Header{{Line: 1, Name: "a", Count: 2}},
			},
		},
		{
			name: "headerless and same-line entries with header counts",
			text: "k=1\n[Core] A-2\n[remote \"Origin\"]b=\"\"\n[a]\n",
			want: File{
				Entries: []Entry{
					{Line: 1, Name: "", Var: "k", Key: "k", Value: "1", HasValue: true},
					{Line: 2, Name: "core", Var: "a-2", Key: "core.a-2", Value: "", HasValue: false},
					{Line: 3, Name: "remote.Origin", Var: "b", Key: "remote.Origin.b", Value: "", HasValue: true},
				},
				Headers: []Header{
					{Line: 2, Name: "core", Count: 1},
					{Line: 3, Name: "remote.Origin", Count: 1},
					{Line: 4, Name: "a", Count: 0},
				},
			},
		},
		{
			name: "old, extended and empty subsections",
			text: "[a.B]c=1\n[a \"B\"]c=2\n[a \"\"]c=3",
			want: File{
				Entries: []Entry{
					{Line: 1, Name: "a.b", Var: "c", Key: "a.b.c", Value: "1", HasValue: true},
					{Line: 2, Name: "a.B", Var: "c", Key: "a.B.c", Value: "2", HasValue: true},
					{Line: 3, Name: "a.", Var: "c", Key: "a..c", Value: "3", HasValue: true},
				},
				Headers: []Header{
					{Line: 1, Name: "a.b", Count: 1},
					{Line: 2, Name: "a.B", Count: 1},
					{Line: 3, Name: "a.", Count: 1},
				},
			},
		},
		{
			name: "adjacent headers and entry after the last header",
			text: "[a][b] x\n",
			want: File{
				Entries: []Entry{{Line: 1, Name: "b", Var: "x", Key: "b.x", Value: "", HasValue: false}},
				Headers: []Header{
					{Line: 1, Name: "a", Count: 0},
					{Line: 1, Name: "b", Count: 1},
				},
			},
		},
		{
			name: "quoted empty text flushes leading pending blanks",
			text: "b=\"\" a",
			want: File{
				Entries: []Entry{{Line: 1, Name: "", Var: "b", Key: "b", Value: "a", HasValue: true}},
			},
		},
		{
			name: "quoted empty text preserves prior pending blank",
			text: "b=a \"\"",
			want: File{
				Entries: []Entry{{Line: 1, Name: "", Var: "b", Key: "b", Value: "a ", HasValue: true}},
			},
		},
		{
			name: "quoted empty text separates pending blanks",
			text: "b=a \"\" b",
			want: File{
				Entries: []Entry{{Line: 1, Name: "", Var: "b", Key: "b", Value: "a  b", HasValue: true}},
			},
		},
		{
			name: "continued lines and preserved interior blanks",
			text: "b=one \\\n  two\\\nthree\n",
			want: File{
				Entries: []Entry{{Line: 1, Name: "", Var: "b", Key: "b", Value: "one   twothree", HasValue: true}},
			},
		},
		{
			name: "bare CR remains value text and line",
			text: "b=x\ry = 2\n",
			want: File{
				Entries: []Entry{{Line: 1, Name: "", Var: "b", Key: "b", Value: "x\ry = 2", HasValue: true}},
			},
		},
		{
			name: "CRLF ends the value",
			text: "b=x\r\n",
			want: File{
				Entries: []Entry{{Line: 1, Name: "", Var: "b", Key: "b", Value: "x", HasValue: true}},
			},
		},
		{
			name: "comment markers end an unquoted value",
			text: "b=a#b;c\n",
			want: File{
				Entries: []Entry{{Line: 1, Name: "", Var: "b", Key: "b", Value: "a", HasValue: true}},
			},
		},
		{
			name: "terminal backslash continues to synthetic newline",
			text: "b=a\\",
			want: File{
				Entries: []Entry{{Line: 1, Name: "", Var: "b", Key: "b", Value: "a", HasValue: true}},
			},
		},
		{
			name: "quoted continuation",
			text: "b=\"a\\\n  b\"",
			want: File{
				Entries: []Entry{{Line: 1, Name: "", Var: "b", Key: "b", Value: "a  b", HasValue: true}},
			},
		},
		{
			name: "supported escapes",
			text: "b=\\t\\b\\n\\\\\\\"",
			want: File{
				Entries: []Entry{{Line: 1, Name: "", Var: "b", Key: "b", Value: "\t\b\n\\\"", HasValue: true}},
			},
		},
		{
			name: "Unicode non-blanks and astral rune remain value text",
			text: "b=\u3000x\u00a0\u0085\u2028\u2029\u200b\ufeff😀",
			want: File{
				Entries: []Entry{{Line: 1, Name: "", Var: "b", Key: "b", Value: "\u3000x\u00a0\u0085\u2028\u2029\u200b\ufeff😀", HasValue: true}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, fault := Parse(test.text)
			if fault != nil {
				t.Fatalf("Parse(%q) fault = %#v, want nil", test.text, fault)
			}
			if len(got.Entries) != len(test.want.Entries) {
				t.Errorf("Parse(%q) has %d entries, want %d: %#v", test.text, len(got.Entries), len(test.want.Entries), got.Entries)
			} else {
				for i, want := range test.want.Entries {
					if got.Entries[i] != want {
						t.Errorf("Parse(%q) entry %d = %#v, want %#v", test.text, i, got.Entries[i], want)
					}
				}
			}
			if len(got.Headers) != len(test.want.Headers) {
				t.Errorf("Parse(%q) has %d headers, want %d: %#v", test.text, len(got.Headers), len(test.want.Headers), got.Headers)
			} else {
				for i, want := range test.want.Headers {
					if got.Headers[i] != want {
						t.Errorf("Parse(%q) header %d = %#v, want %#v", test.text, i, got.Headers[i], want)
					}
				}
			}
		})
	}
}

func TestParseFaults(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		wantFile  File
		wantFault base.Fault
	}{
		{name: "empty section name", text: "[]", wantFile: File{}, wantFault: base.Fault{Code: "CQ001", Line: 1}},
		{name: "line feed inside old header", text: "[a\n", wantFile: File{}, wantFault: base.Fault{Code: "CQ001", Line: 1}},
		{name: "second BOM is not skipped", text: "\ufeff\ufeff", wantFile: File{}, wantFault: base.Fault{Code: "CQ002", Line: 1}},
		{name: "BOM after a space is not skipped", text: " \ufeff", wantFile: File{}, wantFault: base.Fault{Code: "CQ002", Line: 1}},
		{name: "no-break space is not a blank", text: "\u00a0", wantFile: File{}, wantFault: base.Fault{Code: "CQ002", Line: 1}},
		{name: "ideographic space is not a blank", text: "\u3000", wantFile: File{}, wantFault: base.Fault{Code: "CQ002", Line: 1}},
		{name: "NEL is not a line end", text: "\u0085", wantFile: File{}, wantFault: base.Fault{Code: "CQ002", Line: 1}},
		{name: "line separator is not a line end", text: "\u2028", wantFile: File{}, wantFault: base.Fault{Code: "CQ002", Line: 1}},
		{name: "paragraph separator is not a line end", text: "\u2029", wantFile: File{}, wantFault: base.Fault{Code: "CQ002", Line: 1}},
		{name: "zero-width space is not a blank", text: "\u200b", wantFile: File{}, wantFault: base.Fault{Code: "CQ002", Line: 1}},
		{name: "vertical tab is not a blank", text: "\v", wantFile: File{}, wantFault: base.Fault{Code: "CQ002", Line: 1}},
		{name: "form feed is not a blank", text: "\f", wantFile: File{}, wantFault: base.Fault{Code: "CQ002", Line: 1}},
		{name: "digit cannot start a variable", text: "2c=1", wantFile: File{}, wantFault: base.Fault{Code: "CQ002", Line: 1}},
		{name: "space before delimiter is not allowed", text: "b c", wantFile: File{}, wantFault: base.Fault{Code: "CQ003", Line: 1}},
		{name: "hash before delimiter is not allowed", text: "b#c", wantFile: File{}, wantFault: base.Fault{Code: "CQ003", Line: 1}},
		{name: "semicolon before delimiter is not allowed", text: "b;c", wantFile: File{}, wantFault: base.Fault{Code: "CQ003", Line: 1}},
		{name: "bare CR after variable is not skipped", text: "b\r", wantFile: File{}, wantFault: base.Fault{Code: "CQ003", Line: 1}},
		{name: "underscore is not part of variable", text: "b_", wantFile: File{}, wantFault: base.Fault{Code: "CQ003", Line: 1}},
		{name: "unterminated quoted value", text: "b=\"x\n", wantFile: File{}, wantFault: base.Fault{Code: "CQ004", Line: 1}},
		{name: "unknown escape", text: "b=x\\q", wantFile: File{}, wantFault: base.Fault{Code: "CQ005", Line: 1, Arg: `\q`}},
		{name: "CR after backslash is an unknown escape", text: "b=x\\\ry", wantFile: File{}, wantFault: base.Fault{Code: "CQ005", Line: 1, Arg: "\\\r"}},
		{name: "first fault after a continuation is on line two", text: "b=x\\\n\\q", wantFile: File{}, wantFault: base.Fault{Code: "CQ005", Line: 2, Arg: `\q`}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, fault := Parse(test.text)
			if fault == nil {
				t.Fatalf("Parse(%q) fault = nil, want %#v", test.text, test.wantFault)
			}
			if *fault != test.wantFault {
				t.Errorf("Parse(%q) fault = %#v, want %#v", test.text, fault, test.wantFault)
			}
			if len(got.Entries) != len(test.wantFile.Entries) {
				t.Errorf("Parse(%q) retained %d entries, want %d: %#v", test.text, len(got.Entries), len(test.wantFile.Entries), got.Entries)
			} else {
				for i, want := range test.wantFile.Entries {
					if got.Entries[i] != want {
						t.Errorf("Parse(%q) retained entry %d = %#v, want %#v", test.text, i, got.Entries[i], want)
					}
				}
			}
			if len(got.Headers) != len(test.wantFile.Headers) {
				t.Errorf("Parse(%q) retained %d headers, want %d: %#v", test.text, len(got.Headers), len(test.wantFile.Headers), got.Headers)
			} else {
				for i, want := range test.wantFile.Headers {
					if got.Headers[i] != want {
						t.Errorf("Parse(%q) retained header %d = %#v, want %#v", test.text, i, got.Headers[i], want)
					}
				}
			}
		})
	}
}

func TestParseFaultRetainsCompletedPrefix(t *testing.T) {
	text := "[a]\nb=1\n[c\nd=2\n"
	want := File{
		Entries: []Entry{{Line: 2, Name: "a", Var: "b", Key: "a.b", Value: "1", HasValue: true}},
		Headers: []Header{{Line: 1, Name: "a", Count: 1}},
	}
	got, fault := Parse(text)
	if fault == nil || *fault != (base.Fault{Code: "CQ001", Line: 3}) {
		t.Fatalf("Parse(%q) fault = %#v, want CQ001 at line 3", text, fault)
	}
	if len(got.Entries) != len(want.Entries) || len(got.Headers) != len(want.Headers) {
		t.Fatalf("Parse(%q) prefix = %#v, want %#v", text, got, want)
	}
	for i, entry := range want.Entries {
		if got.Entries[i] != entry {
			t.Errorf("Parse(%q) prefix entry %d = %#v, want %#v", text, i, got.Entries[i], entry)
		}
	}
	for i, header := range want.Headers {
		if got.Headers[i] != header {
			t.Errorf("Parse(%q) prefix header %d = %#v, want %#v", text, i, got.Headers[i], header)
		}
	}
}
