package main

import (
	"io"
	"os/exec"
)

type NoLimiter struct {
	BinPath string
	Stdin   io.Reader
	Stdout  io.Writer
}

func NewNoLimiter(binPath string) *NoLimiter {
	return &NoLimiter{BinPath: binPath}
}

func (nl *NoLimiter) SetStdin(stdin io.Reader) {
	nl.Stdin = stdin
}

func (nl *NoLimiter) SetStdout(stdout io.Writer) {
	nl.Stdout = stdout
}
func (nl *NoLimiter) Run() (*LimiterResult, error) {
	cmd := exec.Command(nl.BinPath)
	cmd.Stdin = nl.Stdin
	cmd.Stdout = nl.Stdout
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return &LimiterResult{Result: LIMITER_RESULT_RUN_SUCCESSFUL}, nil
}
