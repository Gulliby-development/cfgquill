package config

import (
	"reflect"
	"testing"

	"github.com/Gulliby-development/cfgquill/internal/base"
)

func TestFind(t *testing.T) {
	var find func(File, string) []Entry = File.Find
	file := File{
		Entries: []Entry{
			{Line: 1, Name: "", Var: "b", Key: "b", Value: "", HasValue: false},
			{Line: 2, Name: "", Var: "b", Key: "b", Value: "", HasValue: true},
			{Line: 3, Name: "a.S", Var: "b", Key: "a.S.b", Value: "one", HasValue: true},
			{Line: 4, Name: "a.s", Var: "b", Key: "a.s.b", Value: "two", HasValue: true},
		},
		Headers: []Header{
			{Line: 3, Name: "a.S", Count: 1},
			{Line: 4, Name: "a.s", Count: 1},
		},
	}
	before := File{
		Entries: []Entry{
			{Line: 1, Name: "", Var: "b", Key: "b", Value: "", HasValue: false},
			{Line: 2, Name: "", Var: "b", Key: "b", Value: "", HasValue: true},
			{Line: 3, Name: "a.S", Var: "b", Key: "a.S.b", Value: "one", HasValue: true},
			{Line: 4, Name: "a.s", Var: "b", Key: "a.s.b", Value: "two", HasValue: true},
		},
		Headers: []Header{
			{Line: 3, Name: "a.S", Count: 1},
			{Line: 4, Name: "a.s", Count: 1},
		},
	}
	tests := []struct {
		canon string
		want  []Entry
	}{
		{canon: "missing", want: []Entry{}},
		{
			canon: "b",
			want: []Entry{
				{Line: 1, Name: "", Var: "b", Key: "b", Value: "", HasValue: false},
				{Line: 2, Name: "", Var: "b", Key: "b", Value: "", HasValue: true},
			},
		},
		{
			canon: "a.S.b",
			want:  []Entry{{Line: 3, Name: "a.S", Var: "b", Key: "a.S.b", Value: "one", HasValue: true}},
		},
		{
			canon: "a.s.b",
			want:  []Entry{{Line: 4, Name: "a.s", Var: "b", Key: "a.s.b", Value: "two", HasValue: true}},
		},
		{canon: "A.S.b", want: []Entry{}},
	}
	for _, test := range tests {
		got := find(file, test.canon)
		if len(got) != len(test.want) {
			t.Errorf("Find(%q) returned %d entries, want %d", test.canon, len(got), len(test.want))
			continue
		}
		for i := range test.want {
			if got[i] != test.want[i] {
				t.Errorf("Find(%q)[%d] = %#v, want %#v", test.canon, i, got[i], test.want[i])
			}
		}
		if !reflect.DeepEqual(file, before) {
			t.Fatalf("Find(%q) changed the File: %#v", test.canon, file)
		}
	}
}

func TestLoadSuccess(t *testing.T) {
	var load func(*base.Inputs, string) (Loaded, error) = Load
	tests := []struct {
		name  string
		arg   string
		files map[string][]byte
		stdin []byte
		want  Loaded
	}{
		{
			name:  "named file preserves BOM and CRLF",
			arg:   "t.cfg",
			files: map[string][]byte{"t.cfg": []byte("\ufeff[a]\r\nB=1")},
			want: Loaded{
				Display: "t.cfg",
				Text:    "\ufeff[a]\r\nB=1",
				File: File{
					Entries: []Entry{{Line: 2, Name: "a", Var: "b", Key: "a.b", Value: "1", HasValue: true}},
					Headers: []Header{{Line: 1, Name: "a", Count: 1}},
				},
			},
		},
		{
			name:  "stdin preserves BOM and CRLF",
			arg:   "-",
			stdin: []byte("\ufeff[a]\r\nB=1"),
			want: Loaded{
				Display: "<stdin>",
				Text:    "\ufeff[a]\r\nB=1",
				File: File{
					Entries: []Entry{{Line: 2, Name: "a", Var: "b", Key: "a.b", Value: "1", HasValue: true}},
					Headers: []Header{{Line: 1, Name: "a", Count: 1}},
				},
			},
		},
		{
			name:  "empty file",
			arg:   "empty.cfg",
			files: map[string][]byte{"empty.cfg": []byte{}},
			want: Loaded{
				Display: "empty.cfg",
				Text:    "",
				File:    File{Entries: []Entry{}, Headers: []Header{}},
			},
		},
		{
			name:  "Unicode display name keeps original read path",
			arg:   "café😀.cfg",
			files: map[string][]byte{"café😀.cfg": []byte("[Core]\nb=2")},
			want: Loaded{
				Display: "caf??.cfg",
				Text:    "[Core]\nb=2",
				File: File{
					Entries: []Entry{{Line: 2, Name: "core", Var: "b", Key: "core.b", Value: "2", HasValue: true}},
					Headers: []Header{{Line: 1, Name: "core", Count: 1}},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			memory := &base.MemIO{Files: test.files, Stdin: test.stdin}
			got, err := load(base.NewInputs(memory), test.arg)
			if err != nil {
				t.Fatalf("Load(%q) error = %v, want nil", test.arg, err)
			}
			if got.Display != test.want.Display || got.Text != test.want.Text {
				t.Errorf("Load(%q) source = (%q, %q), want (%q, %q)", test.arg, got.Display, got.Text, test.want.Display, test.want.Text)
			}
			if len(got.File.Entries) != len(test.want.File.Entries) || len(got.File.Headers) != len(test.want.File.Headers) {
				t.Errorf("Load(%q) model lengths = (%d entries, %d headers), want (%d, %d)", test.arg, len(got.File.Entries), len(got.File.Headers), len(test.want.File.Entries), len(test.want.File.Headers))
			} else {
				for i := range test.want.File.Entries {
					if got.File.Entries[i] != test.want.File.Entries[i] {
						t.Errorf("Load(%q) entry %d = %#v, want %#v", test.arg, i, got.File.Entries[i], test.want.File.Entries[i])
					}
				}
				for i := range test.want.File.Headers {
					if got.File.Headers[i] != test.want.File.Headers[i] {
						t.Errorf("Load(%q) header %d = %#v, want %#v", test.arg, i, got.File.Headers[i], test.want.File.Headers[i])
					}
				}
			}
			if memory.Stdout != "" || memory.Stderr != "" {
				t.Errorf("Load(%q) wrote output: stdout %q, stderr %q", test.arg, memory.Stdout, memory.Stderr)
			}
		})
	}
}

func TestLoadErrors(t *testing.T) {
	var load func(*base.Inputs, string) (Loaded, error) = Load
	tests := []struct {
		name  string
		arg   string
		files map[string][]byte
		stdin []byte
		want  string
	}{
		{name: "missing file", arg: "missing.cfg", want: "missing.cfg: cannot read the file"},
		{name: "empty path", arg: "", want: ": cannot read the file"},
		{
			name:  "syntax fault",
			arg:   "t.cfg",
			files: map[string][]byte{"t.cfg": []byte("[a]\nb c\n")},
			want:  "t.cfg:2: CQ003 expected \"=\" after the variable name",
		},
		{
			name:  "unknown escape",
			arg:   "t.cfg",
			files: map[string][]byte{"t.cfg": []byte("[a]\nb=x\\q")},
			want:  `t.cfg:2: CQ005 unknown escape "\q"`,
		},
		{
			name:  "invalid UTF-8 file bytes",
			arg:   "t.cfg",
			files: map[string][]byte{"t.cfg": []byte{91, 97, 93, 10, 98, 255, 10}},
			want:  "t.cfg: input is not valid UTF-8",
		},
		{
			name:  "invalid UTF-8 stdin bytes",
			arg:   "-",
			stdin: []byte{91, 97, 93, 10, 98, 255, 10},
			want:  "<stdin>: input is not valid UTF-8",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			memory := &base.MemIO{Files: test.files, Stdin: test.stdin}
			got, err := load(base.NewInputs(memory), test.arg)
			if err == nil || err.Error() != test.want {
				t.Errorf("Load(%q) error = %v, want %q", test.arg, err, test.want)
			}
			if !reflect.DeepEqual(got, Loaded{}) {
				t.Errorf("Load(%q) returned partial result %#v on error", test.arg, got)
			}
			if memory.Stdout != "" || memory.Stderr != "" {
				t.Errorf("Load(%q) wrote output: stdout %q, stderr %q", test.arg, memory.Stdout, memory.Stderr)
			}
		})
	}
}

func TestLoadConsumesStdin(t *testing.T) {
	var load func(*base.Inputs, string) (Loaded, error) = Load

	successIO := &base.MemIO{Stdin: []byte("[a]\nb=1")}
	successInputs := base.NewInputs(successIO)
	first, err := load(successInputs, "-")
	if err != nil || first.Display != "<stdin>" || first.Text != "[a]\nb=1" {
		t.Fatalf("first stdin Load = %#v, %v, want successful input", first, err)
	}
	second, err := load(successInputs, "-")
	if err == nil || err.Error() != "standard input is named twice" || !reflect.DeepEqual(second, Loaded{}) {
		t.Errorf("second stdin Load = %#v, %v, want duplicate-input error and zero result", second, err)
	}
	if successIO.Stdout != "" || successIO.Stderr != "" {
		t.Errorf("stdin Load wrote output: stdout %q, stderr %q", successIO.Stdout, successIO.Stderr)
	}

	syntaxIO := &base.MemIO{Stdin: []byte("[a]\nb c")}
	syntaxInputs := base.NewInputs(syntaxIO)
	first, err = load(syntaxInputs, "-")
	if err == nil || err.Error() != "<stdin>:2: CQ003 expected \"=\" after the variable name" || !reflect.DeepEqual(first, Loaded{}) {
		t.Errorf("syntax-fault stdin Load = %#v, %v, want syntax error and zero result", first, err)
	}
	second, err = load(syntaxInputs, "-")
	if err == nil || err.Error() != "standard input is named twice" || !reflect.DeepEqual(second, Loaded{}) {
		t.Errorf("stdin Load after syntax fault = %#v, %v, want duplicate-input error and zero result", second, err)
	}
	if syntaxIO.Stdout != "" || syntaxIO.Stderr != "" {
		t.Errorf("stdin Load wrote output: stdout %q, stderr %q", syntaxIO.Stdout, syntaxIO.Stderr)
	}
}
