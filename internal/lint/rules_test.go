package lint

// BEGIN acceptance rows imports: host-written from the audited spec table. Edit acceptanceActual, not this block.
import (
	acceptancejson "encoding/json"
	acceptancereflect "reflect"
	acceptancetesting "testing"
)

// END acceptance rows imports

import (
	"testing"

	"github.com/Gulliby-development/cfgquill/internal/config"
	"github.com/Gulliby-development/cfgquill/internal/problem"
)

// BEGIN acceptance rows: host-written from the audited spec table. Edit acceptanceActual, not this block.
func TestAcceptanceRowsRules(t *acceptancetesting.T) {
	var rows []struct {
		ID       string `json:"id"`
		Input    any    `json:"input"`
		Expected any    `json:"expected"`
	}
	if err := acceptancejson.Unmarshal([]byte("[{\"id\":\"AC-002:1\",\"input\":[\"\",{\"Entries\":[{\"HasValue\":true,\"Key\":\"a.b\",\"Line\":2,\"Value\":\"1\"},{\"HasValue\":true,\"Key\":\"a.b\",\"Line\":3,\"Value\":\"1\"},{\"HasValue\":true,\"Key\":\"a.b\",\"Line\":4,\"Value\":\"2\"},{\"HasValue\":true,\"Key\":\"a.b\",\"Line\":5,\"Value\":\"1\"}]}],\"expected\":[{\"Code\":\"CQ101\",\"Line\":3,\"Text\":\"duplicate entry \\\"a.b\\\" (first at line 2)\"},{\"Code\":\"CQ101\",\"Line\":5,\"Text\":\"duplicate entry \\\"a.b\\\" (first at line 2)\"}]},{\"id\":\"AC-002:2\",\"input\":[\"\",{\"Entries\":[{\"Key\":\"b\",\"Line\":1},{\"HasValue\":true,\"Key\":\"b\",\"Line\":2},{\"Key\":\"b\",\"Line\":3},{\"HasValue\":true,\"Key\":\"b\",\"Line\":4}]}],\"expected\":[{\"Code\":\"CQ101\",\"Line\":3,\"Text\":\"duplicate entry \\\"b\\\" (first at line 1)\"},{\"Code\":\"CQ101\",\"Line\":4,\"Text\":\"duplicate entry \\\"b\\\" (first at line 2)\"}]},{\"id\":\"AC-002:3\",\"input\":[\"\",{\"Headers\":[{\"Count\":0,\"Line\":1,\"Name\":\"a\"},{\"Count\":0,\"Line\":2,\"Name\":\"b.x y\"},{\"Count\":1,\"Line\":3,\"Name\":\"c\"}]}],\"expected\":[{\"Code\":\"CQ102\",\"Line\":1,\"Text\":\"empty section \\\"a\\\"\"},{\"Code\":\"CQ102\",\"Line\":2,\"Text\":\"empty section \\\"b.x y\\\"\"}]},{\"id\":\"AC-003:1\",\"input\":[\"\",{\"Entries\":[{\"HasValue\":true,\"Key\":\"a.b\",\"Line\":2,\"Value\":\" x\"}]}],\"expected\":[{\"Code\":\"CQ103\",\"Line\":2,\"Text\":\"blank at the edge of value \\\" x\\\"\"}]},{\"id\":\"AC-003:2\",\"input\":[\"\",{\"Entries\":[{\"HasValue\":true,\"Key\":\"a.b\",\"Line\":2,\"Value\":\"\\t\"},{\"HasValue\":true,\"Key\":\"a.c\",\"Line\":3,\"Value\":\"x\\r\"}]}],\"expected\":[{\"Code\":\"CQ103\",\"Line\":2,\"Text\":\"blank at the edge of value \\\"?\\\"\"},{\"Code\":\"CQ103\",\"Line\":3,\"Text\":\"blank at the edge of value \\\"x?\\\"\"}]},{\"id\":\"AC-003:3\",\"input\":[\"\",{\"Entries\":[{\"HasValue\":true,\"Line\":1,\"Value\":\"x\\ty\"},{\"HasValue\":true,\"Line\":2,\"Value\":\"\"},{\"HasValue\":true,\"Line\":3,\"Value\":\"\\nx\\n\"}]}],\"expected\":[]},{\"id\":\"AC-003:4\",\"input\":[\"\",{\"Entries\":[{\"HasValue\":true,\"Key\":\"a.b\",\"Line\":2,\"Value\":\"\u2028\u2029\u0085 ​　\\ufeff \"}]}],\"expected\":[{\"Code\":\"CQ104\",\"Line\":2,\"Text\":\"invisible character U+2028 (LS) in the value of \\\"a.b\\\"\"},{\"Code\":\"CQ104\",\"Line\":2,\"Text\":\"invisible character U+2029 (PS) in the value of \\\"a.b\\\"\"},{\"Code\":\"CQ104\",\"Line\":2,\"Text\":\"invisible character U+0085 (NEL) in the value of \\\"a.b\\\"\"},{\"Code\":\"CQ104\",\"Line\":2,\"Text\":\"invisible character U+00A0 (NBSP) in the value of \\\"a.b\\\"\"},{\"Code\":\"CQ104\",\"Line\":2,\"Text\":\"invisible character U+200B (ZWSP) in the value of \\\"a.b\\\"\"},{\"Code\":\"CQ104\",\"Line\":2,\"Text\":\"invisible character U+3000 (IDEOSP) in the value of \\\"a.b\\\"\"},{\"Code\":\"CQ104\",\"Line\":2,\"Text\":\"invisible character U+FEFF (BOM) in the value of \\\"a.b\\\"\"}]},{\"id\":\"AC-004:1\",\"input\":[\"\\ufeff#😀\\r\\r\\n#\u2028\\r\",{}],\"expected\":[{\"Code\":\"CQ105\",\"Line\":1,\"Text\":\"bare CR at column 4\"},{\"Code\":\"CQ105\",\"Line\":2,\"Text\":\"bare CR at column 3\"}]},{\"id\":\"AC-004:2\",\"input\":[\"\\r\\r\",{\"Headers\":[{\"Line\":1,\"Name\":\"a\"}]}],\"expected\":[{\"Code\":\"CQ102\",\"Line\":1,\"Text\":\"empty section \\\"a\\\"\"},{\"Code\":\"CQ105\",\"Line\":1,\"Text\":\"bare CR at column 1\"},{\"Code\":\"CQ105\",\"Line\":1,\"Text\":\"bare CR at column 2\"}]},{\"id\":\"AC-004:3\",\"input\":[\"[a]\\r\\nb=1\\r\\n\",{}],\"expected\":[]}]"), &rows); err != nil {
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

func TestAcceptanceRowsLiteralCases(t *acceptancetesting.T) {
	tests := []struct {
		name  string
		input string
		want  []problem.Problem
	}{
		{
			name:  "AC-002 first identical origin",
			input: `["",{"Entries":[{"Line":2,"Key":"a.b","Value":"1","HasValue":true},{"Line":3,"Key":"a.b","Value":"1","HasValue":true},{"Line":4,"Key":"a.b","Value":"2","HasValue":true},{"Line":5,"Key":"a.b","Value":"1","HasValue":true}]}]`,
			want: []problem.Problem{
				{Line: 3, Code: "CQ101", Text: `duplicate entry "a.b" (first at line 2)`},
				{Line: 5, Code: "CQ101", Text: `duplicate entry "a.b" (first at line 2)`},
			},
		},
		{
			name:  "AC-002 absent versus empty",
			input: `["",{"Entries":[{"Line":1,"Key":"b"},{"Line":2,"Key":"b","HasValue":true},{"Line":3,"Key":"b"},{"Line":4,"Key":"b","HasValue":true}]}]`,
			want: []problem.Problem{
				{Line: 3, Code: "CQ101", Text: `duplicate entry "b" (first at line 1)`},
				{Line: 4, Code: "CQ101", Text: `duplicate entry "b" (first at line 2)`},
			},
		},
		{
			name:  "AC-002 empty headers",
			input: `["",{"Headers":[{"Line":1,"Name":"a","Count":0},{"Line":2,"Name":"b.x y","Count":0},{"Line":3,"Name":"c","Count":1}]}]`,
			want: []problem.Problem{
				{Line: 1, Code: "CQ102", Text: `empty section "a"`},
				{Line: 2, Code: "CQ102", Text: `empty section "b.x y"`},
			},
		},
		{
			name:  "AC-003 leading space",
			input: `["",{"Entries":[{"Line":2,"Key":"a.b","Value":" x","HasValue":true}]}]`,
			want:  []problem.Problem{{Line: 2, Code: "CQ103", Text: `blank at the edge of value " x"`}},
		},
		{
			name:  "AC-003 tab and CR edges",
			input: `["",{"Entries":[{"Line":2,"Key":"a.b","Value":"\t","HasValue":true},{"Line":3,"Key":"a.c","Value":"x\r","HasValue":true}]}]`,
			want: []problem.Problem{
				{Line: 2, Code: "CQ103", Text: `blank at the edge of value "?"`},
				{Line: 3, Code: "CQ103", Text: `blank at the edge of value "x?"`},
			},
		},
		{
			name:  "AC-003 non-edge blanks",
			input: `["",{"Entries":[{"Line":1,"Value":"x\ty","HasValue":true},{"Line":2,"Value":"","HasValue":true},{"Line":3,"Value":"\nx\n","HasValue":true}]}]`,
			want:  []problem.Problem{},
		},
		{
			name:  "AC-003 ordered distinct marks",
			input: `["",{"Entries":[{"Line":2,"Key":"a.b","Value":"\u2028\u2029\u0085\u00a0\u200b\u3000\ufeff\u00a0","HasValue":true}]}]`,
			want: []problem.Problem{
				{Line: 2, Code: "CQ104", Text: `invisible character U+2028 (LS) in the value of "a.b"`},
				{Line: 2, Code: "CQ104", Text: `invisible character U+2029 (PS) in the value of "a.b"`},
				{Line: 2, Code: "CQ104", Text: `invisible character U+0085 (NEL) in the value of "a.b"`},
				{Line: 2, Code: "CQ104", Text: `invisible character U+00A0 (NBSP) in the value of "a.b"`},
				{Line: 2, Code: "CQ104", Text: `invisible character U+200B (ZWSP) in the value of "a.b"`},
				{Line: 2, Code: "CQ104", Text: `invisible character U+3000 (IDEOSP) in the value of "a.b"`},
				{Line: 2, Code: "CQ104", Text: `invisible character U+FEFF (BOM) in the value of "a.b"`},
			},
		},
		{
			name:  "AC-004 rune columns and LF-only lines",
			input: `["\ufeff#😀\r\r\n#\u2028\r",{}]`,
			want: []problem.Problem{
				{Line: 1, Code: "CQ105", Text: "bare CR at column 4"},
				{Line: 2, Code: "CQ105", Text: "bare CR at column 3"},
			},
		},
		{
			name:  "AC-004 code ordering and stable ties",
			input: `["\r\r",{"Headers":[{"Line":1,"Name":"a"}]}]`,
			want: []problem.Problem{
				{Line: 1, Code: "CQ102", Text: `empty section "a"`},
				{Line: 1, Code: "CQ105", Text: "bare CR at column 1"},
				{Line: 1, Code: "CQ105", Text: "bare CR at column 2"},
			},
		},
		{
			name:  "AC-004 excludes CRLF",
			input: `["[a]\r\nb=1\r\n",{}]`,
			want:  []problem.Problem{},
		},
	}
	for _, test := range tests {
		var input any
		if err := acceptancejson.Unmarshal([]byte(test.input), &input); err != nil {
			t.Fatal(err)
		}
		actual, err := acceptanceActual(input)
		if err != nil {
			t.Fatal(err)
		}
		if got := actual.([]problem.Problem); !acceptancereflect.DeepEqual(got, test.want) {
			t.Errorf("%s: Findings = %#v, want %#v", test.name, got, test.want)
		}
	}
}

// acceptanceActual decodes the harness input into the public configuration model.
func acceptanceActual(input any) (any, error) {
	encoded, err := acceptancejson.Marshal(input)
	if err != nil {
		return nil, err
	}
	var actual []acceptancejson.RawMessage
	if err := acceptancejson.Unmarshal(encoded, &actual); err != nil {
		return nil, err
	}
	var text string
	if err := acceptancejson.Unmarshal(actual[0], &text); err != nil {
		return nil, err
	}
	var file config.File
	if err := acceptancejson.Unmarshal(actual[1], &file); err != nil {
		return nil, err
	}
	return Findings(text, file), nil
}

func TestFindingsEmpty(t *testing.T) {
	got := Findings("", config.File{})
	want := []problem.Problem{}
	if !acceptancereflect.DeepEqual(got, want) {
		t.Fatalf("Findings(empty) = %#v, want %#v", got, want)
	}
}

func TestFindingsTrailingSpace(t *testing.T) {
	file := config.File{Entries: []config.Entry{{Line: 2, Key: "a.b", Value: "x ", HasValue: true}}}
	got := Findings("", file)
	want := []problem.Problem{{Line: 2, Code: "CQ103", Text: "blank at the edge of value \"x \""}}
	if !acceptancereflect.DeepEqual(got, want) {
		t.Fatalf("Findings(trailing space) = %#v, want %#v", got, want)
	}
}

func TestFindingsUnicodeEdgesAreNotBlank(t *testing.T) {
	file := config.File{Entries: []config.Entry{{Line: 2, Key: "a.b", Value: "\u00a0middle\u3000", HasValue: true}}}
	got := Findings("", file)
	want := []problem.Problem{
		{Line: 2, Code: "CQ104", Text: "invisible character U+00A0 (NBSP) in the value of \"a.b\""},
		{Line: 2, Code: "CQ104", Text: "invisible character U+3000 (IDEOSP) in the value of \"a.b\""},
	}
	if !acceptancereflect.DeepEqual(got, want) {
		t.Fatalf("Findings(Unicode edges) = %#v, want %#v", got, want)
	}
}

func TestFindingsMixedOrdering(t *testing.T) {
	text := "x\n#\r"
	file := config.File{
		Entries: []config.Entry{
			{Line: 2, Key: "a.b", Value: "\tx\t\u200b", HasValue: true},
			{Line: 2, Key: "a.b", Value: "\tx\t\u200b", HasValue: true},
		},
		Headers: []config.Header{{Line: 1, Name: "empty", Count: 0}},
	}
	got := Findings(text, file)
	want := []problem.Problem{
		{Line: 1, Code: "CQ102", Text: "empty section \"empty\""},
		{Line: 2, Code: "CQ101", Text: "duplicate entry \"a.b\" (first at line 2)"},
		{Line: 2, Code: "CQ103", Text: "blank at the edge of value \"?x??\""},
		{Line: 2, Code: "CQ103", Text: "blank at the edge of value \"?x??\""},
		{Line: 2, Code: "CQ104", Text: "invisible character U+200B (ZWSP) in the value of \"a.b\""},
		{Line: 2, Code: "CQ104", Text: "invisible character U+200B (ZWSP) in the value of \"a.b\""},
		{Line: 2, Code: "CQ105", Text: "bare CR at column 2"},
	}
	if !acceptancereflect.DeepEqual(got, want) {
		t.Fatalf("Findings(mixed ordering) = %#v, want %#v", got, want)
	}
}

func TestFindingsAcceptanceRawCRRows(t *testing.T) {
	tests := []struct {
		input string
		want  []problem.Problem
	}{
		{
			input: "[\"\\ufeff#😀\\r\\r\\n#\\u2028\\r\",{}]",
			want: []problem.Problem{
				{Line: 1, Code: "CQ105", Text: "bare CR at column 4"},
				{Line: 2, Code: "CQ105", Text: "bare CR at column 3"},
			},
		},
		{
			input: "[\"[a]\\r\\nb=1\\r\\n\",{}]",
			want:  []problem.Problem{},
		},
	}
	for _, test := range tests {
		var input any
		if err := acceptancejson.Unmarshal([]byte(test.input), &input); err != nil {
			t.Fatal(err)
		}
		actual, err := acceptanceActual(input)
		if err != nil {
			t.Fatal(err)
		}
		if got := actual.([]problem.Problem); !acceptancereflect.DeepEqual(got, test.want) {
			t.Errorf("Findings(%s) = %#v, want %#v", test.input, got, test.want)
		}
	}
}

func TestFindingsEscapedBOMRawCRAcceptanceInput(t *testing.T) {
	input := "[\"\\ufeff#😀\\r\\r\\n#\\u2028\\r\",{}]"
	var decoded any
	if err := acceptancejson.Unmarshal([]byte(input), &decoded); err != nil {
		t.Fatal(err)
	}
	actual, err := acceptanceActual(decoded)
	if err != nil {
		t.Fatal(err)
	}
	got := actual.([]problem.Problem)
	want := []problem.Problem{
		{Line: 1, Code: "CQ105", Text: "bare CR at column 4"},
		{Line: 2, Code: "CQ105", Text: "bare CR at column 3"},
	}
	if !acceptancereflect.DeepEqual(got, want) {
		t.Fatalf("Findings(%s) = %#v, want %#v", input, got, want)
	}
}
