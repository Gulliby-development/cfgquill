package lint

import (
	"testing"

	"github.com/Gulliby-development/cfgquill/internal/base"
)

func TestRunContracts(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		files  map[string][]byte
		stdin  []byte
		status int
		err    string
		stdout string
	}{
		{
			name: "missing files",
			err:  "usage: check FILE...",
		},
		{
			name:  "unknown option precedes file reads and arity",
			args:  []string{"--oops", "missing.cfg"},
			err:   "unknown option \"--oops\"",
			files: map[string][]byte{},
		},
		{
			name: "double dash is an unknown option",
			args: []string{"--", "--help"},
			err:  "unknown option \"--\"",
		},
		{
			name: "help is not a check option",
			args: []string{"--help"},
			err:  "unknown option \"--help\"",
		},
		{
			name:   "CQ001 first syntax fault",
			args:   []string{"t.cfg"},
			files:  map[string][]byte{"t.cfg": []byte("[a]\nb=1\n[c\nd=2\n")},
			status: 1,
			stdout: "t.cfg:3: CQ001 bad section header\n1 problem\n",
		},
		{
			name:   "CQ002 first syntax fault",
			args:   []string{"t.cfg"},
			files:  map[string][]byte{"t.cfg": []byte("[a]\n 2c=1\n")},
			status: 1,
			stdout: "t.cfg:2: CQ002 bad variable name\n1 problem\n",
		},
		{
			name:   "CQ003 first syntax fault",
			args:   []string{"t.cfg"},
			files:  map[string][]byte{"t.cfg": []byte("[a]\nb c\n")},
			status: 1,
			stdout: "t.cfg:2: CQ003 expected \"=\" after the variable name\n1 problem\n",
		},
		{
			name:   "CQ004 first syntax fault",
			args:   []string{"t.cfg"},
			files:  map[string][]byte{"t.cfg": []byte("[a]\nb=\"x\n")},
			status: 1,
			stdout: "t.cfg:2: CQ004 unterminated quote\n1 problem\n",
		},
		{
			name:   "CQ005 first syntax fault",
			args:   []string{"t.cfg"},
			files:  map[string][]byte{"t.cfg": []byte("[a]\nb=x\\q\n")},
			status: 1,
			stdout: "t.cfg:2: CQ005 unknown escape \"\\q\"\n1 problem\n",
		},
		{
			name: "syntax suppresses semantics and later files continue",
			args: []string{"a.cfg", "b.cfg"},
			files: map[string][]byte{
				"a.cfg": []byte("[a]\nb=1\nb=1\n[c\n"),
				"b.cfg": []byte("[c]\n"),
			},
			status: 1,
			stdout: "a.cfg:4: CQ001 bad section header\nb.cfg:1: CQ102 empty section \"c\"\n2 problems\n",
		},
		{
			name: "aggregate findings in argument order",
			args: []string{"a.cfg", "b.cfg"},
			files: map[string][]byte{
				"a.cfg": []byte("[a]\nb = 1\nb = 1\n"),
				"b.cfg": []byte("[c]\n"),
			},
			status: 1,
			stdout: "a.cfg:3: CQ101 duplicate entry \"a.b\" (first at line 2)\nb.cfg:1: CQ102 empty section \"c\"\n2 problems\n",
		},
		{
			name: "aggregate findings reverse with argument order",
			args: []string{"b.cfg", "a.cfg"},
			files: map[string][]byte{
				"a.cfg": []byte("[a]\nb = 1\nb = 1\n"),
				"b.cfg": []byte("[c]\n"),
			},
			status: 1,
			stdout: "b.cfg:1: CQ102 empty section \"c\"\na.cfg:3: CQ101 duplicate entry \"a.b\" (first at line 2)\n2 problems\n",
		},
		{
			name: "clean and empty files report ok",
			args: []string{"clean.cfg", "empty.cfg"},
			files: map[string][]byte{
				"clean.cfg": []byte("[a]\nb = 1\n"),
				"empty.cfg": {},
			},
			stdout: "ok\n",
		},
		{
			name: "missing file rolls back earlier findings",
			args: []string{"a.cfg", "missing.cfg"},
			files: map[string][]byte{
				"a.cfg": []byte("[a]\nb=1\nb=1\n"),
			},
			err: "missing.cfg: cannot read the file",
		},
		{
			name: "invalid UTF-8 file rolls back earlier findings",
			args: []string{"a.cfg", "bad.cfg"},
			files: map[string][]byte{
				"a.cfg":   []byte("[a]\nb=1\nb=1\n"),
				"bad.cfg": {0x5b, 0x61, 0x5d, 0x0a, 0x62, 0xff, 0x0a},
			},
			err: "bad.cfg: input is not valid UTF-8",
		},
		{
			name:  "repeated stdin rolls back earlier file and stdin findings",
			args:  []string{"a.cfg", "-", "-"},
			files: map[string][]byte{"a.cfg": []byte("[a]\nb=1\nb=1\n")},
			stdin: []byte("[c]\n"),
			err:   "standard input is named twice",
		},
		{
			name:  "unicode path is displayed as plain text",
			args:  []string{"a.cfg", "café😀.cfg"},
			files: map[string][]byte{"a.cfg": []byte("[a]\nb=1\nb=1\n")},
			err:   "caf??.cfg: cannot read the file",
		},
		{
			name: "syntax report is discarded on later fatal read",
			args: []string{"syntax.cfg", "missing.cfg"},
			files: map[string][]byte{
				"syntax.cfg": []byte("[a]\nb c\n"),
			},
			err: "missing.cfg: cannot read the file",
		},
		{
			name:  "invalid UTF-8 stdin bytes roll back earlier findings",
			args:  []string{"a.cfg", "-"},
			files: map[string][]byte{"a.cfg": []byte("[a]\nb=1\nb=1\n")},
			stdin: []byte{0x5b, 0x61, 0x5d, 0x0a, 0x62, 0xff, 0x0a},
			err:   "<stdin>: input is not valid UTF-8",
		},
		{
			name:   "stdin uses the stdin display name",
			args:   []string{"-"},
			stdin:  []byte("[c]\n"),
			status: 1,
			stdout: "<stdin>:1: CQ102 empty section \"c\"\n1 problem\n",
		},
		{
			name: "first fatal input error wins",
			args: []string{"missing.cfg", "bad.cfg"},
			files: map[string][]byte{
				"bad.cfg": {0x5b, 0x61, 0x5d, 0x0a, 0x62, 0xff, 0x0a},
			},
			err: "missing.cfg: cannot read the file",
		},
		{
			name:   "BOM and CRLF are preserved while trailing CR is reported",
			args:   []string{"t.cfg"},
			files:  map[string][]byte{"t.cfg": []byte("\ufeff[a]\r\nb=1\r")},
			status: 1,
			stdout: "t.cfg:2: CQ105 bare CR at column 4\n1 problem\n",
		},
		{
			name:   "bare CR in comment is reported at its raw column",
			args:   []string{"t.cfg"},
			files:  map[string][]byte{"t.cfg": []byte("[a]\nb=1\n# x\ry\n")},
			status: 1,
			stdout: "t.cfg:3: CQ105 bare CR at column 4\n1 problem\n",
		},
		{
			name:   "CRLF has no bare CR findings",
			args:   []string{"t.cfg"},
			files:  map[string][]byte{"t.cfg": []byte("\ufeff[a]\r\nb = 1\r\n")},
			stdout: "ok\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			io := &base.MemIO{Files: test.files, Stdin: test.stdin}
			status, err := Run(test.args, io)

			if status != test.status {
				t.Errorf("Run(%q) status = %d, want %d", test.args, status, test.status)
			}
			if test.err == "" {
				if err != nil {
					t.Errorf("Run(%q) error = %v, want nil", test.args, err)
				}
			} else if err == nil || err.Error() != test.err {
				t.Errorf("Run(%q) error = %v, want %q", test.args, err, test.err)
			}
			if io.Stdout != test.stdout {
				t.Errorf("Run(%q) stdout = %q, want %q", test.args, io.Stdout, test.stdout)
			}
			if io.Stderr != "" {
				t.Errorf("Run(%q) stderr = %q, want empty", test.args, io.Stderr)
			}
		})
	}
}
