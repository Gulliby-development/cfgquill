package base

// BEGIN acceptance rows imports: host-written from the audited spec table. Edit acceptanceActual, not this block.
import (
	acceptancejson "encoding/json"
	acceptancereflect "reflect"
	acceptancetesting "testing"
)

// END acceptance rows imports

import "testing"

// BEGIN acceptance rows: host-written from the audited spec table. Edit acceptanceActual, not this block.
func TestAcceptanceRowsBase(t *acceptancetesting.T) {
	var rows []struct {
		ID       string `json:"id"`
		Input    any    `json:"input"`
		Expected any    `json:"expected"`
	}
	if err := acceptancejson.Unmarshal([]byte("[{\"id\":\"AC-001:1\",\"input\":\"\",\"expected\":[\"\",\"\\\"\\\"\"]},{\"id\":\"AC-001:2\",\"input\":\"a\\\"b\",\"expected\":[\"a\\\"b\",\"\\\"a\\\"b\\\"\"]},{\"id\":\"AC-001:3\",\"input\":\"café😀\",\"expected\":[\"caf??\",\"\\\"caf??\\\"\"]},{\"id\":\"AC-001:4\",\"input\":\" \\t\\n\\r\",\"expected\":[\" ????\",\"\\\" ????\\\"\"]},{\"id\":\"AC-001:5\",\"input\":\"~\",\"expected\":[\"~\",\"\\\"~\\\"\"]}]"), &rows); err != nil {
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

// acceptanceActual adapts the public text utilities to the host's row harness.
func acceptanceActual(input any) (any, error) {
	text := input.(string)
	return []string{Plain(text), Show(text)}, nil
}

func TestTextUtilities(t *testing.T) {
	pluralCases := []struct {
		n    int
		word string
		want string
	}{
		{0, "problem", "0 problems"},
		{1, "problem", "1 problem"},
		{2, "problem", "2 problems"},
		{-1, "problem", "-1 problems"},
	}
	for _, test := range pluralCases {
		if got := Plural(test.n, test.word); got != test.want {
			t.Errorf("Plural(%d, %q) = %q, want %q", test.n, test.word, got, test.want)
		}
	}

	tableCases := []struct {
		header []string
		rows   [][]string
		want   []string
	}{
		{
			[]string{"A", "B"},
			[][]string{{"😀", "x"}, {"aa", "yz"}},
			[]string{"A   B", "😀   x", "aa  yz"},
		},
		{
			[]string{"LINE", "VALUE"},
			nil,
			[]string{"LINE  VALUE"},
		},
		{
			[]string{"A", "B"},
			[][]string{{"x", ""}},
			[]string{"A  B", "x"},
		},
	}
	for _, test := range tableCases {
		if got := Table(test.header, test.rows); !acceptancereflect.DeepEqual(got, test.want) {
			t.Errorf("Table(%q, %q) = %#v, want %#v", test.header, test.rows, got, test.want)
		}
	}
}

func TestOptions(t *testing.T) {
	parseCases := []struct {
		argv  []string
		specs []string
		want  Parsed
		err   string
	}{
		{
			[]string{"--all", "a.b", "--all", "-"},
			[]string{"--all"},
			Parsed{Flags: map[string]bool{"--all": true}, Args: []string{"a.b", "-"}},
			"",
		},
		{
			[]string{"-x", ""},
			nil,
			Parsed{Flags: map[string]bool{}, Args: []string{"-x", ""}},
			"",
		},
		{[]string{"--"}, []string{"--all"}, Parsed{}, "unknown option \"--\""},
		{[]string{"--help"}, nil, Parsed{}, "unknown option \"--help\""},
		{[]string{"--all=1"}, []string{"--all"}, Parsed{}, "unknown option \"--all=1\""},
		{[]string{"--oops", "--help"}, nil, Parsed{}, "unknown option \"--oops\""},
	}
	for _, test := range parseCases {
		got, err := Parse(test.argv, test.specs)
		if test.err != "" {
			if err == nil || err.Error() != test.err {
				t.Errorf("Parse(%q, %q) error = %v, want %q", test.argv, test.specs, err, test.err)
			}
			continue
		}
		if err != nil {
			t.Errorf("Parse(%q, %q) error = %v", test.argv, test.specs, err)
			continue
		}
		if !acceptancereflect.DeepEqual(got, test.want) {
			t.Errorf("Parse(%q, %q) = %#v, want %#v", test.argv, test.specs, got, test.want)
		}
	}

	needCases := []struct {
		words string
		args  []string
		min   int
		max   int
		want  string
	}{
		{"list FILE", nil, 1, 1, "usage: list FILE"},
		{"list FILE", []string{"t.cfg"}, 1, 1, ""},
		{"list FILE", []string{"t.cfg", "u.cfg"}, 1, 1, "usage: list FILE"},
		{"check FILE...", []string{"a", "b"}, 1, -1, ""},
		{"check FILE...", nil, 0, -1, ""},
	}
	for _, test := range needCases {
		err := Need(test.words, test.args, test.min, test.max)
		if test.want == "" {
			if err != nil {
				t.Errorf("Need(%q, %q, %d, %d) error = %v, want nil", test.words, test.args, test.min, test.max, err)
			}
		} else if err == nil || err.Error() != test.want {
			t.Errorf("Need(%q, %q, %d, %d) error = %v, want %q", test.words, test.args, test.min, test.max, err, test.want)
		}
	}
}

func TestInputsReadFiles(t *testing.T) {
	badBytes := []struct {
		name string
		data []byte
		want string
	}{
		{"bad-utf8", []byte{0x5b, 0x61, 0x5d, 0x0a, 0x62, 0xff, 0x0a}, "bad-utf8: input is not valid UTF-8"},
		{"continuation", []byte{0x80}, "continuation: input is not valid UTF-8"},
		{"overlong", []byte{0xc0, 0xaf}, "overlong: input is not valid UTF-8"},
		{"surrogate", []byte{0xed, 0xa0, 0x80}, "surrogate: input is not valid UTF-8"},
		{"truncated", []byte{0xe2, 0x82}, "truncated: input is not valid UTF-8"},
		{"too-large", []byte{0xf4, 0x90, 0x80, 0x80}, "too-large: input is not valid UTF-8"},
	}
	files := map[string][]byte{
		"café😀.cfg": []byte("\ufeffa\r\n\u0085\u2028\u2029"),
		"empty":     {},
		"valid":     []byte("�😀"),
	}
	for _, test := range badBytes {
		files[test.name] = test.data
	}
	io := &MemIO{Files: files}
	inputs := NewInputs(io)

	got, err := inputs.Read("café😀.cfg")
	if err != nil {
		t.Fatalf("Read valid Unicode file: %v", err)
	}
	if got.Display != "caf??.cfg" || got.Text != "\ufeffa\r\n\u0085\u2028\u2029" {
		t.Errorf("Read Unicode file = %#v, want unchanged text and display caf??.cfg", got)
	}
	got, err = inputs.Read("empty")
	if err != nil || got.Display != "empty" || got.Text != "" {
		t.Errorf("Read empty file = %#v, %v", got, err)
	}
	got, err = inputs.Read("valid")
	if err != nil || got.Text != "�😀" {
		t.Errorf("Read valid replacement and astral runes = %#v, %v", got, err)
	}
	for _, test := range badBytes {
		_, err := inputs.Read(test.name)
		if err == nil || err.Error() != test.want {
			t.Errorf("Read(%q) error = %v, want invalid UTF-8", test.name, err)
		}
	}
	missingCases := []struct {
		path string
		want string
	}{
		{"missing.cfg", "missing.cfg: cannot read the file"},
		{"", ": cannot read the file"},
	}
	for _, test := range missingCases {
		_, err := inputs.Read(test.path)
		if err == nil || err.Error() != test.want {
			t.Errorf("Read(%q) error = %v, want %q", test.path, err, test.want)
		}
	}
	if got, err = inputs.Read("empty"); err != nil || got.Text != "" {
		t.Errorf("repeated file read = %#v, %v", got, err)
	}
}

func TestInputsReadStdinAndMemIO(t *testing.T) {
	io := &MemIO{Stdin: []byte("\ufeffx\r\ny")}
	inputs := NewInputs(io)
	got, err := inputs.Read("-")
	if err != nil || got.Display != "<stdin>" || got.Text != "\ufeffx\r\ny" {
		t.Errorf("Read stdin = %#v, %v", got, err)
	}
	_, err = inputs.Read("-")
	if err == nil || err.Error() != "standard input is named twice" {
		t.Errorf("second Read stdin error = %v", err)
	}

	invalidStdin := &MemIO{Stdin: []byte{0x5b, 0x61, 0x5d, 0x0a, 0x62, 0xff, 0x0a}}
	invalidInputs := NewInputs(invalidStdin)
	_, err = invalidInputs.Read("-")
	if err == nil || err.Error() != "<stdin>: input is not valid UTF-8" {
		t.Errorf("invalid stdin error = %v", err)
	}
	_, err = invalidInputs.Read("-")
	if err == nil || err.Error() != "standard input is named twice" {
		t.Errorf("Read stdin after invalid bytes error = %v", err)
	}

	failed := &failingStdinIO{MemIO: &MemIO{}}
	failedInputs := NewInputs(failed)
	_, err = failedInputs.Read("-")
	if err == nil || err.Error() != "<stdin>: cannot read the file" {
		t.Errorf("failed stdin read error = %v", err)
	}
	_, err = failedInputs.Read("-")
	if err == nil || err.Error() != "standard input is named twice" || failed.calls != 1 {
		t.Errorf("second failed stdin read = %v, calls = %d, want duplicate and one read", err, failed.calls)
	}

	memory := &MemIO{Files: map[string][]byte{"present-empty": {}}}
	data, err := memory.ReadFile("present-empty")
	if err != nil || len(data) != 0 {
		t.Errorf("ReadFile(present-empty) = %v, %v, want empty bytes and nil", data, err)
	}
	if _, err := memory.ReadFile("absent"); err == nil {
		t.Error("ReadFile(absent) succeeded, want an error")
	}
	memory.Out("a")
	memory.Out("b")
	memory.Err("x")
	memory.Err("y")
	if memory.Stdout != "ab" || memory.Stderr != "xy" {
		t.Errorf("MemIO output = %q and %q, want ab and xy", memory.Stdout, memory.Stderr)
	}
}

type failingStdinIO struct {
	*MemIO
	calls int
}

func (f *failingStdinIO) ReadStdin() ([]byte, error) {
	f.calls++
	_, err := f.MemIO.ReadFile("missing")
	return nil, err
}

func TestCursor(t *testing.T) {
	cursorCases := []struct {
		text         string
		initialAtEnd bool
		runes        []rune
		lines        []int
		ends         []bool
	}{
		{
			"a\nb",
			false,
			[]rune{'a', '\n', 'b', '\n', EOF, EOF},
			[]int{1, 1, 2, 2, 2, 2},
			[]bool{false, false, true, true, true, true},
		},
		{
			"a\n",
			false,
			[]rune{'a', '\n', '\n', EOF},
			[]int{1, 1, 2, 2},
			[]bool{false, true, true, true},
		},
		{
			"",
			true,
			[]rune{'\n', EOF, EOF},
			[]int{1, 1, 1},
			[]bool{true, true, true},
		},
		{
			"😀\r\nx\ry\u0085\u2028\u2029",
			false,
			[]rune{'😀', '\n', 'x', '\r', 'y', '\u0085', '\u2028', '\u2029', '\n', EOF},
			[]int{1, 1, 2, 2, 2, 2, 2, 2, 2, 2},
			[]bool{false, false, false, false, false, false, false, true, true, true},
		},
	}
	for _, test := range cursorCases {
		cursor := NewCursor(test.text)
		if cursor.Line() != 1 {
			t.Errorf("Line before reading %q = %d, want 1", test.text, cursor.Line())
		}
		if cursor.AtEnd() != test.initialAtEnd {
			t.Errorf("AtEnd before reading %q = %t", test.text, cursor.AtEnd())
		}
		for i, wantRune := range test.runes {
			if got := cursor.Next(); got != wantRune {
				t.Errorf("Next for %q at %d = %U, want %U", test.text, i, got, wantRune)
			}
			if got := cursor.Line(); got != test.lines[i] {
				t.Errorf("Line for %q at %d = %d, want %d", test.text, i, got, test.lines[i])
			}
			if got := cursor.AtEnd(); got != test.ends[i] {
				t.Errorf("AtEnd for %q at %d = %t, want %t", test.text, i, got, test.ends[i])
			}
		}
	}
	cursor := NewCursor("a\nb")
	cursor.Next()
	cursor.Next()
	cursor.Next()
	fault := cursor.Fault("CQ005", "\\q")
	if fault.Code != "CQ005" || fault.Line != 2 || fault.Arg != "\\q" || fault.Error() != "CQ005" {
		t.Errorf("Fault = %#v / %q, want CQ005 at line 2 with argument \\q", fault, fault.Error())
	}
}

func TestASCIIPredicates(t *testing.T) {
	predicateCases := []struct {
		rune    rune
		space   bool
		alpha   bool
		alnum   bool
		keyChar bool
	}{
		{' ', true, false, false, false},
		{'\t', true, false, false, false},
		{'\n', true, false, false, false},
		{'\r', true, false, false, false},
		{'A', false, true, true, true},
		{'z', false, true, true, true},
		{'0', false, false, true, true},
		{'-', false, false, false, true},
		{'_', false, false, false, false},
		{' ', false, false, false, false},
		{'', false, false, false, false},
		{' ', false, false, false, false},
		{' ', false, false, false, false},
		{'　', false, false, false, false},
		{'​', false, false, false, false},
		{'\ufeff', false, false, false, false},
		{'\u000b', false, false, false, false},
		{'\f', false, false, false, false},
		{'é', false, false, false, false},
		{EOF, false, false, false, false},
	}
	for _, test := range predicateCases {
		if got := IsSpace(test.rune); got != test.space {
			t.Errorf("IsSpace(%U) = %t, want %t", test.rune, got, test.space)
		}
		if got := IsAlpha(test.rune); got != test.alpha {
			t.Errorf("IsAlpha(%U) = %t, want %t", test.rune, got, test.alpha)
		}
		if got := IsAlnum(test.rune); got != test.alnum {
			t.Errorf("IsAlnum(%U) = %t, want %t", test.rune, got, test.alnum)
		}
		if got := IsKeyChar(test.rune); got != test.keyChar {
			t.Errorf("IsKeyChar(%U) = %t, want %t", test.rune, got, test.keyChar)
		}
	}
}

func TestLower(t *testing.T) {
	lowerCases := []struct {
		input string
		want  string
	}{
		{"AZaz09-É😀", "azaz09-É😀"},
		{"", ""},
	}
	for _, test := range lowerCases {
		if got := Lower(test.input); got != test.want {
			t.Errorf("Lower(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}
