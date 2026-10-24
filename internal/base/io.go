package base

import (
	"errors"
	"unicode/utf8"
)

// IO is the input and output boundary used by commands.
type IO interface {
	ReadStdin() ([]byte, error)
	ReadFile(path string) ([]byte, error)
	Out(text string)
	Err(text string)
}

// Source contains the display name and unchanged text of an input.
type Source struct {
	Display string
	Text    string
}

// Inputs reads named files and standard input, validating each byte sequence.
type Inputs struct {
	io        IO
	stdinRead bool
}

// NewInputs creates a reader that tracks standard input use.
func NewInputs(io IO) *Inputs {
	return &Inputs{io: io}
}

// Read reads one argument as a file path or standard input.
func (in *Inputs) Read(arg string) (Source, error) {
	display := Plain(arg)
	var (
		data []byte
		err  error
	)

	if arg == "-" {
		display = "<stdin>"
		if in.stdinRead {
			return Source{}, errors.New("standard input is named twice")
		}
		// Consume stdin before attempting IO so failures cannot be retried.
		in.stdinRead = true
		data, err = in.io.ReadStdin()
	} else {
		data, err = in.io.ReadFile(arg)
	}

	if err != nil {
		return Source{}, errors.New(display + ": cannot read the file")
	}
	if !utf8.Valid(data) {
		return Source{}, errors.New(display + ": input is not valid UTF-8")
	}

	return Source{Display: display, Text: string(data)}, nil
}

// MemIO is an in-memory implementation of IO for deterministic use and tests.
type MemIO struct {
	Stdin  []byte
	Files  map[string][]byte
	Stdout string
	Stderr string
}

// ReadStdin returns the configured standard input bytes.
func (m *MemIO) ReadStdin() ([]byte, error) {
	return m.Stdin, nil
}

// ReadFile returns a configured file or an error when the path is absent.
func (m *MemIO) ReadFile(path string) ([]byte, error) {
	data, ok := m.Files[path]
	if !ok {
		return nil, errors.New("file not found")
	}
	return data, nil
}

// Out appends text to standard output.
func (m *MemIO) Out(text string) {
	m.Stdout += text
}

// Err appends text to standard error.
func (m *MemIO) Err(text string) {
	m.Stderr += text
}
