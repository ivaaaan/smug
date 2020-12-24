package main

import (
	"log"
	"os/exec"
	"strings"
)

type Commander interface {
	Exec(cmd *exec.Cmd) (string, error)
	ExecSilently(cmd *exec.Cmd) error
}

type DefaultCommander struct {
	logger *log.Logger
}

func (c DefaultCommander) Exec(cmd *exec.Cmd) (string, error) {
	if c.logger != nil {
		c.logger.Println(strings.Join(cmd.Args, " "))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		if c.logger != nil {
			c.logger.Println(err)
		}
		return "", &ShellError{strings.Join(cmd.Args, " "), err}
	}

	return strings.TrimSuffix(string(output), "\n"), nil
}

func (c DefaultCommander) ExecSilently(cmd *exec.Cmd) error {
	if c.logger != nil {
		c.logger.Println(strings.Join(cmd.Args, " "))
	}

	err := cmd.Run()
	if err != nil {
		if c.logger != nil {
			c.logger.Println(err)
		}
		return &ShellError{strings.Join(cmd.Args, " "), err}
	}
	return nil
}
