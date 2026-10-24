package problem

import (
	"reflect"
	"testing"

	"github.com/Gulliby-development/cfgquill/internal/base"
)

func TestFaultText(t *testing.T) {
	tests := []struct {
		code string
		arg  string
		want string
	}{
		{code: "CQ001", want: "bad section header"},
		{code: "CQ002", want: "bad variable name"},
		{code: "CQ003", want: "expected \"=\" after the variable name"},
		{code: "CQ004", want: "unterminated quote"},
		{code: "CQ005", arg: `\q`, want: `unknown escape "\q"`},
		{code: "CQ005", arg: "\\\r", want: `unknown escape "\?"`},
		{code: "CQ005", arg: "\\ ", want: `unknown escape "\?"`},
		{code: "CQ005", arg: "\\😀", want: `unknown escape "\?"`},
		{code: "CQ006", want: "text after the end of the value"},
	}
	for _, test := range tests {
		if got := FaultText(test.code, test.arg); got != test.want {
			t.Errorf("FaultText(%q, %q) = %q, want %q", test.code, test.arg, got, test.want)
		}
	}
}

func TestFromFault(t *testing.T) {
	fault := &base.Fault{Code: "CQ005", Line: 17, Arg: `\q`}
	got := FromFault(fault)
	want := Problem{Line: 17, Code: "CQ005", Text: `unknown escape "\q"`}
	if got != want {
		t.Fatalf("FromFault() = %#v, want %#v", got, want)
	}
	if *fault != (base.Fault{Code: "CQ005", Line: 17, Arg: `\q`}) {
		t.Fatalf("FromFault changed its fault: %#v", fault)
	}
}

func TestRender(t *testing.T) {
	tests := []struct {
		display string
		problem Problem
		want    string
	}{
		{display: "settings.cfg", problem: Problem{Line: 3, Code: "CQ002", Text: "bad variable name"}, want: "settings.cfg:3: CQ002 bad variable name"},
		{display: "", problem: Problem{Line: 1, Code: "CQ004", Text: "unterminated quote"}, want: ":1: CQ004 unterminated quote"},
		{display: "<stdin>", problem: Problem{Line: 12, Code: "CQ001", Text: "bad section header"}, want: "<stdin>:12: CQ001 bad section header"},
	}
	for _, test := range tests {
		if got := test.problem.Render(test.display); got != test.want {
			t.Errorf("Render(%q) = %q, want %q", test.display, got, test.want)
		}
	}
}

func TestSorted(t *testing.T) {
	tests := []struct {
		name string
		in   []Problem
		want []Problem
	}{
		{
			name: "mixed lines and codes with stable equal keys",
			in: []Problem{
				{Line: 3, Code: "CQ001", Text: "later line"},
				{Line: 1, Code: "CQ002", Text: "first code"},
				{Line: 1, Code: "CQ001", Text: "second equal-key item"},
				{Line: 1, Code: "CQ001", Text: "first equal-key item"},
				{Line: 1, Code: "CQ003", Text: "last code"},
				{Line: 2, Code: "CQ001", Text: "middle line"},
			},
			want: []Problem{
				{Line: 1, Code: "CQ001", Text: "second equal-key item"},
				{Line: 1, Code: "CQ001", Text: "first equal-key item"},
				{Line: 1, Code: "CQ002", Text: "first code"},
				{Line: 1, Code: "CQ003", Text: "last code"},
				{Line: 2, Code: "CQ001", Text: "middle line"},
				{Line: 3, Code: "CQ001", Text: "later line"},
			},
		},
		{name: "empty", in: []Problem{}, want: []Problem{}},
		{
			name: "singleton",
			in:   []Problem{{Line: 0, Code: "CQ006", Text: "only"}},
			want: []Problem{{Line: 0, Code: "CQ006", Text: "only"}},
		},
	}
	for _, test := range tests {
		before := make([]Problem, len(test.in))
		copy(before, test.in)
		got := Sorted(test.in)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s: Sorted() = %#v, want %#v", test.name, got, test.want)
		}
		if !reflect.DeepEqual(test.in, before) {
			t.Errorf("%s: Sorted() changed input to %#v, want %#v", test.name, test.in, before)
		}
		if len(got) > 0 {
			got[0].Text = "changed result"
			if !reflect.DeepEqual(test.in, before) {
				t.Errorf("%s: result does not have independent slice storage", test.name)
			}
		}
	}
}
