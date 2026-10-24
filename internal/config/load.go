package config

import (
	"errors"

	"github.com/Gulliby-development/cfgquill/internal/base"
	"github.com/Gulliby-development/cfgquill/internal/problem"
)

// Loaded retains the source and parsed model of a configuration file.
type Loaded struct {
	Display string
	Text    string
	File    File
}

// Load reads and parses one input using the caller's shared input state.
func Load(in *base.Inputs, arg string) (Loaded, error) {
	source, err := in.Read(arg)
	if err != nil {
		return Loaded{}, err
	}

	file, fault := Parse(source.Text)
	if fault != nil {
		return Loaded{}, errors.New(problem.FromFault(fault).Render(source.Display))
	}

	return Loaded{Display: source.Display, Text: source.Text, File: file}, nil
}
