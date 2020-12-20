package main

import (
	"os/exec"
	"strings"
)

type Commander interface {
	Exec(cmd *exec.Cmd) (string, error)
	ExecSilently(cmd *exec.Cmd) error
}

type DefaultCommander struct {
}

func (c DefaultCommander) Exec(cmd *exec.Cmd) (string, error) {
	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(string(output), "\n"), nil
}

func (c DefaultCommander) ExecSilently(cmd *exec.Cmd) error {
	return cmd.Run()
}
