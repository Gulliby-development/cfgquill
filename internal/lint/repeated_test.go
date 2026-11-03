package lint

import (
	"testing"

	"github.com/Gulliby-development/cfgquill/internal/base"
)

func TestRepeatedSectionFindings(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		status int
		stdout string
	}{
		{"one header is fine", "[a]\nb = 1\n", 0, "ok\n"},
		{"same section twice", "[a]\nb = 1\n[a]\nc = 2\n", 1, "t.cfg:3: CQ106 repeated section \"a\" (first at line 1)\n1 problem\n"},
		{"names are compared in stored form", "[Core]\nb = 1\n[CORE]\nc = 2\n", 1, "t.cfg:3: CQ106 repeated section \"core\" (first at line 1)\n1 problem\n"},
		{"subsections are case sensitive", "[a \"X\"]\nb = 1\n[a \"x\"]\nc = 2\n", 0, "ok\n"},
		{"same subsection twice", "[a \"X\"]\nb = 1\n[a \"X\"]\nc = 2\n", 1, "t.cfg:3: CQ106 repeated section \"a.X\" (first at line 1)\n1 problem\n"},
		{"section and subsection are different headers", "[a]\nb = 1\n[a \"x\"]\nc = 2\n", 0, "ok\n"},
		{"third header points at the first", "[a]\nb = 1\n[c]\nd = 1\n[a]\ne = 1\n[a]\nf = 1\n", 1, "t.cfg:5: CQ106 repeated section \"a\" (first at line 1)\nt.cfg:7: CQ106 repeated section \"a\" (first at line 1)\n2 problems\n"},
		{"empty repeated section has two findings", "[a]\nb = 1\n[a]\n", 1, "t.cfg:3: CQ102 empty section \"a\"\nt.cfg:3: CQ106 repeated section \"a\" (first at line 1)\n2 problems\n"},
		{"empty first section", "[a]\n[a]\nb = 1\n", 1, "t.cfg:1: CQ102 empty section \"a\"\nt.cfg:2: CQ106 repeated section \"a\" (first at line 1)\n2 problems\n"},
		{"syntax error hides the finding", "[a]\nb = 1\n[a]\nc\n[d\n", 1, "t.cfg:5: CQ001 bad section header\n1 problem\n"},
		{"with a duplicate entry", "[a]\nb = 1\n[a]\nb = 1\n", 1, "t.cfg:3: CQ106 repeated section \"a\" (first at line 1)\nt.cfg:4: CQ101 duplicate entry \"a.b\" (first at line 2)\n2 problems\n"},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			memory := &base.MemIO{Files: map[string][]byte{"t.cfg": []byte(test.text)}}
			status, err := Run([]string{"t.cfg"}, memory)
			if err != nil {
				t.Fatalf("Run error = %v, want nil", err)
			}
			if status != test.status || memory.Stdout != test.stdout || memory.Stderr != "" {
				t.Errorf("status = %d, stdout = %q, stderr = %q, want %d, %q and empty stderr", status, memory.Stdout, memory.Stderr, test.status, test.stdout)
			}
		})
	}
}

func TestRepeatedSectionIsPerFile(t *testing.T) {
	memory := &base.MemIO{Files: map[string][]byte{
		"a.cfg": []byte("[a]\nb = 1\n"),
		"b.cfg": []byte("[a]\nc = 1\n[a]\nd = 1\n"),
	}}
	status, err := Run([]string{"a.cfg", "b.cfg"}, memory)
	want := "b.cfg:3: CQ106 repeated section \"a\" (first at line 1)\n1 problem\n"
	if err != nil || status != 1 || memory.Stdout != want {
		t.Errorf("status = %d, err = %v, stdout = %q, want 1, nil and %q", status, err, memory.Stdout, want)
	}
}
