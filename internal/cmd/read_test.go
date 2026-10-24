package cmd

import (
	"testing"

	"github.com/Gulliby-development/cfgquill/internal/base"
)

var _ func([]string, base.IO) (int, error) = RunList
var _ func([]string, base.IO) (int, error) = RunSections
var _ func([]string, base.IO) (int, error) = RunGet
var _ func([]string, base.IO) (int, error) = RunBool

func TestReadCommandAPIs(t *testing.T) {
	if ListUsage != "list FILE" {
		t.Errorf("ListUsage = %q, want %q", ListUsage, "list FILE")
	}
	if SectionsUsage != "sections FILE" {
		t.Errorf("SectionsUsage = %q, want %q", SectionsUsage, "sections FILE")
	}
	if GetUsage != "get [--all] KEY FILE" {
		t.Errorf("GetUsage = %q, want %q", GetUsage, "get [--all] KEY FILE")
	}
	if BoolUsage != "bool KEY FILE" {
		t.Errorf("BoolUsage = %q, want %q", BoolUsage, "bool KEY FILE")
	}
}

func TestReadValidationPrecedence(t *testing.T) {
	tests := []struct {
		name string
		run  func([]string, base.IO) (int, error)
		args []string
		want string
	}{
		{"first unknown list option", RunList, []string{"--oops", "--bad"}, "unknown option \"--oops\""},
		{"list arity", RunList, nil, "usage: list FILE"},
		{"sections excess files", RunSections, []string{"t.cfg", "missing.cfg"}, "usage: sections FILE"},
		{"get excess arguments", RunGet, []string{"--all", "a.b", "t.cfg", "missing.cfg"}, "usage: get [--all] KEY FILE"},
		{"bool arity before key", RunBool, []string{".b"}, "usage: bool KEY FILE"},
		{"get key before read", RunGet, []string{".b", "missing.cfg"}, "invalid key \".b\""},
		{"bool rejects LF in subsection", RunBool, []string{"a.x\ny.b", "missing.cfg"}, "invalid key \"a.x?y.b\""},
		{"bool rejects get option", RunBool, []string{"--all", "a.b", "t.cfg"}, "unknown option \"--all\""},
		{"list rejects option marker", RunList, []string{"--", "t.cfg"}, "unknown option \"--\""},
		{"sections rejects help option", RunSections, []string{"--help"}, "unknown option \"--help\""},
		{"get requires exact option", RunGet, []string{"--all=1", "a.b", "t.cfg"}, "unknown option \"--all=1\""},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			memory := &base.MemIO{Files: map[string][]byte{
				"t.cfg": []byte("[a]\nb = 1\n"),
			}}
			_, err := test.run(test.args, memory)
			if err == nil || err.Error() != test.want {
				t.Errorf("error = %v, want %q", err, test.want)
			}
			if memory.Stdout != "" || memory.Stderr != "" {
				t.Errorf("output = %q, stderr = %q, want both empty", memory.Stdout, memory.Stderr)
			}
		})
	}
}

func TestReadListInputLoading(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		hasFile    bool
		file       []byte
		stdin      []byte
		wantStatus int
		wantOut    string
		wantErr    string
	}{
		{"missing file", "missing.cfg", false, nil, nil, 0, "", "missing.cfg: cannot read the file"},
		{"sanitized display", "café😀.cfg", false, nil, nil, 0, "", "caf??.cfg: cannot read the file"},
		{"empty path", "", false, nil, nil, 0, "", ": cannot read the file"},
		{"invalid UTF-8 file", "t.cfg", true, []byte{0x5b, 0x61, 0x5d, 0x0a, 0x62, 0xff, 0x0a}, nil, 0, "", "t.cfg: input is not valid UTF-8"},
		{"valid stdin", "-", false, nil, []byte("[a]\nb = 1\n"), 0, "LINE  KEY  VALUE\n2     a.b  \"1\"\n", ""},
		{"invalid UTF-8 stdin", "-", false, nil, []byte{0x5b, 0x61, 0x5d, 0x0a, 0x62, 0xff, 0x0a}, 0, "", "<stdin>: input is not valid UTF-8"},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			memory := &base.MemIO{Stdin: test.stdin, Files: map[string][]byte{}}
			if test.hasFile {
				memory.Files[test.path] = test.file
			}
			status, err := RunList([]string{test.path}, memory)
			if test.wantErr != "" {
				if err == nil || err.Error() != test.wantErr {
					t.Errorf("error = %v, want %q", err, test.wantErr)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if status != test.wantStatus {
				t.Errorf("status = %d, want %d", status, test.wantStatus)
			}
			if memory.Stdout != test.wantOut || memory.Stderr != "" {
				t.Errorf("output = %q, stderr = %q, want output %q and empty stderr", memory.Stdout, memory.Stderr, test.wantOut)
			}
		})
	}
}

func TestReadListOutput(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{"empty file suppresses header", "", ""},
		{"comments and blanks suppress header", "# nothing here\n; nor here\n\n   \n", ""},
		{"empty and absent values", "[a]\nb =\nc = \"\"\nd\n", "LINE  KEY  VALUE\n2     a.b  \"\"\n3     a.c  \"\"\n4     a.d  -\n"},
		{"CRLF and missing final LF", "[a]\r\nb = 1\r\nc = 2", "LINE  KEY  VALUE\n2     a.b  \"1\"\n3     a.c  \"2\"\n"},
		{"initial BOM", "\ufeff[a]\nb = 1\n", "LINE  KEY  VALUE\n2     a.b  \"1\"\n"},
		{"stored subsection case and order", "[a.B]\nc = 1\n[a \"B\"]\nc = 2\n[A]\nC = 3\n", "LINE  KEY    VALUE\n2     a.b.c  \"1\"\n4     a.B.c  \"2\"\n6     a.c    \"3\"\n"},
		{"rune-width key and sanitized value", "[a \"😀\"]\nb = x \n[a]\nb\n", "LINE  KEY    VALUE\n2     a.😀.b  \"x?\"\n4     a.b    -\n"},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			memory := &base.MemIO{Files: map[string][]byte{"t.cfg": []byte(test.text)}}
			status, err := RunList([]string{"t.cfg"}, memory)
			if err != nil || status != 0 {
				t.Errorf("RunList status/error = %d, %v, want 0, nil", status, err)
			}
			if memory.Stdout != test.want || memory.Stderr != "" {
				t.Errorf("output = %q, stderr = %q, want %q and empty stderr", memory.Stdout, memory.Stderr, test.want)
			}
		})
	}
}

func TestReadSectionsOutput(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{"headerless entries suppress header", "k = 1\n", ""},
		{"empty file suppresses header", "", ""},
		{"ordered section counts and first-dot split", "k = 1\n[core]\n a = 1\n b = 2\n[remote \"origin\"]\n url = u\n[remote \"up stream\"]\n[Branch.Main]\n x\n y\n[a.b \"C\"] z = 1\n[a]\n", "LINE  SECTION  SUBSECTION   ENTRIES\n2     core     -            2\n5     remote   \"origin\"     1\n7     remote   \"up stream\"  0\n8     branch   \"main\"       2\n11    a        \"b.C\"        1\n12    a        -            0\n"},
		{"empty subsection", "[a \"\"]\n", "LINE  SECTION  SUBSECTION  ENTRIES\n1     a        \"\"          0\n"},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			memory := &base.MemIO{Files: map[string][]byte{"t.cfg": []byte(test.text)}}
			status, err := RunSections([]string{"t.cfg"}, memory)
			if err != nil || status != 0 {
				t.Errorf("RunSections status/error = %d, %v, want 0, nil", status, err)
			}
			if memory.Stdout != test.want || memory.Stderr != "" {
				t.Errorf("output = %q, stderr = %q, want %q and empty stderr", memory.Stdout, memory.Stderr, test.want)
			}
		})
	}
}

func TestReadGetValues(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		text       string
		wantStatus int
		wantOut    string
		wantErr    string
	}{
		{"last value", []string{"a.b", "t.cfg"}, "[a]\nb = 1\nb = 2\n", 0, "2\n", ""},
		{"repeated all option preserves order", []string{"--all", "a.b", "--all", "t.cfg"}, "[a]\nb = 1\nb = 2\n", 0, "1\n2\n", ""},
		{"absent and empty values with all", []string{"--all", "a.b", "t.cfg"}, "[a]\nb\nb = \"\"\n", 0, "\n\n", ""},
		{"canonical section and variable", []string{"A.S.B", "t.cfg"}, "[a \"S\"]\nb = 1\n", 0, "1\n", ""},
		{"subsection case mismatch", []string{"a.s.b", "t.cfg"}, "[a \"S\"]\nb = 1\n", 1, "", ""},
		{"old form subsection mismatch", []string{"a.S.b", "t.cfg"}, "[a.S]\nb = 1\n", 1, "", ""},
		{"old form lowercased match", []string{"a.s.b", "t.cfg"}, "[a.S]\nb = 1\n", 0, "1\n", ""},
		{"empty subsection", []string{"a..c", "t.cfg"}, "[a \"\"]\nc = 1\n", 0, "1\n", ""},
		{"quoted blanks preserved", []string{"a.b", "t.cfg"}, "[a]\nb = \"   \"\n", 0, "   \n", ""},
		{"bare CR remains in value", []string{"a.b", "t.cfg"}, "[a]\nb = x\ry = 2\n", 0, "x\ry = 2\n", ""},
		{"bare CR does not create a key", []string{"a.y", "t.cfg"}, "[a]\nb = x\ry = 2\n", 1, "", ""},
		{"raw Unicode remains in value", []string{"a.b", "t.cfg"}, "[a]\nb = \"x\u2028y\u2029z\ufeffw\u0085v\"\n", 0, "x\u2028y\u2029z\ufeffw\u0085v\n", ""},
		{"continued value", []string{"a.b", "t.cfg"}, "[a]\nb = one \\\n  two\\\nthree\n", 0, "one   twothree\n", ""},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			memory := &base.MemIO{Files: map[string][]byte{"t.cfg": []byte(test.text)}}
			status, err := RunGet(test.args, memory)
			if status != test.wantStatus {
				t.Errorf("RunGet status = %d, want %d", status, test.wantStatus)
			}
			if test.wantErr == "" {
				if err != nil {
					t.Errorf("RunGet error = %v, want nil", err)
				}
			} else if err == nil || err.Error() != test.wantErr {
				t.Errorf("RunGet error = %v, want %q", err, test.wantErr)
			}
			if memory.Stdout != test.wantOut || memory.Stderr != "" {
				t.Errorf("RunGet output = %q, stderr = %q, want %q and empty stderr", memory.Stdout, memory.Stderr, test.wantOut)
			}
		})
	}
}

func TestReadBoolValues(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		text       string
		wantStatus int
		wantOut    string
		wantErr    string
	}{
		{"absent value is true", []string{"a.b", "t.cfg"}, "[a]\nb\n", 0, "true\n", ""},
		{"empty value is false", []string{"a.b", "t.cfg"}, "[a]\nb = \"\"\n", 0, "false\n", ""},
		{"case-insensitive true word", []string{"a.b", "t.cfg"}, "[a]\nb = TrUe\n", 0, "true\n", ""},
		{"false word", []string{"a.b", "t.cfg"}, "[a]\nb = off\n", 0, "false\n", ""},
		{"zero digits", []string{"a.b", "t.cfg"}, "[a]\nb = 000\n", 0, "false\n", ""},
		{"nine nonzero digits", []string{"a.b", "t.cfg"}, "[a]\nb = 123456789\n", 0, "true\n", ""},
		{"ten digits rejected", []string{"a.b", "t.cfg"}, "[a]\nb = 0123456789\n", 0, "", "t.cfg:2: bad boolean \"0123456789\" for a.b"},
		{"leading blank rejected", []string{"a.b", "t.cfg"}, "[a]\nb = \" true\"\n", 0, "", "t.cfg:2: bad boolean \" true\" for a.b"},
		{"NBSP rejected", []string{"a.b", "t.cfg"}, "[a]\nb = true \n", 0, "", "t.cfg:2: bad boolean \"true?\" for a.b"},
		{"earlier invalid value ignored", []string{"a.b", "t.cfg"}, "[a]\nb = maybe\nb = yes\n", 0, "true\n", ""},
		{"last invalid value reported", []string{"a.b", "t.cfg"}, "[a]\nb = yes\nb = maybe\n", 0, "", "t.cfg:3: bad boolean \"maybe\" for a.b"},
		{"missing key is negative answer", []string{"a.b", "t.cfg"}, "[a]\nc = 1\n", 1, "", ""},
		{"canonical key in error", []string{"A.S.B", "t.cfg"}, "[a \"S\"]\nb = 2x\n", 0, "", "t.cfg:2: bad boolean \"2x\" for a.S.b"},
		{"signed number rejected", []string{"a.b", "t.cfg"}, "[a]\nb = +1\n", 0, "", "t.cfg:2: bad boolean \"+1\" for a.b"},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			memory := &base.MemIO{Files: map[string][]byte{"t.cfg": []byte(test.text)}}
			status, err := RunBool(test.args, memory)
			if status != test.wantStatus {
				t.Errorf("RunBool status = %d, want %d", status, test.wantStatus)
			}
			if test.wantErr == "" {
				if err != nil {
					t.Errorf("RunBool error = %v, want nil", err)
				}
			} else if err == nil || err.Error() != test.wantErr {
				t.Errorf("RunBool error = %v, want %q", err, test.wantErr)
			}
			if memory.Stdout != test.wantOut || memory.Stderr != "" {
				t.Errorf("RunBool output = %q, stderr = %q, want %q and empty stderr", memory.Stdout, memory.Stderr, test.wantOut)
			}
		})
	}
}

func TestReadCommandsParseCompleteFile(t *testing.T) {
	tests := []struct {
		name string
		run  func([]string, base.IO) (int, error)
		args []string
	}{
		{"list late syntax error", RunList, []string{"t.cfg"}},
		{"sections late syntax error", RunSections, []string{"t.cfg"}},
		{"get matching key before late syntax error", RunGet, []string{"a.b", "t.cfg"}},
		{"get absent key before late syntax error", RunGet, []string{"a.z", "t.cfg"}},
		{"bool invalid value before late syntax error", RunBool, []string{"a.b", "t.cfg"}},
		{"bool absent key before late syntax error", RunBool, []string{"a.z", "t.cfg"}},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			memory := &base.MemIO{Files: map[string][]byte{"t.cfg": []byte("[a]\nb = maybe\n[c\n")}}
			status, err := test.run(test.args, memory)
			if status != 0 || err == nil || err.Error() != "t.cfg:3: CQ001 bad section header" {
				t.Errorf("command result = %d, %v, want 0 and late syntax error", status, err)
			}
			if memory.Stdout != "" || memory.Stderr != "" {
				t.Errorf("output = %q, stderr = %q, want both empty", memory.Stdout, memory.Stderr)
			}
		})
	}

	syntaxCases := []struct {
		name string
		text string
		want string
	}{
		{"unterminated quote", "[a]\nb = \"x\n", "t.cfg:2: CQ004 unterminated quote"},
		{"escape before raw CR", "[a]\nb = x\\\ry\n", "t.cfg:2: CQ005 unknown escape \"\\?\""},
		{"BOM inside file", "[a]\n\ufeffb = 1\n", "t.cfg:2: CQ002 bad variable name"},
	}
	for _, test := range syntaxCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			memory := &base.MemIO{Files: map[string][]byte{"t.cfg": []byte(test.text)}}
			status, err := RunList([]string{"t.cfg"}, memory)
			if status != 0 || err == nil || err.Error() != test.want {
				t.Errorf("RunList result = %d, %v, want error %q", status, err, test.want)
			}
			if memory.Stdout != "" || memory.Stderr != "" {
				t.Errorf("output = %q, stderr = %q, want both empty", memory.Stdout, memory.Stderr)
			}
		})
	}
}
