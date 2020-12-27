package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

func fakeCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(command, cs...)
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	return cmd
}

func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)

	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No command\n")
		os.Exit(2)
	}

	cmd, args := args[0], args[1:]
	fmt.Println(cmd)
	switch cmd {
	case "echo":
		iargs := []interface{}{}
		for _, s := range args {
			iargs = append(iargs, s)
		}
		fmt.Println(iargs...)

	case "exit":
		n, _ := strconv.Atoi(args[0])
		os.Exit(n)
	}
}

func TestExec(t *testing.T) {
	commander := DefaultCommander{}
	t.Run("test execute echo", func(t *testing.T) {
		cmd := fakeCommand("echo", "1")
		output, err := commander.Exec(cmd)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}

		if output != "1" {
			t.Errorf("expected 1, got %q", output)
		}

	})

	t.Run("test outputs error", func(t *testing.T) {
		cmd := fakeCommand("exit", "42")
		_, err := commander.Exec(cmd)
		expected := &ShellError{strings.Join(cmd.Args, " "), err}

		if err != expected {
			t.Errorf("expected %v, got %v", expected, err)
		}
	})
}
