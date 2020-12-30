package main

import (
	"errors"
	"strings"

	"github.com/spf13/pflag"
)

const (
	CommandStart = "start"
	CommandStop  = "stop"
)

type Options struct {
	Command string
	Project string
	Windows []string
	Attach  bool
	Debug   bool
}

var ErrHelp = errors.New("help requested")

const (
	WindowsUsage = "List of windows to start. If session exists, those windows will be attached to current session."
	AttachUsage  = "Force switch client for a session"
	DebugUsage   = "Print all commands to ~/.config/smug/smug.log"
)

func ParseOptions(argv []string, helpRequested func()) (Options, error) {
	if len(argv) < 2 {
		helpRequested()
		return Options{}, ErrHelp
	}

	cmd := argv[0]
	project := argv[1]

	flags := pflag.NewFlagSet(cmd, 0)
	windows := flags.StringArrayP("windows", "w", []string{}, WindowsUsage)
	attach := flags.BoolP("attach", "a", false, AttachUsage)
	debug := flags.BoolP("debug", "d", false, DebugUsage)

	err := flags.Parse(argv)
	if err == pflag.ErrHelp {
		return Options{}, ErrHelp
	}

	if err != nil {
		return Options{}, err
	}

	if strings.Contains(project, ":") {
		parts := strings.Split(project, ":")
		project = parts[0]
		wl := strings.Split(parts[1], ",")
		windows = &wl
	}

	return Options{cmd, project, *windows, *attach, *debug}, nil
}
