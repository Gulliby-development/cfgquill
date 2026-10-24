package main

import (
	"io"
	"os"

	"github.com/Gulliby-development/cfgquill/internal/base"
	"github.com/Gulliby-development/cfgquill/internal/cli"
)

type standardIO struct{}

func (standardIO) ReadStdin() ([]byte, error) {
	return io.ReadAll(os.Stdin)
}

func (standardIO) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (standardIO) Out(text string) {
	_, _ = os.Stdout.Write([]byte(text))
}

func (standardIO) Err(text string) {
	_, _ = os.Stderr.Write([]byte(text))
}

func main() {
	os.Exit(cli.Run(os.Args[1:], base.IO(standardIO{})))
}
