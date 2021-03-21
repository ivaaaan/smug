package main

import (
	"errors"
	"strings"

	"github.com/spf13/pflag"
)

const (
	CommandStart = "start"
	CommandStop  = "stop"
	CommandNew   = "new"
	CommandEdit  = "edit"
	CommandList  = "list"
	CommandPrint = "print"
)

var validCommands = []string{CommandStart, CommandStop, CommandNew, CommandEdit, CommandList, CommandPrint}

type Options struct {
	Command string
	Project string
	Config  string
	Windows []string
	Attach  bool
	Debug   bool
}

var ErrHelp = errors.New("help requested")

const (
	WindowsUsage = "List of windows to start. If session exists, those windows will be attached to current session."
	AttachUsage  = "Force switch client for a session"
	DebugUsage   = "Print all commands to ~/.config/smug/smug.log"
	FileUsage    = "A custom path to a config file"
)

// Creates a new FlagSet.
// Moved it to a variable to be able to override it in the tests.
var NewFlagSet = func(cmd string) *pflag.FlagSet {
	f := pflag.NewFlagSet(cmd, pflag.ContinueOnError)
	return f
}

func ParseOptions(argv []string, helpRequested func()) (Options, error) {
	if len(argv) == 0 {
		helpRequested()
		return Options{}, ErrHelp
	}

	if argv[0] == "--help" || argv[0] == "-h" {
		helpRequested()
		return Options{}, ErrHelp
	}

	cmd := argv[0]
	if !Contains(validCommands, cmd) {
		helpRequested()
		return Options{}, ErrHelp
	}

	flags := NewFlagSet(cmd)

	config := flags.StringP("file", "f", "", FileUsage)
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

	var project string
	if *config == "" && len(argv) > 1 {
		project = argv[1]
	}

	if strings.Contains(project, ":") {
		parts := strings.Split(project, ":")
		project = parts[0]
		wl := strings.Split(parts[1], ",")
		windows = &wl
	}

	return Options{cmd, project, *config, *windows, *attach, *debug}, nil
}
