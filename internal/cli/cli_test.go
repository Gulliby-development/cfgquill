package cli

// BEGIN acceptance rows imports: host-written from the audited spec table. Edit acceptanceActual, not this block.
import (
	acceptancejson "encoding/json"
	acceptancereflect "reflect"
	acceptancetesting "testing"
)

// END acceptance rows imports

import (
	"bytes"
	"errors"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/Gulliby-development/cfgquill/internal/base"
)

// BEGIN acceptance rows: host-written from the audited spec table. Edit acceptanceActual, not this block.
func TestAcceptanceRowsCli(t *acceptancetesting.T) {
	var rows []struct {
		ID       string `json:"id"`
		Input    any    `json:"input"`
		Expected any    `json:"expected"`
	}
	if err := acceptancejson.Unmarshal([]byte("[{\"id\":\"AC-001:1\",\"input\":{\"argv\":[],\"files\":{},\"stdin\":\"\"},\"expected\":{\"status\":2,\"stderr\":\"cfgquill: missing command (list, sections, get, bool, unquote, quote or check)\\n\",\"stdout\":\"\"}},{\"id\":\"AC-001:2\",\"input\":{\"argv\":[\"--help\"],\"files\":{},\"stdin\":\"\"},\"expected\":{\"status\":0,\"stderr\":\"\",\"stdout\":\"usage: cfgquill list FILE\\nusage: cfgquill sections FILE\\nusage: cfgquill get [--all] KEY FILE\\nusage: cfgquill bool KEY FILE\\nusage: cfgquill unquote TEXT\\nusage: cfgquill quote TEXT\\nusage: cfgquill check FILE...\\n\"}},{\"id\":\"AC-001:3\",\"input\":{\"argv\":[\"--help\",\"list\"],\"files\":{},\"stdin\":\"\"},\"expected\":{\"status\":2,\"stderr\":\"cfgquill: unknown command \\\"--help\\\"\\n\",\"stdout\":\"\"}},{\"id\":\"AC-001:4\",\"input\":{\"argv\":[\"frob\",\"--oops\"],\"files\":{},\"stdin\":\"\"},\"expected\":{\"status\":2,\"stderr\":\"cfgquill: unknown command \\\"frob\\\"\\n\",\"stdout\":\"\"}},{\"id\":\"AC-001:5\",\"input\":{\"argv\":[\"😀\\n\"],\"files\":{},\"stdin\":\"\"},\"expected\":{\"status\":2,\"stderr\":\"cfgquill: unknown command \\\"??\\\"\\n\",\"stdout\":\"\"}},{\"id\":\"AC-001:6\",\"input\":{\"argv\":[\"list\",\"--help\"],\"files\":{},\"stdin\":\"\"},\"expected\":{\"status\":2,\"stderr\":\"cfgquill: unknown option \\\"--help\\\"\\n\",\"stdout\":\"\"}},{\"id\":\"AC-001:7\",\"input\":{\"argv\":[\"list\",\"--oops\",\"--bad\"],\"files\":{},\"stdin\":\"\"},\"expected\":{\"status\":2,\"stderr\":\"cfgquill: unknown option \\\"--oops\\\"\\n\",\"stdout\":\"\"}},{\"id\":\"AC-001:8\",\"input\":{\"argv\":[\"list\"],\"files\":{},\"stdin\":\"\"},\"expected\":{\"status\":2,\"stderr\":\"cfgquill: usage: list FILE\\n\",\"stdout\":\"\"}},{\"id\":\"AC-001:9\",\"input\":{\"argv\":[\"get\",\"--all\",\"a.b\",\"t.cfg\",\"u.cfg\"],\"files\":{},\"stdin\":\"\"},\"expected\":{\"status\":2,\"stderr\":\"cfgquill: usage: get [--all] KEY FILE\\n\",\"stdout\":\"\"}},{\"id\":\"AC-001:10\",\"input\":{\"argv\":[\"get\",\".b\",\"missing.cfg\"],\"files\":{},\"stdin\":\"\"},\"expected\":{\"status\":2,\"stderr\":\"cfgquill: invalid key \\\".b\\\"\\n\",\"stdout\":\"\"}},{\"id\":\"AC-001:11\",\"input\":{\"argv\":[\"list\",\"missing.cfg\"],\"files\":{},\"stdin\":\"\"},\"expected\":{\"status\":2,\"stderr\":\"cfgquill: missing.cfg: cannot read the file\\n\",\"stdout\":\"\"}},{\"id\":\"AC-001:12\",\"input\":{\"argv\":[\"list\",\"café😀.cfg\"],\"files\":{},\"stdin\":\"\"},\"expected\":{\"status\":2,\"stderr\":\"cfgquill: caf??.cfg: cannot read the file\\n\",\"stdout\":\"\"}},{\"id\":\"AC-001:13\",\"input\":{\"argv\":[\"list\",\"t.cfg\"],\"files\":{\"t.cfg\":\"[a]\\nb = 1\\n\"},\"stdin\":\"\"},\"expected\":{\"status\":0,\"stderr\":\"\",\"stdout\":\"LINE  KEY  VALUE\\n2     a.b  \\\"1\\\"\\n\"}},{\"id\":\"AC-001:14\",\"input\":{\"argv\":[\"sections\",\"t.cfg\"],\"files\":{\"t.cfg\":\"[a]\\nb = 1\\n\"},\"stdin\":\"\"},\"expected\":{\"status\":0,\"stderr\":\"\",\"stdout\":\"LINE  SECTION  SUBSECTION  ENTRIES\\n1     a        -           1\\n\"}},{\"id\":\"AC-001:15\",\"input\":{\"argv\":[\"get\",\"a.b\",\"t.cfg\"],\"files\":{\"t.cfg\":\"[a]\\nb = 1\\nb = 2\\n\"},\"stdin\":\"\"},\"expected\":{\"status\":0,\"stderr\":\"\",\"stdout\":\"2\\n\"}},{\"id\":\"AC-001:16\",\"input\":{\"argv\":[\"get\",\"--all\",\"--all\",\"a.b\",\"t.cfg\"],\"files\":{\"t.cfg\":\"[a]\\nb = 1\\nb = 2\\n\"},\"stdin\":\"\"},\"expected\":{\"status\":0,\"stderr\":\"\",\"stdout\":\"1\\n2\\n\"}},{\"id\":\"AC-001:17\",\"input\":{\"argv\":[\"get\",\"a.c\",\"t.cfg\"],\"files\":{\"t.cfg\":\"[a]\\nb = 1\\n\"},\"stdin\":\"\"},\"expected\":{\"status\":1,\"stderr\":\"\",\"stdout\":\"\"}},{\"id\":\"AC-001:18\",\"input\":{\"argv\":[\"bool\",\"a.b\",\"t.cfg\"],\"files\":{\"t.cfg\":\"[a]\\nb\\n\"},\"stdin\":\"\"},\"expected\":{\"status\":0,\"stderr\":\"\",\"stdout\":\"true\\n\"}},{\"id\":\"AC-001:19\",\"input\":{\"argv\":[\"bool\",\"a.c\",\"t.cfg\"],\"files\":{\"t.cfg\":\"[a]\\nb\\n\"},\"stdin\":\"\"},\"expected\":{\"status\":1,\"stderr\":\"\",\"stdout\":\"\"}},{\"id\":\"AC-001:20\",\"input\":{\"argv\":[\"bool\",\"a.b\",\"t.cfg\"],\"files\":{\"t.cfg\":\"[a]\\nb = maybe\\n[\\n\"},\"stdin\":\"\"},\"expected\":{\"status\":2,\"stderr\":\"cfgquill: t.cfg:3: CQ001 bad section header\\n\",\"stdout\":\"\"}},{\"id\":\"AC-001:21\",\"input\":{\"argv\":[\"unquote\",\"a \\\"\\\" b\"],\"files\":{},\"stdin\":\"\"},\"expected\":{\"status\":0,\"stderr\":\"\",\"stdout\":\"a  b\"}},{\"id\":\"AC-001:22\",\"input\":{\"argv\":[\"unquote\",\"a\\\"b\"],\"files\":{},\"stdin\":\"\"},\"expected\":{\"status\":2,\"stderr\":\"cfgquill: unquote: CQ004 unterminated quote\\n\",\"stdout\":\"\"}},{\"id\":\"AC-001:23\",\"input\":{\"argv\":[\"quote\",\" a\"],\"files\":{},\"stdin\":\"\"},\"expected\":{\"status\":0,\"stderr\":\"\",\"stdout\":\"\\\" a\\\"\"}},{\"id\":\"AC-001:24\",\"input\":{\"argv\":[\"check\",\"t.cfg\"],\"files\":{\"t.cfg\":\"[a]\\nb = 1\\n\"},\"stdin\":\"\"},\"expected\":{\"status\":0,\"stderr\":\"\",\"stdout\":\"ok\\n\"}},{\"id\":\"AC-001:25\",\"input\":{\"argv\":[\"check\",\"t.cfg\"],\"files\":{\"t.cfg\":\"[a]\\n\"},\"stdin\":\"\"},\"expected\":{\"status\":1,\"stderr\":\"\",\"stdout\":\"t.cfg:1: CQ102 empty section \\\"a\\\"\\n1 problem\\n\"}},{\"id\":\"AC-001:26\",\"input\":{\"argv\":[\"check\",\"t.cfg\",\"missing.cfg\"],\"files\":{\"t.cfg\":\"[a]\\n\"},\"stdin\":\"\"},\"expected\":{\"status\":2,\"stderr\":\"cfgquill: missing.cfg: cannot read the file\\n\",\"stdout\":\"\"}},{\"id\":\"AC-001:27\",\"input\":{\"argv\":[\"check\",\"-\",\"-\"],\"files\":{},\"stdin\":\"\"},\"expected\":{\"status\":2,\"stderr\":\"cfgquill: standard input is named twice\\n\",\"stdout\":\"\"}}]"), &rows); err != nil {
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
	fields, ok := input.(map[string]any)
	if !ok {
		return nil, errors.New("invalid acceptance input")
	}

	argvValues, ok := fields["argv"].([]any)
	if !ok {
		return nil, errors.New("invalid acceptance argv")
	}
	argv := make([]string, len(argvValues))
	for index, value := range argvValues {
		argv[index], ok = value.(string)
		if !ok {
			return nil, errors.New("invalid acceptance argument")
		}
	}

	filesValue, ok := fields["files"].(map[string]any)
	if !ok {
		return nil, errors.New("invalid acceptance files")
	}
	files := make(map[string][]byte, len(filesValue))
	for path, value := range filesValue {
		text, ok := value.(string)
		if !ok {
			return nil, errors.New("invalid acceptance file")
		}
		files[path] = []byte(text)
	}
	stdin, ok := fields["stdin"].(string)
	if !ok {
		return nil, errors.New("invalid acceptance stdin")
	}

	memory := &base.MemIO{Stdin: []byte(stdin), Files: files}
	status := Run(argv, memory)
	return map[string]any{
		"status": status,
		"stdout": memory.Stdout,
		"stderr": memory.Stderr,
	}, nil
}

type countingIO struct {
	*base.MemIO
	outCalls int
	errCalls int
}

func (io *countingIO) Out(text string) {
	io.outCalls++
	io.MemIO.Out(text)
}

func (io *countingIO) Err(text string) {
	io.errCalls++
	io.MemIO.Err(text)
}

func TestCLIUsageAndHelpCall(t *testing.T) {
	want := "usage: cfgquill list FILE\n" +
		"usage: cfgquill sections FILE\n" +
		"usage: cfgquill get [--all] KEY FILE\n" +
		"usage: cfgquill bool KEY FILE\n" +
		"usage: cfgquill unquote TEXT\n" +
		"usage: cfgquill quote TEXT\n" +
		"usage: cfgquill check FILE...\n"
	if Usage != want {
		t.Fatalf("Usage = %q, want %q", Usage, want)
	}

	memory := &countingIO{MemIO: &base.MemIO{}}
	if status := Run([]string{"--help"}, memory); status != 0 {
		t.Fatalf("status = %d, want 0", status)
	}
	if memory.Stdout != want || memory.outCalls != 1 {
		t.Errorf("help output = %q with %d Out calls, want exact usage with one call", memory.Stdout, memory.outCalls)
	}
	if memory.Stderr != "" || memory.errCalls != 0 {
		t.Errorf("help stderr = %q with %d Err calls, want empty with no calls", memory.Stderr, memory.errCalls)
	}
}

func TestCLIErrorIsReportedOnce(t *testing.T) {
	memory := &countingIO{MemIO: &base.MemIO{}}
	if status := Run([]string{"list", "--oops"}, memory); status != 2 {
		t.Fatalf("status = %d, want 2", status)
	}
	if memory.Stdout != "" || memory.outCalls != 0 {
		t.Errorf("stdout = %q with %d Out calls, want empty with no calls", memory.Stdout, memory.outCalls)
	}
	if memory.Stderr != "cfgquill: unknown option \"--oops\"\n" || memory.errCalls != 1 {
		t.Errorf("stderr = %q with %d Err calls, want one prefixed diagnostic", memory.Stderr, memory.errCalls)
	}
}

func TestDispatchAcceptanceRowsExplicit(t *testing.T) {
	rows := []struct {
		input  string
		status int
		stdout string
		stderr string
	}{
		{`{"argv":["😀\n"],"files":{},"stdin":""}`, 2, "", "cfgquill: unknown command \"??\"\n"},
		{`{"argv":["list","t.cfg"],"files":{"t.cfg":"[a]\nb = 1\n"},"stdin":""}`, 0, "LINE  KEY  VALUE\n2     a.b  \"1\"\n", ""},
		{`{"argv":["sections","t.cfg"],"files":{"t.cfg":"[a]\nb = 1\n"},"stdin":""}`, 0, "LINE  SECTION  SUBSECTION  ENTRIES\n1     a        -           1\n", ""},
		{`{"argv":["get","a.b","t.cfg"],"files":{"t.cfg":"[a]\nb = 1\nb = 2\n"},"stdin":""}`, 0, "2\n", ""},
		{`{"argv":["get","--all","--all","a.b","t.cfg"],"files":{"t.cfg":"[a]\nb = 1\nb = 2\n"},"stdin":""}`, 0, "1\n2\n", ""},
		{`{"argv":["get","a.c","t.cfg"],"files":{"t.cfg":"[a]\nb = 1\n"},"stdin":""}`, 1, "", ""},
		{`{"argv":["bool","a.b","t.cfg"],"files":{"t.cfg":"[a]\nb\n"},"stdin":""}`, 0, "true\n", ""},
		{`{"argv":["bool","a.c","t.cfg"],"files":{"t.cfg":"[a]\nb\n"},"stdin":""}`, 1, "", ""},
		{`{"argv":["bool","a.b","t.cfg"],"files":{"t.cfg":"[a]\nb = maybe\n[\n"},"stdin":""}`, 2, "", "cfgquill: t.cfg:3: CQ001 bad section header\n"},
		{`{"argv":["check","t.cfg"],"files":{"t.cfg":"[a]\nb = 1\n"},"stdin":""}`, 0, "ok\n", ""},
	}
	for _, row := range rows {
		var input any
		if err := acceptancejson.Unmarshal([]byte(row.input), &input); err != nil {
			t.Fatal(err)
		}
		result, err := acceptanceActual(input)
		if err != nil {
			t.Fatal(err)
		}
		actual := result.(map[string]any)
		if actual["status"] != row.status {
			t.Errorf("status = %v, want %d for %s", actual["status"], row.status, row.input)
		}
		if actual["stdout"] != row.stdout {
			t.Errorf("stdout = %q, want %q for %s", actual["stdout"], row.stdout, row.input)
		}
		if actual["stderr"] != row.stderr {
			t.Errorf("stderr = %q, want %q for %s", actual["stderr"], row.stderr, row.input)
		}
	}
}

func TestExecutableRejectsInvalidStdin(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "cfgquill")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}

	build := exec.Command("go", "build", "-o", binary, "../..")
	var buildOutput bytes.Buffer
	build.Stdout = &buildOutput
	build.Stderr = &buildOutput
	if err := build.Run(); err != nil {
		t.Fatalf("go build failed: %v\n%s", err, buildOutput.String())
	}

	program := exec.Command(binary, "list", "-")
	program.Stdin = bytes.NewReader([]byte{0x5B, 0x61, 0x5D, 0x0A, 0x62, 0xFF, 0x0A})
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	program.Stdout = &stdout
	program.Stderr = &stderr
	err := program.Run()
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 2 {
		t.Fatalf("exit error = %v, want exit code 2", err)
	}
	if stdout.String() != "" {
		t.Errorf("stdout = %q, want empty", stdout.String())
	}
	if stderr.String() != "cfgquill: <stdin>: input is not valid UTF-8\n" {
		t.Errorf("stderr = %q, want exact invalid UTF-8 diagnostic", stderr.String())
	}
}
