package main

import (
	"strings"

	"github.com/docopt/docopt-go"
)

const usage = `Smug - tmux session manager. Version v0.1.3.1

Usage:
	smug <command> <project> [-w <window>]... [--attach] [--debug]

Options:
	-w List of windows to start. If session exists, those windows will be attached to current session.

Examples:
	$ smug start blog
	$ smug start blog:win1
	$ smug start blog -w win1
	$ smug start blog:win1,win2
	$ smug stop blog
`

type Options struct {
	Command string
	Project string
	Windows []string
	Attach  bool
	Debug   bool
}

func ParseOptions(p docopt.Parser, argv []string) (Options, error) {
	arguments, err := p.ParseArgs(usage, argv, "")
	if err != nil {
		return Options{}, err
	}

	cmd, err := arguments.String("<command>")
	if err != nil {
		return Options{}, err
	}

	project, err := arguments.String("<project>")
	if err != nil {
		return Options{}, err
	}

	attach, err := arguments.Bool("--attach")
	if err != nil {
		return Options{}, err
	}

	debug, err := arguments.Bool("--debug")
	if err != nil {
		return Options{}, err
	}

	var windows []string

	if strings.Contains(project, ":") {
		parts := strings.Split(project, ":")
		project = parts[0]
		windows = strings.Split(parts[1], ",")
	} else {
		windows = arguments["-w"].([]string)
	}

	return Options{cmd, project, windows, attach, debug}, nil
}
