package cmd

// BEGIN acceptance rows imports: host-written from the audited spec table. Edit acceptanceActual, not this block.
import (
	acceptancejson "encoding/json"
	acceptancereflect "reflect"
	acceptancetesting "testing"
)

// END acceptance rows imports

import (
	"errors"

	"github.com/Gulliby-development/cfgquill/internal/base"
)

// BEGIN acceptance rows: host-written from the audited spec table. Edit acceptanceActual, not this block.
func TestAcceptanceRowsText(t *acceptancetesting.T) {
	var rows []struct {
		ID       string `json:"id"`
		Input    any    `json:"input"`
		Expected any    `json:"expected"`
	}
	if err := acceptancejson.Unmarshal([]byte("[{\"id\":\"AC-001:1\",\"input\":[\"unquote\",[]],\"expected\":\"usage: unquote TEXT\"},{\"id\":\"AC-001:2\",\"input\":[\"quote\",[]],\"expected\":\"usage: quote TEXT\"},{\"id\":\"AC-001:3\",\"input\":[\"unquote\",[\"a\",\"b\"]],\"expected\":\"usage: unquote TEXT\"},{\"id\":\"AC-001:4\",\"input\":[\"quote\",[\"a\",\"b\"]],\"expected\":\"usage: quote TEXT\"},{\"id\":\"AC-001:5\",\"input\":[\"unquote\",[\"--\"]],\"expected\":\"unknown option \\\"--\\\"\"},{\"id\":\"AC-001:6\",\"input\":[\"quote\",[\"--help\"]],\"expected\":\"unknown option \\\"--help\\\"\"},{\"id\":\"AC-001:7\",\"input\":[\"unquote\",[\"a\\\"b\",\"--oops\",\"--next\"]],\"expected\":\"unknown option \\\"--oops\\\"\"},{\"id\":\"AC-001:8\",\"input\":[\"quote\",[\"a\",\"b\",\"--all=1\"]],\"expected\":\"unknown option \\\"--all=1\\\"\"},{\"id\":\"AC-001:9\",\"input\":[\"unquote\",[\"-\"]],\"expected\":null},{\"id\":\"AC-001:10\",\"input\":[\"quote\",[\"-\"]],\"expected\":null},{\"id\":\"AC-003:1\",\"input\":\"a\\\"b\",\"expected\":\"unquote: CQ004 unterminated quote\"},{\"id\":\"AC-003:2\",\"input\":\"\\\"a\\nb\\\"\",\"expected\":\"unquote: CQ004 unterminated quote\"},{\"id\":\"AC-003:3\",\"input\":\"a\\r\\nb\",\"expected\":\"unquote: CQ006 text after the end of the value\"},{\"id\":\"AC-003:4\",\"input\":\"a\\nb\",\"expected\":\"unquote: CQ006 text after the end of the value\"},{\"id\":\"AC-003:5\",\"input\":\"a\\\\\\rb\",\"expected\":\"unquote: CQ005 unknown escape \\\"\\\\?\\\"\"},{\"id\":\"AC-003:6\",\"input\":\"x\\\\ \",\"expected\":\"unquote: CQ005 unknown escape \\\"\\\\?\\\"\"},{\"id\":\"AC-003:7\",\"input\":\"x\\\\q\",\"expected\":\"unquote: CQ005 unknown escape \\\"\\\\q\\\"\"}]"), &rows); err != nil {
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
	var (
		command string
		args    []string
	)

	switch value := input.(type) {
	case string:
		command = "unquote"
		args = []string{value}
	case []any:
		if len(value) != 2 {
			return nil, errors.New("invalid acceptance command input")
		}
		var ok bool
		command, ok = value[0].(string)
		if !ok {
			return nil, errors.New("invalid acceptance command name")
		}
		values, ok := value[1].([]any)
		if !ok {
			return nil, errors.New("invalid acceptance command arguments")
		}
		args = make([]string, len(values))
		for index, argument := range values {
			args[index], ok = argument.(string)
			if !ok {
				return nil, errors.New("invalid acceptance command argument")
			}
		}
	default:
		return nil, errors.New("invalid acceptance input")
	}

	var run func([]string, base.IO) (int, error)
	switch command {
	case "unquote":
		run = RunUnquote
	case "quote":
		run = RunQuote
	default:
		return nil, errors.New("invalid acceptance command")
	}

	memory := &base.MemIO{}
	status, err := run(args, memory)
	if memory.Stderr != "" {
		return nil, errors.New("unexpected stderr")
	}
	if err != nil {
		if status != 0 {
			return nil, errors.New("unexpected status on error")
		}
		if memory.Stdout != "" {
			return nil, errors.New("unexpected stdout on error")
		}
		return err.Error(), nil
	}
	if status != 0 {
		return nil, errors.New("unexpected status on success")
	}
	return nil, nil
}

func TestRunUnquote(t *acceptancetesting.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty argument", "", ""},
		{"only spaces", "  ", ""},
		{"empty quotes", "\"\"", ""},
		{"comment", "# comment", ""},
		{"tab inside value", "a\tb", "a\tb"},
		{"quoted edge spaces", "\"  a b  \"", "  a b  "},
		{"empty quotes drop leading blank", "\"\" a", "a"},
		{"empty quotes keep preceding blank", "a \"\"", "a "},
		{"empty quotes retain blanks around text", "a \"\" b", "a  b"},
		{"adjacent quotes", "\"\"\"a\"", "a"},
		{"comment after text", "a#b", "a"},
		{"terminal backslash", "a\\", "a"},
		{"continued line", "a\\\nb", "ab"},
		{"continued line in quotes", "\"a\\\n  b\"", "a  b"},
		{"terminal LF", "a\n", "a"},
		{"terminal CRLF", "a\r\n", "a"},
		{"bare CR inside text", "a\rb", "a\rb"},
		{"terminal bare CR", "a\r", "a"},
		{"quoted bare CR", "\"a\r\"", "a\r"},
		{"NBSP is not a blank", "x ", "x "},
		{"ideographic space is not a blank", "　x", "　x"},
		{"line separator is ordinary", "\" \"", " "},
		{"supported escapes", "\\t\\b\\n\\\\\\\"", "\t\b\n\\\""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *acceptancetesting.T) {
			memory := &base.MemIO{}
			status, err := RunUnquote([]string{test.input}, memory)
			if status != 0 {
				t.Errorf("status = %d, want 0", status)
			}
			if err != nil {
				t.Errorf("error = %v, want nil", err)
			}
			if memory.Stdout != test.want {
				t.Errorf("stdout = %q, want %q", memory.Stdout, test.want)
			}
			if memory.Stderr != "" {
				t.Errorf("stderr = %q, want empty", memory.Stderr)
			}
		})
	}
}

func TestRunQuote(t *acceptancetesting.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty argument", "", ""},
		{"space is quoted", " ", "\" \""},
		{"leading space is quoted", " a", "\" a\""},
		{"trailing space is quoted", "a ", "\"a \""},
		{"quote is escaped", "a\"b", "a\\\"b"},
		{"backslash is escaped", "a\\b", "a\\\\b"},
		{"LF is escaped", "a\nb", "a\\nb"},
		{"bare CR is quoted", "a\rb", "\"a\rb\""},
		{"trailing CR is quoted", "a\r", "\"a\r\""},
		{"backspace stays raw", "a\bb", "a\bb"},
		{"NBSP is not a quote trigger", "x ", "x "},
		{"line separator stays raw", "a b", "a b"},
		{"edge tabs are escaped without quotes", "\ta\t", "\\ta\\t"},
		{"semicolon is quoted", "a;b", "\"a;b\""},
		{"hash is quoted", "a#b", "\"a#b\""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *acceptancetesting.T) {
			memory := &base.MemIO{}
			status, err := RunQuote([]string{test.input}, memory)
			if status != 0 {
				t.Errorf("status = %d, want 0", status)
			}
			if err != nil {
				t.Errorf("error = %v, want nil", err)
			}
			if memory.Stdout != test.want {
				t.Errorf("stdout = %q, want %q", memory.Stdout, test.want)
			}
			if memory.Stderr != "" {
				t.Errorf("stderr = %q, want empty", memory.Stderr)
			}
		})
	}
}

func TestTextCommandUsageAndSingleDash(t *acceptancetesting.T) {
	if UnquoteUsage != "unquote TEXT" {
		t.Errorf("UnquoteUsage = %q, want %q", UnquoteUsage, "unquote TEXT")
	}
	if QuoteUsage != "quote TEXT" {
		t.Errorf("QuoteUsage = %q, want %q", QuoteUsage, "quote TEXT")
	}

	commands := []struct {
		name string
		run  func([]string, base.IO) (int, error)
	}{
		{"unquote", RunUnquote},
		{"quote", RunQuote},
	}
	for _, command := range commands {
		t.Run(command.name, func(t *acceptancetesting.T) {
			memory := &base.MemIO{}
			status, err := command.run([]string{"-"}, memory)
			if status != 0 {
				t.Errorf("status = %d, want 0", status)
			}
			if err != nil {
				t.Errorf("error = %v, want nil", err)
			}
			if memory.Stdout != "-" {
				t.Errorf("stdout = %q, want %q", memory.Stdout, "-")
			}
			if memory.Stderr != "" {
				t.Errorf("stderr = %q, want empty", memory.Stderr)
			}
		})
	}
}
